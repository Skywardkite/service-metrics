package main

import (
	"log"

	"github.com/Skywardkite/service-metrics/internal/app"
	agentConfig "github.com/Skywardkite/service-metrics/internal/config/agent_config"
)

func main() {
    cfg, err := agentConfig.ParseFlags()
    if err != nil {
        log.Fatal("Error to parse flags:", err)
    }
    
    a := app.NewApp(&cfg)
    a.Run()
}