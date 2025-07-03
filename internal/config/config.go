package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	RPCEndpoint  string
	Validators   []string
	PollInterval time.Duration
}

func Load() (*Config, error) {
	rpc := os.Getenv("SOLANA_RPC_ENDPOINT")
	if rpc == "" {
		rpc = "https://api.mainnet-beta.solana.com"
	}
	validators := os.Getenv("VALIDATOR_IDENTITIES")
	var validatorList []string
	if validators != "" {
		validatorList = strings.Split(validators, ",")
	}
	interval := os.Getenv("POLL_INTERVAL")
	if interval == "" {
		interval = "30s"
	}
	pollInterval, err := time.ParseDuration(interval)
	if err != nil {
		pollInterval = 30 * time.Second
	}
	return &Config{
		RPCEndpoint:  rpc,
		Validators:   validatorList,
		PollInterval: pollInterval,
	}, nil
}
