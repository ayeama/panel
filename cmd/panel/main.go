package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"panel/internal/middleware"
	"panel/pkg/api"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/bindings/volumes"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/gorilla/websocket"
	"github.com/opencontainers/runtime-spec/specs-go"
	nettypes "go.podman.io/common/libnetwork/types"
	"go.yaml.in/yaml/v3"
)

type Server struct {
	podman *context.Context
}

type ImageHandler struct {
	srv *Server
}

func (h *ImageHandler) handle_read(w http.ResponseWriter, r *http.Request) {
	var imageFilters = map[string]struct{}{
		"label":     {},
		"reference": {},
	}
	filters := make(map[string][]string)
	filters["dangling"] = []string{"false"}
	filters["label"] = []string{"com.github.ayeama.panel.server.name"}
	for key, values := range r.URL.Query() {
		if _, ok := imageFilters[key]; !ok {
			continue
		}
		for _, v := range values {
			if v == "" {
				continue
			}
			// TODO validation
			filters[key] = append(filters[key], v)
		}
	}

	options := &images.ListOptions{
		Filters: filters,
	}
	l, err := images.List(*h.srv.podman, options)
	if err != nil {
		fmt.Println("image handle create list:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	items := make([]api.ImageResponse, 0, len(l))
	for _, i := range l {
		image, _ := images.GetImage(*h.srv.podman, i.ID, nil)

		ref := image.RepoTags[0] // TODO assume
		repo, _ := reference.Parse(ref)
		named, _ := repo.(reference.Named)
		name := named.Name()
		tagged, _ := repo.(reference.Tagged)
		tag := tagged.Tag()

		env := make(map[string]string)
		for _, e := range image.Config.Env {
			s := strings.SplitN(e, "=", 2)
			if strings.HasPrefix(s[0], "PANEL_") { // TODO prefix
				env[strings.TrimPrefix(s[0], "PANEL_")] = s[1]
			}
		}

		items = append(items, api.ImageResponse{
			Id:         i.ID,
			Reference:  repo.String(),
			Repository: name,
			Tag:        tag,
			Env:        env,
		})
	}

	response := api.ImageListResponse{
		Items: items,
	}

	json, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (h *ImageHandler) Register(m *http.ServeMux) {
	m.HandleFunc("GET /images", h.handle_read)
}

type ServerHandler struct {
	srv *Server
}

func (h *ServerHandler) handle_create(w http.ResponseWriter, r *http.Request) {
	var request api.ServerCreateRequest
	d := json.NewDecoder(r.Body)
	err := d.Decode(&request)
	if err != nil {
		fmt.Println("server handle create parse request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	img, err := images.GetImage(*h.srv.podman, request.Image, nil)
	if err != nil {
		fmt.Println("server handle create get image:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// s := specgen.NewSpecGenerator(request.Image, false)
	s := specgen.NewSpecGenerator(img.ID, false)

	s.Labels = map[string]string{"com.github.ayeama.panel.server.id": "TODO"} // TODO

	if s.Env == nil {
		s.Env = make(map[string]string)
	}
	for k, v := range request.Env {
		s.Env[fmt.Sprintf("PANEL_%s", k)] = v // TODO check prefix
	}

	cpus := 2.0
	cpuPeriod := uint64(100000)
	cpuQuota := int64(float64(cpuPeriod) * cpus)
	memLimit := int64(2000000000)

	s.ResourceLimits = &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Period: &cpuPeriod,
			Quota:  &cpuQuota,
		},
		Memory: &specs.LinuxMemory{
			Limit: &memLimit,
		},
	}

	stdin := true
	s.Stdin = &stdin

	terminal := true
	s.Terminal = &terminal

	// publish := true
	// s.PublishExposedPorts = &publish

	// TODO aweful logic
	s.PortMappings = make([]nettypes.PortMapping, 0, len(img.Config.ExposedPorts))
	for k := range img.Config.ExposedPorts {
		exposedPortParts := strings.Split(k, "/")
		containerPort, err := strconv.Atoi(exposedPortParts[0])
		if err != nil {
			// TODO handle error
		}
		minPort := 45000
		maxPort := 45000 + 1024
		hostPort := rand.Intn(maxPort-minPort+1) + minPort

		s.PortMappings = append(s.PortMappings, nettypes.PortMapping{
			ContainerPort: uint16(containerPort),
			HostPort:      uint16(hostPort),
			Protocol:      "tcp,udp",
		})
	}

	c, err := containers.CreateWithSpec(*h.srv.podman, s, nil)
	if err != nil {
		fmt.Println("server handle create container spec:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = containers.Start(*h.srv.podman, c.ID, nil)
	if err != nil {
		fmt.Println("server handle create start:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	i, err := containers.Inspect(*h.srv.podman, c.ID, nil)
	if err != nil {
		fmt.Println("server handle create container inspect:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ports := make([]string, 0)
	for port, hostports := range i.NetworkSettings.Ports {
		for _, hostport := range hostports {
			ports = append(ports, fmt.Sprintf("%s %s:%s", port, hostport.HostIP, hostport.HostPort))
		}
	}

	response := api.ServerResponse{
		Id:     i.ID,
		Name:   i.Name,
		Image:  img.RepoTags[0], // TODO check
		Status: i.State.Status,
		Ports:  ports,
	}

	json, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

func (h *ServerHandler) handle_read(w http.ResponseWriter, r *http.Request) {
	all := true
	filters := map[string][]string{"label": {"com.github.ayeama.panel.server.id"}} // TODO labels
	options := &containers.ListOptions{
		All:     &all,
		Filters: filters,
	}
	l, err := containers.List(*h.srv.podman, options)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	items := make([]api.ServerResponse, 0, len(l))
	for _, c := range l {
		items = append(items, api.ServerResponse{
			Id:     c.ID,
			Name:   c.Names[0], // TODO: check
			Image:  c.Image,
			Status: c.State, // NOTE: state not status
			// TODO ports
		})
	}

	response := api.ServerListResponse{
		Items: items,
	}

	json, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (h *ServerHandler) handle_read_one(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i, err := containers.Inspect(*h.srv.podman, id, nil)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	img, err := images.GetImage(*h.srv.podman, i.Image, nil)
	if err != nil {
		fmt.Println("server handle create get image:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ports := make([]string, 0, len(i.NetworkSettings.Ports))
	for _, hostports := range i.NetworkSettings.Ports {
		for _, hostport := range hostports {
			ports = append(ports, hostport.HostPort)
		}
	}

	response := api.ServerResponse{
		Id:     i.ID,
		Name:   i.Name,
		Image:  img.RepoTags[0], // TODO check
		Status: i.State.Status,
		Ports:  ports,
	}

	json, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (h *ServerHandler) handle_delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	depend := true
	force := true
	volumes := true
	timeout := uint(0)
	options := &containers.RemoveOptions{
		Depend:  &depend,
		Force:   &force,
		Volumes: &volumes,
		Timeout: &timeout,
	}
	_, err := containers.Remove(*h.srv.podman, id, options)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ServerHandler) handle_start(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := containers.Start(*h.srv.podman, id, nil)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	i, err := containers.Inspect(*h.srv.podman, id, nil)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ports := make([]string, 0)
	for port, hostports := range i.NetworkSettings.Ports {
		for _, hostport := range hostports {
			ports = append(ports, fmt.Sprintf("%s %s:%s", port, hostport.HostIP, hostport.HostPort))
		}
	}

	response := api.ServerResponse{
		Id:     i.ID,
		Name:   i.Name,
		Image:  i.Image,
		Status: i.State.Status,
		Ports:  ports,
	}

	json, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (h *ServerHandler) handle_stop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// i, err := containers.Inspect(*h.srv.podman, id, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// TODO first try sending stop command
	// img, err := images.GetImage(*h.srv.podman, i.Image, nil)
	// if err != nil {
	// 	fmt.Println("server handle stop get image:", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// cmd := img.Config.Labels["com.github.ayeama.panel.server.stop"]
	// if cmd != "" {
	// 	options := &handlers.ExecCreateConfig{
	// 		Cmd:          []string{cmd},
	// 		AttachStdout: false,
	// 		AttachStderr: false,
	// 		AttachStdin:  false,
	// 	}
	// 	exec, _ := containers.ExecCreate(*h.srv.podman, i.ID, options)
	// 	// TODO handle error

	// 	_ = containers.ExecStart(*h.srv.podman, exec, nil)
	// 	// TODO handle error

	// }

	timeout := uint(0)
	options := &containers.StopOptions{
		Timeout: &timeout,
	}
	err := containers.Stop(*h.srv.podman, id, options)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get the container again
	i, err := containers.Inspect(*h.srv.podman, id, nil)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ports := make([]string, 0)
	for port, hostports := range i.NetworkSettings.Ports {
		for _, hostport := range hostports {
			ports = append(ports, fmt.Sprintf("%s %s:%s", port, hostport.HostIP, hostport.HostPort))
		}
	}

	response := api.ServerResponse{
		Id:     i.ID,
		Name:   i.Name,
		Image:  i.Image,
		Status: i.State.Status,
		Ports:  ports,
	}

	json, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (h *ServerHandler) handle_attach(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("attach connection upgrade:", err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(3)

	ready := make(chan bool, 1)

	stdoutReader, stdoutWriter := io.Pipe()
	stdinReader, stdinWriter := io.Pipe()

	cleanup := sync.OnceFunc(func() {
		cancel()
		conn.Close()

		stdoutReader.Close()
		stdoutWriter.Close()

		stdinReader.Close()
		stdinWriter.Close()
	})

	logs := true
	stream := true
	options := &containers.AttachOptions{
		Logs:   &logs,
		Stream: &stream,
	}

	go func() {
		defer wg.Done()
		defer cleanup()

		err := containers.Attach(*h.srv.podman, id, stdinReader, stdoutWriter, stdoutWriter, ready, options)
		if err != nil {
			if !errors.Is(err, io.ErrClosedPipe) && !errors.Is(err, io.EOF) {
				fmt.Println("attach podman attach:", err)
			}
		}
	}()

	select {
	case <-ready:
	case <-ctx.Done():
		wg.Wait()
		return
	}

	go func() {
		defer wg.Done()
		defer cleanup()

		buf := make([]byte, 4096)
		for {
			n, err := stdoutReader.Read(buf)
			if err != nil {
				if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrClosedPipe) {
					fmt.Println("attach stdout read:", err)
				}
				return
			}

			err = conn.WriteMessage(websocket.BinaryMessage, buf[:n])
			if err != nil {
				if !websocket.IsCloseError(
					err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived,
				) {
					fmt.Println("attach websocket write:", err)
				}
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		defer cleanup()

		for {
			_, buf, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(
					err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived,
				) {
					fmt.Println("attach websocket read:", err)
				}
				return
			}

			_, err = stdinWriter.Write(buf)
			if err != nil {
				if !errors.Is(err, io.ErrClosedPipe) {
					fmt.Println("attach stdin write:", err)
				}
				return
			}
		}
	}()

	wg.Wait()
}

func (h *ServerHandler) handle_stats(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("stats connection upgrade:", err)
		return
	}
	defer conn.Close()

	options := &containers.StatsOptions{}
	options = options.WithAll(false).WithInterval(1).WithStream(true)
	sr, err := containers.Stats(*h.srv.podman, []string{id}, options)
	if err != nil {
		fmt.Println("stats podman stats:", err)
		return
	}

	for report := range sr {
		for _, stat := range report.Stats {
			var netRx uint64
			var netTx uint64

			for _, net := range stat.Network {
				netRx += net.RxBytes
				netTx += net.TxBytes
			}

			err := conn.WriteMessage(1, []byte(fmt.Sprintf("cpu %.2f mem %.2f net %d %d\n", stat.CPU, stat.MemPerc, netRx, netTx)))
			if err != nil {
				if !websocket.IsCloseError(
					err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived,
				) {
					fmt.Println("stats websocket write:", err)
				}
				return
			}
		}
	}
}

func (h *ServerHandler) handle_backup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i, err := containers.Inspect(*h.srv.podman, id, nil)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var tmpd string
	if len(i.Mounts) > 0 {
		tmpd, err = os.MkdirTemp("", "")
		if err != nil {
			fmt.Println("server handle backup create tmpdir:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer os.RemoveAll(tmpd)
	} else {
		fmt.Println("server handle backup no volume mounts found")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// create manifest file
	fmanifest, err := os.Create(filepath.Join(tmpd, "manifest.yaml"))
	if err != nil {
		fmt.Println("server handle backup manifest file:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer fmanifest.Close()

	image, _ := images.GetImage(*h.srv.podman, i.Image, nil)

	manifest := api.ServerBackupManifest{
		Server: api.ServerBackupManifestServer{
			Id:    i.ID,
			Name:  i.Name,
			Image: image.RepoTags[0], // TODO check
		},
		Created: time.Now().UTC(),
	}
	manifest.Server.Mounts = make([]api.ServerBackupManifestServerMount, len(i.Mounts))

	// export volumes
	for j, m := range i.Mounts {
		fexport, err := os.Create(filepath.Join(tmpd, fmt.Sprintf("%s.tar", m.Name))) // TODO clean up
		if err != nil {
			fmt.Println("server handle backup create volume export file:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = volumes.Export(*h.srv.podman, m.Name, fexport)
		if err != nil {
			fmt.Println("server handle backup volume export:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		manifest.Server.Mounts[j].Name = m.Name
		manifest.Server.Mounts[j].Destination = m.Destination

		fexport.Close()
	}

	// write manifest
	bmanifest, err := yaml.Marshal(manifest)
	if err != nil {
		fmt.Println("server handle backup manifest serializing:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmanifest.Write(bmanifest)

	w.Header().Set("Content-Type", "application/gzip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.tar.gz", i.Name))

	gz := gzip.NewWriter(w)
	defer gz.Close()

	tw := tar.NewWriter(gz)
	defer tw.Close()

	err = filepath.Walk(tmpd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == tmpd {
			return nil
		}

		rel, err := filepath.Rel(tmpd, path)
		if err != nil {
			return err
		}

		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = filepath.Join(i.Name, rel)

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)
		return err
	})

	if err != nil {
		fmt.Println("server handle backup tar gz:", err)
	}
}

func (h *ServerHandler) Register(m *http.ServeMux) {
	m.HandleFunc("POST /servers", h.handle_create)
	m.HandleFunc("GET /servers", h.handle_read)
	m.HandleFunc("GET /servers/{id}", h.handle_read_one)
	m.HandleFunc("DELETE /servers/{id}", h.handle_delete)

	m.HandleFunc("POST /servers/{id}/start", h.handle_start)
	m.HandleFunc("POST /servers/{id}/stop", h.handle_stop)

	m.HandleFunc("GET /servers/{id}/attach", h.handle_attach)
	m.HandleFunc("GET /servers/{id}/stats", h.handle_stats)

	m.HandleFunc("GET /servers/{id}/backup", h.handle_backup)
}

func main() {
	ctx := context.Background()

	podman_uri := "unix:/run/user/1000/podman/podman.sock"
	podman, err := bindings.NewConnection(ctx, podman_uri)
	if err != nil {
		log.Fatal(err)
	}

	srv := Server{
		podman: &podman,
	}

	mux := http.NewServeMux()

	image_handler := ImageHandler{
		srv: &srv,
	}
	image_handler.Register(mux)

	server_handler := ServerHandler{
		srv: &srv,
	}
	server_handler.Register(mux)

	handler := middleware.Log(middleware.Cors(mux))

	err = http.ListenAndServe("0.0.0.0:8000", handler)
	if err != nil {
		log.Fatal(err)
	}
}
