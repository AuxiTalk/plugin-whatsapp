package main

import (
	"context"
	"fmt"
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

	// Exemplo de QR no terminal (será movido para dentro do client depois)
	fmt.Println("[whatsapp] aguardando QR... (implementação completa em próximo passo)")

	if err := runtime.Listen(); err != nil {
		log.Fatalf("runtime: %v", err)
	}

	_ = context.Background()
}
