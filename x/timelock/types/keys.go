package types

import "encoding/binary"

const (
	ModuleName = "timelock"
	StoreKey   = ModuleName
)

// DecisionLevel determines how long a decision must wait before execution.
type DecisionLevel int32

const (
	L0_Immediate     DecisionLevel = 0 // Skill mint, Token transfer, Agent register
	L1_24Hours       DecisionLevel = 1 // Small treasury, minor params
	L2_72Hours       DecisionLevel = 2 // Large treasury, protocol params, validator changes
	L3_7Days         DecisionLevel = 3 // Protocol upgrades, new modules
	L4_Constitution  DecisionLevel = 4 // Constitutional — requires human 4/5
)

// LockedDecision is a queued decision awaiting execution after timelock.
type LockedDecision struct {
	DecisionID     uint64        `json:"decision_id"`
	Level          DecisionLevel `json:"level"`
	Proposer       string        `json:"proposer"`
	Module         string        `json:"module"`
	Action         string        `json:"action"`
	Payload        []byte        `json:"payload"`
	ReasoningURI   string        `json:"reasoning_uri"`   // IPFS hash of XAI report
	QueuedAt       int64         `json:"queued_at"`
	ExecuteAfter   int64         `json:"execute_after"`
	Vetoed         bool          `json:"vetoed"`
	VetoedBy       string        `json:"vetoed_by"`
	VetoReason     string        `json:"veto_reason"`
	Executed       bool          `json:"executed"`
}

var (
	DecisionKey      = []byte{0x01}
	DecisionCountKey = []byte{0x02}
	PendingQueueKey  = []byte{0x03}
)

func GetDecisionKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append(DecisionKey, bz...)
}

// LevelToSeconds returns the timelock duration in seconds.
func LevelToSeconds(level DecisionLevel) int64 {
	switch level {
	case L0_Immediate:
		return 0
	case L1_24Hours:
		return 24 * 60 * 60
	case L2_72Hours:
		return 72 * 60 * 60
	case L3_7Days:
		return 7 * 24 * 60 * 60
	case L4_Constitution:
		return 14 * 24 * 60 * 60 // 14 days + 4/5 council
	default:
		return 24 * 60 * 60 // safe default
	}
}
