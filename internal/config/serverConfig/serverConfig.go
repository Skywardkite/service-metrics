package serverconfig

import (
	"errors"
	"flag"
)

type Config struct {
	FlagRunAddr  	string
}

func ParseFlags() (Config, error){
    var cfg Config
    flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
    flag.Parse()

	if len(flag.Args()) > 0 {
		return cfg, errors.New("Unknown flags")
	}

    return cfg, nil
}