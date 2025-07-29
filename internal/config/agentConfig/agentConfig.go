package agentConfig

import (
	"errors"
	"flag"
	"time"
)

type AgentConfig struct {
	FlagRunAddr  	string
	ReportInterval 	time.Duration
	PollInterval   	time.Duration
}

func ParseFlags() (AgentConfig, error){
    var cfg AgentConfig
    var report, poll int

    flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
    flag.IntVar(&report, "r", 10, "frequency of sending metrics")
    flag.IntVar(&poll, "p", 2, "metrics polling frequency")
    flag.Parse()

	if len(flag.Args()) > 0 {
		return cfg, errors.New("Unknown flags")
	}

	cfg.ReportInterval = time.Duration(report) * time.Second
	cfg.PollInterval = time.Duration(poll) * time.Second

    return cfg, nil
}