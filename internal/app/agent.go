package app

import (
	"net/http"
	"strings"
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
			url := app.cfg.FlagRunAddr
			if !strings.HasPrefix(app.cfg.FlagRunAddr, "http://") && !strings.HasPrefix(app.cfg.FlagRunAddr, "https://") {
				url = "http://" + app.cfg.FlagRunAddr
			}
			handler.SendMetrics(client, store, url + "/update")
		}
	}
}