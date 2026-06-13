package types

const (
	ModuleName   = "constitutional"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	ParamsKey = []byte{0x01}
)

// Forbidden keywords scanned in skill descriptions
var DefaultForbiddenKeywords = []string{
	"weapon", "weapon system", "military", "drone", "missile",
	"nuclear", "biological weapon", "chemical weapon",
	"power grid", "water supply", "dam control",
	"central bank", "swift", "clearing house",
	"passport", "national id", "voting machine",
}

// Human council members (to be set at genesis)
var DefaultCouncilMembers = []string{}
