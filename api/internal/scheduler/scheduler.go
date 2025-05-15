package scheduler

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ayeama/panel/api/internal/broker"
	"github.com/ayeama/panel/api/internal/domain"
)

type Scheduler struct {
	broker broker.Broker

	m      sync.RWMutex
	agents map[string]*domain.AgentStat
}

func New(broker broker.Broker) *Scheduler {
	return &Scheduler{broker: broker}
}

func (s *Scheduler) Start() {
	slog.Info("starting")

	// TODO
	// filter out agents (too high cpu usage etc)
	// optionally filter usage for scoring (rolling average etc)
	// then score agents that passed filtering (strategy pattern with option custom weights)

	var wg sync.WaitGroup

	s.agents = make(map[string]*domain.AgentStat) // TODO is this the right spot etc?

	s.handleAgentStat(&wg)
	s.handleEventSchedule(&wg)

	wg.Wait()
}

func (s *Scheduler) handleAgentStat(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			event := s.broker.ReadEventAgentStat()
			stat := &domain.AgentStat{
				Uptime: event.Uptime,
				Cpu:    event.Cpu,
				Memory: event.Memory,
				Time:   event.Time,
			}

			s.m.Lock()
			s.agents[event.Hostname] = stat
			s.m.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		for {
			time.Sleep(time.Second * 1)

			s.m.RLock()
			for _, agent := range s.agents {
				fmt.Println(agent.Cpu, agent.Online())
			}
			s.m.RUnlock()
		}
	}()
}

func (s *Scheduler) handleEventSchedule(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			time.Sleep(time.Second * 1)
		}
	}()
}
