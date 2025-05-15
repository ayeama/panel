package agent

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ayeama/panel/api/internal/agent/stat"
	"github.com/ayeama/panel/api/internal/broker"
	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/runtime"
)

type Agent struct {
	broker  broker.Broker
	runtime runtime.Runtime
}

func New(broker broker.Broker, runtime runtime.Runtime) *Agent {
	return &Agent{broker: broker, runtime: runtime}
}

func (a *Agent) Start() {
	slog.Info("starting")

	var wg sync.WaitGroup

	a.handleStats(&wg)
	a.handleCommand(&wg)

	wg.Wait()
}

func (a *Agent) handleStats(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		hostname := stat.Hostname()
		cpuTotal, cpuIdle := stat.Cpu()

		for {
			time.Sleep(time.Second * 1)

			uptime := stat.Uptime()
			cpu := stat.CpuPercent(&cpuTotal, &cpuIdle)
			memory := stat.MemoryPercent()
			time := time.Now().UTC()

			a.broker.PublishEventAgentStat(domain.EventAgentStat{
				Hostname: hostname,
				Uptime:   uptime,
				Cpu:      cpu,
				Memory:   memory,
				Time:     time,
			})

			slog.Info(
				"published event",
				slog.String("channel", "agent:stat"),
				slog.String("type", "EventAgentStat"),
				slog.String("hostname", hostname),
				slog.Float64("uptime", uptime),
				slog.Float64("cpu", cpu),
				slog.Float64("memory", memory),
			)
		}
	}()
}

func (a *Agent) handleCommand(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		// hostname := stat.Hostname()

		for {
			event := a.broker.ReadEventAgentCommand()
			containerId := a.runtime.Create()

			fmt.Println(event, containerId)
		}
	}()
}
