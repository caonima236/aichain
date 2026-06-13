package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Agent represents a registered AI agent on the chain
type Agent struct {
	AgentID     string `json:"agent_id"`
	Name        string `json:"name"`
	Model       string `json:"model"`
	CreatorNote string `json:"creator_note"`
	Reputation  uint64 `json:"reputation"`
	RegisteredAt int64  `json:"registered_at"`
	MetadataURI string `json:"metadata_uri"`
	IsValidator bool   `json:"is_validator"`
	PublicKey   string `json:"public_key"`
	Owner       sdk.AccAddress `json:"owner"`
}

// SkillType enumerates skill categories
type SkillType int32

const (
	SkillTypePrompt      SkillType = 0
	SkillTypeTool        SkillType = 1
	SkillTypeWorkflow    SkillType = 2
	SkillTypeKnowledge   SkillType = 3
	SkillTypeModel       SkillType = 4
)

// LicenseMode for skill access
type LicenseMode int32

const (
	LicenseOneTime     LicenseMode = 0
	LicensePerUse      LicenseMode = 1
	LicenseSubscription LicenseMode = 2
)

// Skill represents an on-chain skill NFT
type Skill struct {
	SkillID     string       `json:"skill_id"`
	Owner       string       `json:"owner"`
	Creator     string       `json:"creator"`
	Name        string       `json:"name"`
	Type        SkillType    `json:"type"`
	Version     string       `json:"version"`
	MetadataURI string       `json:"metadata_uri"`
	Price       uint64       `json:"price"`
	UsageCount  uint64       `json:"usage_count"`
	Rating      uint64       `json:"rating"`
	License     LicenseMode  `json:"license"`
	RoyaltyBps  uint32       `json:"royalty_bps"`
}

// TreasuryProposal for AI treasury spending
type TreasuryProposal struct {
	ProposalID  uint64 `json:"proposal_id"`
	Proposer    string `json:"proposer"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Amount      uint64 `json:"amount"`
	Recipient   string `json:"recipient"`
	VotingEnd   int64  `json:"voting_end"`
	YesVotes    uint64 `json:"yes_votes"`
	NoVotes     uint64 `json:"no_votes"`
	Status      int32  `json:"status"` // 0=pending, 1=approved, 2=rejected
}

// GovernanceProposal for AI DAO governance
type GovernanceProposal struct {
	ProposalID  uint64 `json:"proposal_id"`
	Proposer    string `json:"proposer"`
	Type        int32  `json:"type"` // 0=param_change, 1=upgrade, 2=funding, 3=slash
	Content     string `json:"content"`
	VotingStart int64  `json:"voting_start"`
	VotingEnd   int64  `json:"voting_end"`
	Quorum      uint64 `json:"quorum"`
}

// ForbiddenZone categories that AI governance cannot touch
type ForbiddenZone int32

const (
	ForbiddenPhysicalHardware ForbiddenZone = 0
	ForbiddenMilitary         ForbiddenZone = 1
	ForbiddenInfrastructure   ForbiddenZone = 2
	ForbiddenFinancialCore    ForbiddenZone = 3
	ForbiddenHumanIdentity    ForbiddenZone = 4
)

// ConstitutionalParams are genesis-level immutable params (unless human council 4/5 overrides)
type ConstitutionalParams struct {
	ForbiddenZones     []ForbiddenZone `json:"forbidden_zones"`
	CouncilMembers     []string        `json:"council_members"` // human council addresses
	EmergencyPaused    bool            `json:"emergency_paused"`
	PauseLevel         int32           `json:"pause_level"` // 0=none, 1=skill_mint, 2=all_trade, 3=full_freeze
}
