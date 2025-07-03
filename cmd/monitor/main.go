package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/maestroi/solana-offnode-monitor/internal/config"
	"github.com/maestroi/solana-offnode-monitor/internal/metrics"
	"github.com/maestroi/solana-offnode-monitor/internal/solana"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	solanaClient := solana.NewClient(cfg.RPCEndpoint)
	metrics.Register(cfg.Validators)

	go metrics.CollectLoop(solanaClient, cfg.Validators, cfg.PollInterval)

	http.Handle("/metrics", promhttp.Handler())
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("Serving metrics at %s/metrics", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
