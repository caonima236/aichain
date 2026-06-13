package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/xiaoran/aichain/x/aichain/types"
	daotypes "github.com/xiaoran/aichain/x/agentdao/types"
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
	ValidateProposal(proposal types.GovernanceProposal) error
	IsHumanCouncilMember(addr string) bool
}

func NewKeeper(storeService store.KVStoreService, logger log.Logger, ak AgentKeeper, ck ConstitutionalKeeper) Keeper {
	return Keeper{storeService: storeService, logger: logger, agentKeeper: ak, constitutional: ck}
}

func (k Keeper) CreateProposal(ctx context.Context, proposer, content string, proposalType int32, quorum uint64) (uint64, error) {
	// Type 4 (constitutional amendment): only human council members can propose.
	// Council members skip the registered-agent check.
	if proposalType == 4 {
		if k.constitutional.IsHumanCouncilMember(proposer) {
			// Allowed — skip agent check
		} else {
			return 0, fmt.Errorf("ERR_HUMAN_COUNCIL_REQUIRED: constitutional amendments require human council member approval (code 2)")
		}
	} else {
		agent, err := k.agentKeeper.GetAgent(ctx, proposer)
		if err != nil {
			return 0, fmt.Errorf("proposer not a registered agent: %w", err)
		}
		_ = agent
	}
	store := k.storeService.OpenKVStore(ctx)
	countBytes, _ := store.Get(daotypes.ProposalCountKey)
	count := uint64(0)
	if countBytes != nil {
		count = binary.BigEndian.Uint64(countBytes)
	}
	count++
	proposal := types.GovernanceProposal{
		ProposalID: count, Proposer: proposer, Type: proposalType,
		Content: content, Quorum: quorum,
	}
	if err := k.constitutional.ValidateProposal(proposal); err != nil {
		return 0, err
	}
	bz, _ := json.Marshal(proposal)
	store.Set(daotypes.GetProposalKey(count), bz)
	countBz := make([]byte, 8)
	binary.BigEndian.PutUint64(countBz, count)
	store.Set(daotypes.ProposalCountKey, countBz)
	k.logger.Info("governance proposal created", "id", count, "proposer", proposer)
	return count, nil
}
