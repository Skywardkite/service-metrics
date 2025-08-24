package agentConfig

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type AgentConfig struct {
	FlagRunAddr  	string
	ReportInterval 	time.Duration
	PollInterval   	time.Duration
	UseBatch 		bool
}

func ParseFlags() (AgentConfig, error){
    var cfg AgentConfig
	var report, poll int

	flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
    flag.IntVar(&report, "r", 10, "frequency of sending metrics")
    flag.IntVar(&poll, "p", 2, "metrics polling frequency")	
	flag.BoolVar(&cfg.UseBatch, "b", false, "use batch API")
    flag.Parse()


	if envFlagRunAddr := os.Getenv("ADDRESS"); envFlagRunAddr != "" {
        cfg.FlagRunAddr = envFlagRunAddr
    }

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		num, err := strconv.Atoi(envReportInterval)
		if err != nil {
			return cfg, fmt.Errorf("invalid REPORT_INTERVAL: %s", envReportInterval)
		}
        report = num
    }

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		num, err := strconv.Atoi(envPollInterval)
		if err != nil {
			return cfg, fmt.Errorf("invalid POLL_INTERVAL: %s", envPollInterval)
		}
        poll = num
    }

	cfg.ReportInterval = time.Duration(report) * time.Second
	cfg.PollInterval = time.Duration(poll) * time.Second

	if envUseBatch := os.Getenv("USE_BATCH_API"); envUseBatch != "" {
        useBatch, err := strconv.ParseBool(envUseBatch)
        if err != nil {
            return cfg, fmt.Errorf("invalid USE_BATCH_API: %s", envUseBatch)
        }
        cfg.UseBatch = useBatch
    }

    return cfg, nil
}