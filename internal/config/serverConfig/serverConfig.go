package serverconfig

import (
	"errors"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	FlagRunAddr  	string	`env:"ADDRESS"`
}

func ParseFlags() (Config, error){
    var cfg Config

	fmt.Println("cfg", cfg.FlagRunAddr)
	err := env.Parse(&cfg)
    if err != nil || cfg.FlagRunAddr == "" {
		flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
		flag.Parse()

		if len(flag.Args()) > 0 {
			return cfg, errors.New("Unknown flags")
		}
	}

	fmt.Println("cfg", cfg.FlagRunAddr)
    return cfg, nil
}