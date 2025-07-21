package main

import (
	"errors"
	"flag"
)

var flagRunAddr string

func parseFlagsServer() error {
    flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
    flag.Parse()

    if len(flag.Args()) > 0 {
		return errors.New("Unknown flags")
	}

    return nil
}
