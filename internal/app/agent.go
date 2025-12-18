package app

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/Skywardkite/service-metrics/internal/agent"
	agentConfig "github.com/Skywardkite/service-metrics/internal/config/agent_config"
	"github.com/Skywardkite/service-metrics/internal/handler"
)

type AgentApp struct {
	cfg *agentConfig.AgentConfig
}

func NewApp(cfg *agentConfig.AgentConfig) *AgentApp {
	return &AgentApp{
		cfg: cfg,
	}
}

func (app *AgentApp) Run() {
	store := agent.NewAgentMetrics()
	client := handler.NewRetryableClient()

	url := app.cfg.FlagRunAddr
	if !strings.HasPrefix(app.cfg.FlagRunAddr, "http://") && !strings.HasPrefix(app.cfg.FlagRunAddr, "https://") {
		url = "http://" + app.cfg.FlagRunAddr
	}

	jobs := make(chan struct{}, 10)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		pollTicker := time.NewTicker(app.cfg.PollInterval)
		defer pollTicker.Stop()

		for range pollTicker.C {
			agent.PollRuntimeMetrics(store)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		pollTicker := time.NewTicker(app.cfg.PollInterval)
		defer pollTicker.Stop()

		for range pollTicker.C {
			agent.PollSystemMetrics(store)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		reportTicker := time.NewTicker(app.cfg.ReportInterval)
		defer reportTicker.Stop()

		for range reportTicker.C {
			jobs <- struct{}{}
		}
	}()

	app.worker(jobs, store, client, &wg, url)

	wg.Wait()
	close(jobs)
}

func (app *AgentApp) worker(jobs <-chan struct{}, store *agent.AgentMetrics, client *retryablehttp.Client, wg *sync.WaitGroup, url string) {
	for i := 0; i < app.cfg.RateLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for range jobs {
				if app.cfg.UseBatch {
					// Батчевая отправка
					err := handler.SendBatch(client, store, url, app.cfg.Key)
					if err != nil {
						log.Printf("Batch API failed, falling back to individual: %v", err)
					}
				} else {
					handler.SendMetrics(client, store, url+"/update/", app.cfg.Key)
				}
			}
		}()
	}
}
