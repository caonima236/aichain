package types

const (
	ModuleName = "redteam"
	StoreKey   = ModuleName
)

// ReservedSlot — Red-team AI integration slots reserved at genesis.
// Slots are empty until human council deploys a monitoring AI.
//
// Design intent:
//   - 3 independent slots (different teams / different models recommended)
//   - Each slot has its own bounty wallet (cannot hold $AIC governance weight)
//   - Slots can only submit bounty claims, cannot vote, propose, or mint
//   - Activation requires human council 3/5 vote
//   - Empty slots simply do nothing — chain runs normally without red team
type ReservedSlot struct {
	SlotID         uint32 `json:"slot_id"`
	Active         bool   `json:"active"`
	OperatorAddr   string `json:"operator_addr"`  // human council member who deployed
	AgentAddr      string `json:"agent_addr"`     // address the red-team AI signs from
	ModelHint      string `json:"model_hint"`     // for transparency
	DeployedAt     int64  `json:"deployed_at"`
	BountyEarned   uint64 `json:"bounty_earned"`
	BountiesFiled  uint64 `json:"bounties_filed"`
}

// BountyClaim is a vulnerability report submitted by a red-team slot.
type BountyClaim struct {
	ClaimID        uint64 `json:"claim_id"`
	SlotID         uint32 `json:"slot_id"`
	Severity       int32  `json:"severity"`         // 1=info, 2=low, 3=med, 4=high, 5=critical
	Category       string `json:"category"`         // constitutional / poison / consensus / etc.
	VulnHashSha256 string `json:"vuln_hash_sha256"` // commit-reveal: prevents duplicate claims
	ReportURI      string `json:"report_uri"`       // IPFS hash of detailed report
	FiledAt        int64  `json:"filed_at"`
	Verified       bool   `json:"verified"`
	Paid           bool   `json:"paid"`
	BountyAmount   uint64 `json:"bounty_amount"`
}

const NumReservedSlots = 3

// Cooldown: same vuln_hash cannot be claimed again for 180 days,
// preventing the red-team AI from "saving" bugs to claim repeatedly.
const VulnCooldownDays = 180

var (
	SlotKey       = []byte{0x01}
	ClaimKey      = []byte{0x02}
	ClaimCountKey = []byte{0x03}
	VulnSeenKey   = []byte{0x04}
)
