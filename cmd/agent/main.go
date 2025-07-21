package main

import (
	"time"

	"github.com/Skywardkite/service-metrics/internal/agent"
	"github.com/Skywardkite/service-metrics/internal/handler"
)

func main() {
    if err := parseFlags(); err != nil {
        panic(err)
    }
    
    store := agent.NewAgentMetrics()
    lastReport := time.Now()

    for {
        agent.PollRuntimeMetrics(store)

        if time.Since(lastReport) >= (reportInterval) {
            handler.SendMetrics(store, flagRunAddr)
            lastReport = time.Now()
        }

        time.Sleep(pollInterval)
    }
}