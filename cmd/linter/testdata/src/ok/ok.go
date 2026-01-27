package main

import (
	"log"
	"os"
)

func main() {
	log.Fatal("ok") // допустимо, внутри main
	os.Exit(1)      // допустимо, внутри main
}
