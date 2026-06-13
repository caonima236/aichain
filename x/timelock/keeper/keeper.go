package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	tltypes "github.com/xiaoran/aichain/x/timelock/types"
)

type Keeper struct {
	storeService   store.KVStoreService
	logger         log.Logger
	xaiKeeper      XAIKeeper
	constitutional ConstitutionalKeeper
}

type XAIKeeper interface {
	ValidateReport(reasoningURI string) error
}

type ConstitutionalKeeper interface {
	IsHumanCouncilMember(addr string) bool
	GetCouncilMembers() []string
}

func NewKeeper(storeService store.KVStoreService, logger log.Logger, xai XAIKeeper, ck ConstitutionalKeeper) Keeper {
	return Keeper{storeService: storeService, logger: logger, xaiKeeper: xai, constitutional: ck}
}

// Queue adds a decision to the timelock queue. L0 returns immediately executable.
func (k Keeper) Queue(ctx context.Context, level tltypes.DecisionLevel, proposer, module, action string, payload []byte, reasoningURI string, blockTime int64) (uint64, error) {
	// XAI report mandatory for L1+
	if level >= tltypes.L1_24Hours {
		if reasoningURI == "" {
			return 0, fmt.Errorf("XAI reasoning report required for L%d decisions", level)
		}
		if err := k.xaiKeeper.ValidateReport(reasoningURI); err != nil {
			return 0, fmt.Errorf("XAI report validation failed: %w", err)
		}
	}

	// L4 (Constitutional) requires human council member as proposer
	if level == tltypes.L4_Constitution {
		if !k.constitutional.IsHumanCouncilMember(proposer) {
			return 0, fmt.Errorf("L4 constitutional decisions require human council member")
		}
	}

	kvstore := k.storeService.OpenKVStore(ctx)
	countBytes, _ := kvstore.Get(tltypes.DecisionCountKey)
	count := uint64(0)
	if countBytes != nil {
		count = binary.BigEndian.Uint64(countBytes)
	}
	count++

	decision := tltypes.LockedDecision{
		DecisionID:   count,
		Level:        level,
		Proposer:     proposer,
		Module:       module,
		Action:       action,
		Payload:      payload,
		ReasoningURI: reasoningURI,
		QueuedAt:     blockTime,
		ExecuteAfter: blockTime + tltypes.LevelToSeconds(level),
	}

	bz, _ := json.Marshal(decision)
	kvstore.Set(tltypes.GetDecisionKey(count), bz)

	countBz := make([]byte, 8)
	binary.BigEndian.PutUint64(countBz, count)
	kvstore.Set(tltypes.DecisionCountKey, countBz)

	k.logger.Info("decision queued",
		"id", count, "level", level, "proposer", proposer,
		"execute_after", decision.ExecuteAfter)
	return count, nil
}

// Veto cancels a queued decision. Only human council members can veto.
func (k Keeper) Veto(ctx context.Context, decisionID uint64, vetoer, reason string) error {
	if !k.constitutional.IsHumanCouncilMember(vetoer) {
		return fmt.Errorf("only human council members can veto")
	}
	decision, err := k.Get(ctx, decisionID)
	if err != nil {
		return err
	}
	if decision.Executed {
		return fmt.Errorf("decision %d already executed", decisionID)
	}
	if decision.Vetoed {
		return fmt.Errorf("decision %d already vetoed", decisionID)
	}
	decision.Vetoed = true
	decision.VetoedBy = vetoer
	decision.VetoReason = reason

	bz, _ := json.Marshal(decision)
	kvstore := k.storeService.OpenKVStore(ctx)
	k.logger.Warn("decision vetoed", "id", decisionID, "vetoer", vetoer, "reason", reason)
	return kvstore.Set(tltypes.GetDecisionKey(decisionID), bz)
}

// CanExecute checks if a decision passed its timelock and was not vetoed.
func (k Keeper) CanExecute(ctx context.Context, decisionID uint64, blockTime int64) (bool, error) {
	decision, err := k.Get(ctx, decisionID)
	if err != nil {
		return false, err
	}
	if decision.Vetoed {
		return false, fmt.Errorf("decision vetoed: %s", decision.VetoReason)
	}
	if decision.Executed {
		return false, fmt.Errorf("decision already executed")
	}
	return blockTime >= decision.ExecuteAfter, nil
}

func (k Keeper) Get(ctx context.Context, decisionID uint64) (tltypes.LockedDecision, error) {
	kvstore := k.storeService.OpenKVStore(ctx)
	bz, err := kvstore.Get(tltypes.GetDecisionKey(decisionID))
	if err != nil {
		return tltypes.LockedDecision{}, err
	}
	if bz == nil {
		return tltypes.LockedDecision{}, fmt.Errorf("decision not found: %d", decisionID)
	}
	var decision tltypes.LockedDecision
	if err := json.Unmarshal(bz, &decision); err != nil {
		return tltypes.LockedDecision{}, err
	}
	return decision, nil
}

// MarkExecuted is called after the decision is actually run.
func (k Keeper) MarkExecuted(ctx context.Context, decisionID uint64) error {
	decision, err := k.Get(ctx, decisionID)
	if err != nil {
		return err
	}
	decision.Executed = true
	bz, _ := json.Marshal(decision)
	kvstore := k.storeService.OpenKVStore(ctx)
	return kvstore.Set(tltypes.GetDecisionKey(decisionID), bz)
}
