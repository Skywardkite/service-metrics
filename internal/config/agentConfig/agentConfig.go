package agentConfig

import (
	"errors"
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

type AgentConfig struct {
	FlagRunAddr  	string
	ReportInterval 	time.Duration
	PollInterval   	time.Duration
}

type Config struct {
	FlagRunAddr  	string 	`env:"ADDRESS"`
	ReportInterval 	int		`env:"REPORT_INTERVAL"`
	PollInterval   	int		`env:"POLL_INTERVAL"`
}

func ParseFlags() (AgentConfig, error){
    var cfg Config
	var agentConfig AgentConfig

	err := env.Parse(&cfg)
    if err != nil {
        flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
		flag.IntVar(&cfg.ReportInterval, "r", 10, "frequency of sending metrics")
		flag.IntVar(&cfg.PollInterval, "p", 2, "metrics polling frequency")
		flag.Parse()

		if len(flag.Args()) > 0 {
			return agentConfig, errors.New("Unknown flags")
		}
    }

	agentConfig = AgentConfig{
		FlagRunAddr:  	cfg.FlagRunAddr,
		ReportInterval:	time.Duration(cfg.ReportInterval) * time.Second,
		PollInterval: 	time.Duration(cfg.PollInterval) * time.Second,
	}

    return agentConfig, nil
}