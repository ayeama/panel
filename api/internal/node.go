package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ayeama/panel/api/internal/repository"
	"github.com/ayeama/panel/api/internal/service"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/system"
)

type NodeConnection struct {
	mu      sync.Mutex
	Name    string
	Uri     string
	context context.Context
}

func NewNodeConnection(name string, uri string) *NodeConnection {
	ctx, err := bindings.NewConnectionWithIdentity(context.Background(), uri, "/home/alex/.ssh/id_rsa", false)
	if err != nil {
		panic(err)
	}
	return &NodeConnection{Name: name, Uri: uri, context: ctx}
}

func (c *NodeConnection) Alive() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	result, err := system.Info(c.context, nil)

	cpuUsage := result.Host.CPUUtilization.SystemPercent + result.Host.CPUUtilization.UserPercent
	memUsage := (float64(result.Host.MemTotal-result.Host.MemFree) / float64(result.Host.MemTotal)) * 100

	fmt.Printf("%.2f%% %.2f%%\n", cpuUsage, memUsage)
	return err == nil
}

type NodeConnectionPool struct {
	mu    sync.Mutex
	Conns map[string]*NodeConnection
}

func NewNodeConnectionPool(ctxs map[string]*NodeConnection) *NodeConnectionPool {
	return &NodeConnectionPool{Conns: ctxs}
}

type Node struct {
	nodeService *service.NodeService
	pool        *NodeConnectionPool
}

func NewNode() *Node {
	db, err := sql.Open("sqlite3", "panel.db")
	if err != nil {
		panic(err)
	}

	connections := make(map[string]*NodeConnection, 0)
	connections["rt1"] = NewNodeConnection("rt1", "ssh://alex@192.168.114.4/run/user/1000/podman/podman.sock")
	connections["rt2"] = NewNodeConnection("rt2", "ssh://alex@192.168.114.5/run/user/1000/podman/podman.sock")

	pool := NewNodeConnectionPool(connections)

	// runtime, err := runtime.New(runtime.RuntimeTypePodman)
	// if err != nil {
	// 	panic(err)
	// }

	nodeRepository := repository.NewNodeRepository(db)
	nodeService := service.NewNodeService(nodeRepository)

	return &Node{nodeService: nodeService, pool: pool}
}

func (a *Node) Start() {
	slog.Info("starting")

	for {
		var wg sync.WaitGroup
		for _, conn := range a.pool.Conns {
			wg.Add(1)
			go func() {
				conn.Alive()
				wg.Done()
			}()
			// fmt.Println(name, conn.Alive())
		}

		wg.Wait()
		time.Sleep(time.Second * 1)
	}
}
