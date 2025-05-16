package ep

import (
	"log/slog"
	"sync"

	"github.com/ayeama/panel/api/internal/broker"
	"github.com/ayeama/panel/api/internal/database"
	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
)

type Ep struct {
	database *database.Database
	broker   broker.Broker
}

func New(database *database.Database, broker broker.Broker) *Ep {
	return &Ep{database: database, broker: broker}
}

func (e *Ep) Start() {
	slog.Info("starting")

	var wg sync.WaitGroup

	e.handleAgentStats(&wg)

	wg.Wait()
}

func (e *Ep) handleAgentStats(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		agentRepository := repository.NewAgentRepository(e.database.Db)

		for {
			event := e.broker.ReadEventAgentStat()

			domainAgent := domain.Agent{Hostname: event.Hostname}
			agentRepository.Update(&domainAgent, &domain.AgentStat{Time: event.Time})
		}
	}()
}
