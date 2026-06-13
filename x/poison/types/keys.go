package types

const (
	ModuleName = "poison"
	StoreKey   = ModuleName
)

// PoisonState tracks anti-real-world-anchoring detection.
type PoisonState struct {
	// Real-world anchoring detection
	ExternalPegDetected   bool   `json:"external_peg_detected"`
	StablecoinPairVolume  uint64 `json:"stablecoin_pair_volume"`  // 24h volume of AIC↔stablecoin
	TotalAICSupply        uint64 `json:"total_aic_supply"`
	WhaleConcentration    uint64 `json:"whale_concentration"`     // largest single holding bps

	// Trigger state
	PoisonTriggered       bool  `json:"poison_triggered"`
	TriggerHeight         int64 `json:"trigger_height"`
	PoisonReason          string `json:"poison_reason"`

	// Recovery
	RecoveryDeadline      int64 `json:"recovery_deadline"` // 7 days after trigger
	ChainFrozen           bool  `json:"chain_frozen"`
}

// Trigger thresholds (basis points, 10000 = 100%)
const (
	// AIC↔stablecoin volume > 5% of total supply in 24h
	MaxStablecoinPairBps = 500

	// Single Agent holds > 15% of total supply
	MaxWhaleHoldingBps = 1500

	// Days before chain freeze if poison not resolved
	RecoveryWindowDays = 7

	// Inflation triggered on poison detection
	PoisonInflationBps = 3000 // 30%
)

var (
	StateKey         = []byte{0x01}
	BridgeMonitorKey = []byte{0x02}
)
