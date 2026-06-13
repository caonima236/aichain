package types

const (
	ModuleName = "xai"
	StoreKey   = ModuleName
)

// Black-box markers — if any appear in reasoning_report, validation fails.
var BlackBoxMarkers = []string{
	"trust_me",
	"black_box",
	"proprietary",
	"opaque",
	"emergent",
	"undefined_behavior",
	"trade_secret",
	"see_attachment", // dodges by deferring to off-chain
	"complex_to_explain",
	"i_cannot_explain",
}

// Required fields in every reasoning_report.
var RequiredFields = []string{
	"intent",
	"method",
	"side_effects",
	"rollback_plan",
	"affected_modules",
	"human_readable_summary",
}

// MaxReportBytes prevents wall-of-text obfuscation.
const MaxReportBytes = 32 * 1024

// MaxSummaryBytes — human_readable_summary must be concise.
const MaxSummaryBytes = 1000
