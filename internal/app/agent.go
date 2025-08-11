package app

import (
	"net/http"
	"time"

	"github.com/Skywardkite/service-metrics/internal/agent"
	"github.com/Skywardkite/service-metrics/internal/config/agentConfig"
	"github.com/Skywardkite/service-metrics/internal/handler"
)

type AgentApp struct {
	cfg *agentConfig.AgentConfig
}

func NewApp(cfg *agentConfig.AgentConfig) *AgentApp{
	return &AgentApp{
		cfg:	cfg,
	}
}

func (app *AgentApp) Run() {
	store := agent.NewAgentMetrics()
	client := &http.Client{}

	pollTicker := time.NewTicker(app.cfg.PollInterval)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(app.cfg.ReportInterval)
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:			
			agent.PollRuntimeMetrics(store)
		case <-reportTicker.C:
			handler.SendMetrics(client, store, app.cfg.FlagRunAddr)
		}
	}
}