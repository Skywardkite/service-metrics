package main

import (
	"errors"
	"flag"
	"time"
)

var (
    flagRunAddr string
    reportInterval time.Duration
    pollInterval time.Duration
)

func parseFlags() error {
    var report, poll int

    flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
    flag.IntVar(&report, "r", 10, "frequency of sending metrics")
    flag.IntVar(&poll, "p", 2, "metrics polling frequency")
    flag.Parse()

	if len(flag.Args()) > 0 {
		return errors.New("Unknown flags")
	}

	reportInterval = time.Duration(report) * time.Second
	pollInterval = time.Duration(poll) * time.Second

    return nil
}