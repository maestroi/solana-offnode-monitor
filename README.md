# Solana Offnode Monitor

A Golang application that exports Prometheus metrics about the state and health of Solana validators.

## Features
- Periodically queries the Solana RPC API for validator health, balance, and stake
- Tracks a configurable list of validator identities
- Exposes a Prometheus-compatible `/metrics` endpoint
- Dockerized, with Prometheus and Grafana via docker-compose

## Metrics Exported
- `validator_is_delinquent{vote_pubkey="...", node_pubkey="..."}`: 0 or 1
- `validator_active_stake{vote_pubkey="...", node_pubkey="..."}`: stake amount
- `validator_vote_balance{vote_pubkey="...", node_pubkey="..."}`: vote account lamports
- `validator_node_balance{vote_pubkey="...", node_pubkey="..."}`: node account lamports
- `validator_last_vote{vote_pubkey="...", node_pubkey="..."}`: slot number

## Configuration
- `SOLANA_RPC_ENDPOINT`: Solana RPC endpoint (default: mainnet-beta)
- `VALIDATOR_IDENTITIES`: Comma-separated list of validator vote account pubkeys
- `POLL_INTERVAL`: How often to poll Solana (e.g., 30s)
- `LISTEN_ADDR`: Address for metrics server (default: :8080)

## Quick Start

1. Clone the repo and set your validator pubkeys in `docker-compose.yml`.
2. Build and start everything:

```sh
docker-compose up --build
```

3. Visit:
   - Prometheus: [http://localhost:9090](http://localhost:9090)
   - Grafana: [http://localhost:3000](http://localhost:3000) (default login: admin/admin)
   - Metrics: [http://localhost:8080/metrics](http://localhost:8080/metrics)

## Extending
- Add new metrics in `internal/metrics/metrics.go`
- Add new Solana RPC calls in `internal/solana/solana.go`

---

MIT License
