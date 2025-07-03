package metrics

import (
	"log"
	"time"

	"github.com/maestroi/solana-offnode-monitor/internal/solana"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	validatorIsDelinquent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "validator_is_delinquent",
			Help: "Whether the validator is delinquent (1) or not (0)",
		},
		[]string{"vote_pubkey", "node_pubkey", "name"},
	)
	validatorActiveStake = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "validator_active_stake",
			Help: "Active stake of the validator",
		},
		[]string{"vote_pubkey", "node_pubkey", "name"},
	)
	validatorVoteBalance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "validator_vote_balance",
			Help: "Balance of the validator's vote account (lamports)",
		},
		[]string{"vote_pubkey", "node_pubkey", "name"},
	)
	validatorNodeBalance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "validator_node_balance",
			Help: "Balance of the validator's node account (lamports)",
		},
		[]string{"vote_pubkey", "node_pubkey", "name"},
	)
	validatorLastVote = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "validator_last_vote",
			Help: "Last vote slot of the validator",
		},
		[]string{"vote_pubkey", "node_pubkey", "name"},
	)
	validatorCommission = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "validator_commission",
			Help: "Commission percentage of the validator",
		},
		[]string{"vote_pubkey", "node_pubkey", "name"},
	)
	epochNumber = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solana_epoch_number",
			Help: "Current Solana epoch number",
		},
	)
)

func Register(validators []string) {
	prometheus.MustRegister(validatorIsDelinquent)
	prometheus.MustRegister(validatorActiveStake)
	prometheus.MustRegister(validatorVoteBalance)
	prometheus.MustRegister(validatorNodeBalance)
	prometheus.MustRegister(validatorLastVote)
	prometheus.MustRegister(validatorCommission)
	prometheus.MustRegister(epochNumber)
}

// Note: vote_pubkey = vote account, node_pubkey = validator identity
// validator_is_delinquent: 0 = not delinquent, 1 = delinquent, -1 = not found
func CollectLoop(client *solana.Client, validators []string, interval time.Duration) {
	for {
		log.Printf("Fetching validator data from Solana RPC...")
		// Fetch and export epoch number
		epoch, err := client.GetEpochInfo()
		if err != nil {
			log.Printf("error getting epoch info: %v", err)
		} else {
			epochNumber.Set(float64(epoch))
		}
		// Fetch validator info
		validatorInfo, err := client.GetValidatorInfo()
		if err != nil {
			log.Printf("error getting validator info: %v", err)
		}
		voteAccounts, err := client.GetVoteAccounts()
		if err != nil {
			log.Printf("error getting vote accounts: %v", err)
			continue
		}
		voteList, _ := voteAccounts["current"].([]interface{})
		delinquentList, _ := voteAccounts["delinquent"].([]interface{})
		voteMap := map[string]map[string]interface{}{}
		for _, v := range voteList {
			m, _ := v.(map[string]interface{})
			voteMap[m["votePubkey"].(string)] = m
		}
		delinquentMap := map[string]map[string]interface{}{}
		for _, v := range delinquentList {
			m, _ := v.(map[string]interface{})
			delinquentMap[m["votePubkey"].(string)] = m
		}
		for _, id := range validators {
			var m map[string]interface{}
			var delinquent bool
			if mm, ok := delinquentMap[id]; ok {
				m = mm
				delinquent = true
			} else if mm, ok := voteMap[id]; ok {
				m = mm
				delinquent = false
			} else {
				log.Printf("Warning: vote account %s not found in current or delinquent lists", id)
				validatorIsDelinquent.WithLabelValues(id, "", "").Set(-1)
				validatorActiveStake.WithLabelValues(id, "", "").Set(0)
				validatorLastVote.WithLabelValues(id, "", "").Set(0)
				validatorVoteBalance.WithLabelValues(id, "", "").Set(0)
				validatorNodeBalance.WithLabelValues(id, "", "").Set(0)
				continue
			}
			votePubkey := m["votePubkey"].(string)
			nodePubkey := ""
			if np, ok := m["nodePubkey"].(string); ok {
				nodePubkey = np
			}
			name := ""
			if info, ok := validatorInfo[nodePubkey]; ok {
				name = info.Name
			}
			if delinquent {
				validatorIsDelinquent.WithLabelValues(votePubkey, nodePubkey, name).Set(1)
			} else {
				validatorIsDelinquent.WithLabelValues(votePubkey, nodePubkey, name).Set(0)
			}
			validatorActiveStake.WithLabelValues(votePubkey, nodePubkey, name).Set(m["activatedStake"].(float64))
			validatorLastVote.WithLabelValues(votePubkey, nodePubkey, name).Set(m["lastVote"].(float64))
			if commission, ok := m["commission"].(float64); ok {
				validatorCommission.WithLabelValues(votePubkey, nodePubkey, name).Set(commission)
			}
			voteBalance, err := client.GetBalance(votePubkey)
			if err != nil {
				log.Printf("error getting vote account balance for %s: %v", votePubkey, err)
				validatorVoteBalance.WithLabelValues(votePubkey, nodePubkey, name).Set(0)
			} else {
				validatorVoteBalance.WithLabelValues(votePubkey, nodePubkey, name).Set(float64(voteBalance))
			}
			if nodePubkey != "" {
				nodeBalance, err := client.GetBalance(nodePubkey)
				if err != nil {
					log.Printf("error getting node account balance for %s: %v", nodePubkey, err)
					validatorNodeBalance.WithLabelValues(votePubkey, nodePubkey, name).Set(0)
				} else {
					validatorNodeBalance.WithLabelValues(votePubkey, nodePubkey, name).Set(float64(nodeBalance))
				}
			}
		}
		time.Sleep(interval)
	}
}
