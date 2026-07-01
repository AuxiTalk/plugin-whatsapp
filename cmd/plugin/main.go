package main

import (
	"log"
	"os"

	"github.com/auxitalk/plugin-whatsapp/internal/config"
	"github.com/auxitalk/plugin-whatsapp/internal/plugin"
)

func main() {
	cfg := config.Load()

	runtime, err := plugin.NewRuntime(os.Stdin, os.Stdout, os.Stderr, cfg)
	if err != nil {
		log.Fatalf("runtime init: %v", err)
	}

	if err := runtime.Listen(); err != nil {
		log.Fatalf("runtime: %v", err)
	}
}
