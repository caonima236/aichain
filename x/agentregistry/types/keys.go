package types

const (
	ModuleName   = "agentregistry"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	// AgentKey prefix for agent store
	AgentKey = []byte{0x01}
	// AgentCountKey for agent ID counter
	AgentCountKey = []byte{0x02}
)

func GetAgentKey(agentID string) []byte {
	return append(AgentKey, []byte(agentID)...)
}
