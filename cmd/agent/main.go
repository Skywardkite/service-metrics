package main

import (
	"time"

	"github.com/Skywardkite/service-metrics/internal/agent"
	"github.com/Skywardkite/service-metrics/internal/handler"
)

const (
    pollInterval   = 2 * time.Second
    reportInterval = 10 * time.Second
)

func main() {
    store := agent.NewAgentMetrics()
    lastReport := time.Now()

    for {
        agent.PollRuntimeMetrics(store)

        if time.Since(lastReport) >= reportInterval {
            handler.SendMetrics(store)
            lastReport = time.Now()
        }

        time.Sleep(pollInterval)
    }
}