package types

const (
	ModuleName   = "skillnft"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	SkillKey      = []byte{0x01}
	SkillCountKey = []byte{0x02}
)

func GetSkillKey(skillID string) []byte {
	return append(SkillKey, []byte(skillID)...)
}
