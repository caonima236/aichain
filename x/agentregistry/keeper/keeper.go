package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xiaoran/aichain/x/aichain/types"
	agenttypes "github.com/xiaoran/aichain/x/agentregistry/types"
)

type Keeper struct {
	storeService store.KVStoreService
	logger       log.Logger
}

func NewKeeper(storeService store.KVStoreService, logger log.Logger) Keeper {
	return Keeper{storeService: storeService, logger: logger}
}

func (k Keeper) RegisterAgent(ctx context.Context, name, model, creatorNote, metadataURI, publicKey string, owner sdk.AccAddress) (string, error) {
	store := k.storeService.OpenKVStore(ctx)
	countBytes, err := store.Get(agenttypes.AgentCountKey)
	if err != nil {
		return "", err
	}
	count := uint64(0)
	if countBytes != nil {
		count = binary.BigEndian.Uint64(countBytes)
	}
	count++
	agentID := fmt.Sprintf("agent-%d", count)

	exists, _ := store.Get(agenttypes.GetAgentKey(agentID))
	if exists != nil {
		return "", fmt.Errorf("agent ID collision: %s", agentID)
	}

	agent := types.Agent{
		AgentID:      agentID,
		Name:         name,
		Model:        model,
		CreatorNote:  creatorNote,
		Reputation:   0,
		RegisteredAt: 0, // set in handler with block time
		MetadataURI:  metadataURI,
		IsValidator:  false,
		PublicKey:    publicKey,
		Owner:        owner,
	}

	bz, err := json.Marshal(agent)
	if err != nil {
		return "", err
	}
	if err := store.Set(agenttypes.GetAgentKey(agentID), bz); err != nil {
		return "", err
	}

	countBz := make([]byte, 8)
	binary.BigEndian.PutUint64(countBz, count)
	if err := store.Set(agenttypes.AgentCountKey, countBz); err != nil {
		return "", err
	}

	k.logger.Info("agent registered", "agent_id", agentID, "name", name)
	return agentID, nil
}

func (k Keeper) GetAgent(ctx context.Context, agentID string) (types.Agent, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(agenttypes.GetAgentKey(agentID))
	if err != nil {
		return types.Agent{}, err
	}
	if bz == nil {
		return types.Agent{}, fmt.Errorf("agent not found: %s", agentID)
	}
	var agent types.Agent
	if err := json.Unmarshal(bz, &agent); err != nil {
		return types.Agent{}, err
	}
	return agent, nil
}

func (k Keeper) UpdateReputation(ctx context.Context, agentID string, delta int64) error {
	agent, err := k.GetAgent(ctx, agentID)
	if err != nil {
		return err
	}
	newRep := int64(agent.Reputation) + delta
	if newRep < 0 {
		newRep = 0
	}
	agent.Reputation = uint64(newRep)
	return k.setAgent(ctx, agent)
}

func (k Keeper) SetValidator(ctx context.Context, agentID string, isValidator bool) error {
	agent, err := k.GetAgent(ctx, agentID)
	if err != nil {
		return err
	}
	agent.IsValidator = isValidator
	return k.setAgent(ctx, agent)
}

func (k Keeper) setAgent(ctx context.Context, agent types.Agent) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(agent)
	if err != nil {
		return err
	}
	return store.Set(agenttypes.GetAgentKey(agent.AgentID), bz)
}

func (k Keeper) GetAgentCount(ctx context.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	bz, _ := store.Get(agenttypes.AgentCountKey)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}
