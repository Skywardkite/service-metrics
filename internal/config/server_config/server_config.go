package server_config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	FlagRunAddr     string
	StoreInternal   time.Duration
	FileStoragePath string
	Restore         bool
	DatabaseDSN     string
	Key             string
	AuditFile       string
	AuditURL        string
	CryptoKeyPath   string
}

func ParseFlags() (Config, error) {
	var cfg Config
	var storeInternal int
	flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&storeInternal, "i", 300, "server metrics update frequency")
	flag.StringVar(&cfg.FileStoragePath, "f", "./tmp/metrics.json", "path to storage")
	flag.BoolVar(&cfg.Restore, "r", true, "need to restore")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "database connection")
	flag.StringVar(&cfg.Key, "k", "", "hash key")
	flag.StringVar(&cfg.AuditFile, "audit-file", "", "path to audit log file")
	flag.StringVar(&cfg.AuditURL, "audit-url", "", "URL to send audit logs")
	flag.StringVar(&cfg.CryptoKeyPath, "crypto-key", "", "path to public key")
	flag.Parse()

	if envFlagRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.FlagRunAddr = envFlagRunAddr
	}

	if envStoreInternal, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		num, err := strconv.Atoi(envStoreInternal)
		if err != nil {
			return cfg, fmt.Errorf("invalid STORE_INTERVAL: %s", envStoreInternal)
		}
		storeInternal = num
	}
	cfg.StoreInternal = time.Duration(storeInternal) * time.Second

	if envFileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = envFileStoragePath
	}

	if envRestore, ok := os.LookupEnv("RESTORE"); ok {
		restore, err := strconv.ParseBool(envRestore)
		if err != nil {
			return cfg, fmt.Errorf("invalid RESTORE: %s", envRestore)
		}
		cfg.Restore = restore
	}

	if envDatabaseDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
		cfg.DatabaseDSN = envDatabaseDSN
	}

	if envFlagKey, ok := os.LookupEnv("KEY"); ok {
		cfg.Key = envFlagKey
	}

	if envAuditFile, ok := os.LookupEnv("AUDIT_FILE"); ok {
		cfg.AuditFile = envAuditFile
	}

	if envAuditURL, ok := os.LookupEnv("AUDIT_URL"); ok {
		cfg.AuditURL = envAuditURL
	}

	if сryptoKeyPath, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		cfg.CryptoKeyPath = сryptoKeyPath
	}

	return cfg, nil
}
