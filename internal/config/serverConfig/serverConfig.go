package serverconfig

import (
	"flag"
	"os"
)

type Config struct {
	FlagRunAddr  	string
}

func ParseFlags() (Config, error){
    var cfg Config
    flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
    flag.Parse()

	if envFlagRunAddr := os.Getenv("ADDRESS"); envFlagRunAddr != "" {
        cfg.FlagRunAddr = envFlagRunAddr
    }

    return cfg, nil
}