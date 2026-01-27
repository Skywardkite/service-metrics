package somepkg

import (
	"log"
	"os"
)

func errCheckPanic() {
	panic("oops") // want "usage of panic is forbidden"
}

func errCheckLogFatal() {
	log.Fatal("oops") // want "call to log.Fatal is only allowed in main"
}

func errCheckOSExit() {
	os.Exit(1) // want "call to os.Exit is only allowed in main"
}
