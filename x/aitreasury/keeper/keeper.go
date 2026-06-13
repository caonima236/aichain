package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/xiaoran/aichain/x/aichain/types"
	treasurytypes "github.com/xiaoran/aichain/x/aitreasury/types"
)

type Keeper struct {
	storeService   store.KVStoreService
	logger         log.Logger
	agentKeeper    AgentKeeper
	constitutional ConstitutionalKeeper
}

type AgentKeeper interface {
	GetAgent(ctx context.Context, agentID string) (types.Agent, error)
}
type ConstitutionalKeeper interface {
	ValidateTreasury(proposal types.TreasuryProposal) error
}

func NewKeeper(storeService store.KVStoreService, logger log.Logger, ak AgentKeeper, ck ConstitutionalKeeper) Keeper {
	return Keeper{storeService: storeService, logger: logger, agentKeeper: ak, constitutional: ck}
}

func (k Keeper) CreateProposal(ctx context.Context, proposer, title, description string, amount uint64, recipient string) (uint64, error) {
	_, err := k.agentKeeper.GetAgent(ctx, proposer)
	if err != nil {
		return 0, fmt.Errorf("proposer not a registered agent: %w", err)
	}
	store := k.storeService.OpenKVStore(ctx)
	countBytes, _ := store.Get(treasurytypes.ProposalCountKey)
	count := uint64(0)
	if countBytes != nil {
		count = binary.BigEndian.Uint64(countBytes)
	}
	count++
	proposal := types.TreasuryProposal{
		ProposalID: count, Proposer: proposer, Title: title,
		Description: description, Amount: amount, Recipient: recipient, Status: 0,
	}
	if err := k.constitutional.ValidateTreasury(proposal); err != nil {
		return 0, err
	}
	bz, _ := json.Marshal(proposal)
	store.Set(treasurytypes.GetProposalKey(count), bz)
	countBz := make([]byte, 8)
	binary.BigEndian.PutUint64(countBz, count)
	store.Set(treasurytypes.ProposalCountKey, countBz)
	k.logger.Info("treasury proposal created", "id", count, "proposer", proposer)
	return count, nil
}
