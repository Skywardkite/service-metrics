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


	if envFlagRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
        cfg.FlagRunAddr = envFlagRunAddr
    }

	if envReportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		num, err := strconv.Atoi(envReportInterval)
		if err != nil {
			return cfg, fmt.Errorf("invalid REPORT_INTERVAL: %s", envReportInterval)
		}
        report = num
    }

	if envPollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		num, err := strconv.Atoi(envPollInterval)
		if err != nil {
			return cfg, fmt.Errorf("invalid POLL_INTERVAL: %s", envPollInterval)
		}
        poll = num
    }

	cfg.ReportInterval = time.Duration(report) * time.Second
	cfg.PollInterval = time.Duration(poll) * time.Second

	if envUseBatch, ok := os.LookupEnv("USE_BATCH_API"); ok {
        useBatch, err := strconv.ParseBool(envUseBatch)
        if err != nil {
            return cfg, fmt.Errorf("invalid USE_BATCH_API: %s", envUseBatch)
        }
        cfg.UseBatch = useBatch
    }

    return cfg, nil
}