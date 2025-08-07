package serverconfig

import (
	"errors"
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	FlagRunAddr  	string	`env:"ADDRESS"`
}

func ParseFlags() (Config, error){
    var cfg Config

	err := env.Parse(&cfg)
    if err != nil || cfg.FlagRunAddr == "" {
		flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
		flag.Parse()

		if len(flag.Args()) > 0 {
			return cfg, errors.New("Unknown flags")
		}
	}
    return cfg, nil
}