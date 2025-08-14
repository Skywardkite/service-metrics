package serverconfig

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	FlagRunAddr  	string
    StoreInternal 	time.Duration
    FileStoragePath string
    Restore         bool
}

func ParseFlags() (Config, error){
    var cfg Config
    var storeInternal int
    flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
    flag.IntVar(&storeInternal, "i", 300, "server metrics update frequency")
    flag.StringVar(&cfg.FileStoragePath, "f", "./tmp/metrics.json", "path to storage")
    flag.BoolVar(&cfg.Restore, "r", true, "need to restore")
    flag.Parse()

	if envFlagRunAddr := os.Getenv("ADDRESS"); envFlagRunAddr != "" {
        cfg.FlagRunAddr = envFlagRunAddr
    }

    if envStoreInternal := os.Getenv("STORE_INTERVAL"); envStoreInternal != "" {
		num, err := strconv.Atoi(envStoreInternal)
		if err != nil {
			return cfg, fmt.Errorf("invalid STORE_INTERVAL: %s", envStoreInternal)
		}
        storeInternal = num
    }
    cfg.StoreInternal = time.Duration(storeInternal) * time.Second

    if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
        cfg.FileStoragePath = envFileStoragePath
    }

    if envRestore := os.Getenv("RESTORE"); envRestore != "" {
        restore, err1 := strconv.ParseBool(envRestore)
        if err1 != nil {
            return cfg, fmt.Errorf("invalid RESTORE: %s", envRestore)
        }
        cfg.Restore = restore
    }

    return cfg, nil
}