package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := LoadConfig()
	server, err := BuildServer(cfg)
	if err != nil {
		log.Fatalf("build server: %v", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	log.Printf("WindTurbineCtl console listening on %s with data dir %s", cfg.Addr, cfg.DataDir)
	if err := server.StartWithContext(ctx); err != nil {
		log.Fatalf("serve: %v", err)
	}
	log.Printf("WindTurbineCtl stopped")
}
