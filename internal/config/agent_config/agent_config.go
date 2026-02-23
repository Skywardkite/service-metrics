package agentConfig

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type AgentConfig struct {
	FlagRunAddr    string
	ReportInterval time.Duration
	PollInterval   time.Duration
	UseBatch       bool
	Key            string
	RateLimit      int
	CryptoKeyPath  string
}

func ParseFlags() (AgentConfig, error) {
	var cfg AgentConfig
	var report, poll int

	flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&report, "r", 10, "frequency of sending metrics")
	flag.IntVar(&poll, "p", 2, "metrics polling frequency")
	flag.BoolVar(&cfg.UseBatch, "b", false, "use batch API")
	flag.StringVar(&cfg.Key, "k", "", "hash key")
	flag.IntVar(&cfg.RateLimit, "l", 1, "rate limit (max parallel requests)")
	flag.StringVar(&cfg.CryptoKeyPath, "crypto-key", "", "path to public key")
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

	if envFlagKey, ok := os.LookupEnv("KEY"); ok {
		cfg.Key = envFlagKey
	}

	if rateLimit, ok := os.LookupEnv("RATE_LIMIT"); ok {
		num, err := strconv.Atoi(rateLimit)
		if err != nil {
			return cfg, fmt.Errorf("invalid RATE_LIMIT: %s", rateLimit)
		}
		cfg.RateLimit = num
	}

	if сryptoKeyPath, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		cfg.CryptoKeyPath = сryptoKeyPath
	}

	return cfg, nil
}
