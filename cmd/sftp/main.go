package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	sshFxfRead   = 0x00000001
	sshFxfWrite  = 0x00000002
	sshFxfAppend = 0x00000004
	sshFxfCreat  = 0x00000008
	sshFxfTrunc  = 0x00000010
	sshFxfExcl   = 0x00000020
)

func resolveRoot(podman context.Context, username string) (string, error) {
	parts := strings.SplitN(username, ".", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid username format (expected server.user)")
	}

	server := parts[0]

	inspect, err := containers.Inspect(podman, server, nil)
	if err != nil {
		return "", fmt.Errorf("server not found: %w", err)
	}

	for _, mount := range inspect.Mounts {
		return mount.Source, nil
	}

	return "", fmt.Errorf("no volume found for server %s", server)
}

func secureJoin(root, reqPath string) (string, error) {
	clean := filepath.Clean("/" + reqPath)
	full := filepath.Join(root, clean)

	rel, err := filepath.Rel(root, full)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("invalid path")
	}

	return full, nil
}

type writerAt struct {
	*os.File
}

func (w *writerAt) WriteAt(p []byte, off int64) (int, error) {
	return w.File.WriteAt(p, off)
}

func (w *writerAt) Close() error {
	return w.File.Close()
}

type rootHandler struct {
	root string
}

func (h *rootHandler) Fileread(r *sftp.Request) (io.ReaderAt, error) {
	path, err := secureJoin(h.root, r.Filepath)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

func (h *rootHandler) Filewrite(r *sftp.Request) (io.WriterAt, error) {
	path, err := secureJoin(h.root, r.Filepath)
	if err != nil {
		return nil, err
	}

	// safer, GVFS-compatible flags
	flags := os.O_RDWR | os.O_CREATE

	// only truncate if explicitly requested
	if r.Flags&sshFxfTrunc != 0 {
		flags |= os.O_TRUNC
	}

	if r.Flags&sshFxfAppend != 0 {
		flags |= os.O_APPEND
	}

	f, err := os.OpenFile(path, flags, 0644)
	if err != nil {
		return nil, err
	}

	return &writerAt{f}, nil
}

func (h *rootHandler) Filecmd(r *sftp.Request) error {
	path, err := secureJoin(h.root, r.Filepath)
	if err != nil {
		return err
	}

	switch r.Method {
	case "Setstat":
		if r.Attributes() != nil {
			// ignore errors — Nautilus doesn't care
			_ = os.Chmod(path, os.FileMode(r.Attributes().Mode))
		}
		return nil
	case "Remove":
		return os.Remove(path)
	case "Mkdir":
		return os.Mkdir(path, 0755)
	case "Rmdir":
		return os.Remove(path)
	case "Rename":
		target, err := secureJoin(h.root, r.Target)
		if err != nil {
			return err
		}
		_ = os.Remove(target) // MUST allow overwrite
		return os.Rename(path, target)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

func (h *rootHandler) Filelist(r *sftp.Request) (sftp.ListerAt, error) {
	path, err := secureJoin(h.root, r.Filepath)
	if err != nil {
		return nil, err
	}

	switch r.Method {
	case "List":
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		var fis []os.FileInfo
		for _, e := range entries {
			fi, err := e.Info()
			if err != nil {
				continue
			}
			fis = append(fis, fi)
		}
		return listerAt(fis), nil

	case "Stat", "Lstat":
		fi, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		return listerAt([]os.FileInfo{fi}), nil

	default:
		return nil, fmt.Errorf("unsupported list method: %s", r.Method)
	}
}

type listerAt []os.FileInfo

func (l listerAt) ListAt(ls []os.FileInfo, offset int64) (int, error) {
	if offset >= int64(len(l)) {
		return 0, io.EOF
	}

	n := copy(ls, l[offset:])
	if n < len(ls) {
		return n, io.EOF
	}
	return n, nil
}

func main() {
	ctx := context.Background()
	podman, err := bindings.NewConnection(ctx, "unix:///run/user/1000/podman/podman.sock")
	if err != nil {
		log.Fatal(err)
	}

	key, err := os.ReadFile("/home/alex/.ssh/id_rsa")
	if err != nil {
		log.Fatal(err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal(err)
	}

	sshConfig := &ssh.ServerConfig{}
	sshConfig.AddHostKey(signer)
	sshConfig.PasswordCallback = func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
		return nil, nil
	}

	addr := "127.0.0.1:2222"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("listening on %s\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		go func() {
			defer conn.Close()

			sshConn, chans, reqs, err := ssh.NewServerConn(conn, sshConfig)
			if err != nil {
				fmt.Println("ssh server error:", err)
				return
			}
			defer sshConn.Close()

			fmt.Printf("connection from %s (%s)\n", sshConn.RemoteAddr(), sshConn.User())

			go ssh.DiscardRequests(reqs)

			// TODO
			root, err := resolveRoot(podman, sshConn.User())
			if err != nil {
				log.Println("resolve error:", err)
				return
			}

			for nc := range chans {
				if nc.ChannelType() != "session" {
					nc.Reject(ssh.UnknownChannelType, "unknown channel type")
					continue
				}

				c, reqs, err := nc.Accept()
				if err != nil {
					fmt.Println("channel accept error:", err)
					continue
				}

				go func() {
					defer c.Close()

					for req := range reqs {
						switch req.Type {
						case "subsystem":
							if string(req.Payload[4:]) == "sftp" {
								req.Reply(true, nil)

								handlers := sftp.Handlers{
									FileGet:  &rootHandler{root: root},
									FilePut:  &rootHandler{root: root},
									FileCmd:  &rootHandler{root: root},
									FileList: &rootHandler{root: root},
								}

								// server, err := sftp.NewServer(
								// 	c,
								// 	sftp.WithDebug(nil),
								// )
								// if err != nil {
								// 	log.Println("sftp server error:", err)
								// 	return
								// }
								server := sftp.NewRequestServer(c, handlers)

								if err := server.Serve(); err == io.EOF {
									server.Close()
								} else if err != nil {
									log.Println("sftp serve error:", err)
								}
								return
							}
						default:
							req.Reply(false, nil)
						}
					}
				}()
			}

		}()
	}
}
