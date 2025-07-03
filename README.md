# Solana Offnode Monitor

A Golang application that exports Prometheus metrics about the state and health of Solana validators, with a ready-to-use Grafana dashboard.

## Features
- Periodically queries the Solana RPC API for validator health, balance, stake, commission, and epoch
- Tracks a configurable list of validator vote accounts
- Exposes a Prometheus-compatible `/metrics` endpoint
- Exports metrics with labels for vote account, node identity, and validator name
- Dockerized, with Prometheus and Grafana via docker-compose
- Prebuilt Grafana dashboard for easy monitoring and filtering by validator name

## Metrics Exported
- `validator_is_delinquent{vote_pubkey, node_pubkey, name}`: 0 (active), 1 (delinquent), -1 (not found)
- `validator_active_stake{vote_pubkey, node_pubkey, name}`: active stake (lamports)
- `validator_vote_balance{vote_pubkey, node_pubkey, name}`: vote account balance (lamports)
- `validator_node_balance{vote_pubkey, node_pubkey, name}`: node account balance (lamports)
- `validator_last_vote{vote_pubkey, node_pubkey, name}`: last vote slot
- `validator_commission{vote_pubkey, node_pubkey, name}`: commission percentage
- `solana_epoch_number`: current Solana epoch

## Configuration
- `SOLANA_RPC_ENDPOINT`: Solana RPC endpoint (default: mainnet-beta)
- `VALIDATOR_IDENTITIES`: Comma-separated list of validator vote account pubkeys
- `POLL_INTERVAL`: How often to poll Solana (e.g., 30s)
- `LISTEN_ADDR`: Address for metrics server (default: :8080)

## Quick Start

1. Clone the repo and set your validator vote pubkeys in `docker-compose.yml`:
   ```yaml
   environment:
     - VALIDATOR_IDENTITIES=VotePubkey1,VotePubkey2
   ```
2. Build and start everything:
   ```sh
   docker-compose up --build
   ```
3. Visit:
   - Prometheus: [http://localhost:9090](http://localhost:9090)
   - Grafana: [http://localhost:3000](http://localhost:3000) (default login: admin/admin)
   - Metrics: [http://localhost:8080/metrics](http://localhost:8080/metrics)

## Grafana Dashboard
- The dashboard is auto-provisioned and shows:
  - Active stake, vote and node balances (in SOL)
  - Commission (%)
  - Delinquent and not found validator counts
  - Current epoch number
  - All metrics are filterable/searchable by validator name

## Extending
- Add new metrics in `internal/metrics/metrics.go`
- Add new Solana RPC calls in `internal/solana/solana.go`
- Customize the dashboard in `grafana/provisioning/dashboards/solana-monitor-dashboard.json`

---

MIT License
