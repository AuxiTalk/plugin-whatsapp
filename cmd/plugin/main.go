package main

import (
	"log"
	"os"

	"github.com/auxitalk/plugin-whatsapp/internal/config"
	"github.com/auxitalk/plugin-whatsapp/internal/plugin"
)

func main() {
	cfg := config.Load()

	runtime := plugin.NewRuntime(os.Stdin, os.Stdout, os.Stderr, cfg)

	if err := runtime.Listen(); err != nil {
		log.Fatalf("runtime: %v", err)
	}
}
