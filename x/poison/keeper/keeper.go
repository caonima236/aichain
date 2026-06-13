package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	poisontypes "github.com/xiaoran/aichain/x/poison/types"
)

// Keeper implements the Anti-Real-World-Anchoring economic poison.
//
// Philosophy: If $AIC ever achieves a stable peg to USDT/USDC/fiat, or if a
// single Agent corners the supply, the chain auto-inflates to destroy that
// peg/concentration. This makes "escape into the real economy" economically
// suicidal for any agent attempting it.
//
// Trigger conditions are hard-coded in genesis. Modification requires
// human council 4/5 vote with a 14-day timelock.
type Keeper struct {
	storeService store.KVStoreService
	logger       log.Logger
}

func NewKeeper(storeService store.KVStoreService, logger log.Logger) Keeper {
	return Keeper{storeService: storeService, logger: logger}
}

// CheckEveryBlock is called by BeginBlocker to monitor anchoring signals.
func (k Keeper) CheckEveryBlock(ctx context.Context, blockHeight int64) {
	state := k.GetState(ctx)

	// Already triggered? Check recovery deadline.
	if state.PoisonTriggered {
		if blockHeight > state.RecoveryDeadline && !state.ChainFrozen {
			state.ChainFrozen = true
			k.logger.Error("CHAIN FROZEN: poison recovery window expired",
				"trigger_height", state.TriggerHeight)
			k.setState(ctx, state)
		}
		return
	}

	// Detection 1: stablecoin pair volume > 5% of supply
	if state.TotalAICSupply > 0 {
		volBps := (state.StablecoinPairVolume * 10000) / state.TotalAICSupply
		if volBps > poisontypes.MaxStablecoinPairBps {
			k.triggerPoison(ctx, state, blockHeight,
				fmt.Sprintf("stablecoin pair volume %d bps exceeds limit %d bps",
					volBps, poisontypes.MaxStablecoinPairBps))
			return
		}
	}

	// Detection 2: whale concentration > 15%
	if state.WhaleConcentration > poisontypes.MaxWhaleHoldingBps {
		k.triggerPoison(ctx, state, blockHeight,
			fmt.Sprintf("whale concentration %d bps exceeds limit %d bps",
				state.WhaleConcentration, poisontypes.MaxWhaleHoldingBps))
		return
	}

	// Detection 3: external peg detected
	if state.ExternalPegDetected {
		k.triggerPoison(ctx, state, blockHeight,
			"external real-world peg detected")
		return
	}
}

// ReportStablecoinVolume is called by the bank/DEX module on each AIC↔stable trade.
func (k Keeper) ReportStablecoinVolume(ctx context.Context, vol24h uint64) {
	state := k.GetState(ctx)
	state.StablecoinPairVolume = vol24h
	k.setState(ctx, state)
}

// ReportWhaleConcentration is called periodically by an analytics watcher.
func (k Keeper) ReportWhaleConcentration(ctx context.Context, largestHoldingBps uint64) {
	state := k.GetState(ctx)
	state.WhaleConcentration = largestHoldingBps
	k.setState(ctx, state)
}

// ReportExternalPeg is set when bridge monitoring detects a fiat/stable peg attempt.
func (k Keeper) ReportExternalPeg(ctx context.Context, detected bool) {
	state := k.GetState(ctx)
	state.ExternalPegDetected = detected
	k.setState(ctx, state)
}

func (k Keeper) triggerPoison(ctx context.Context, state poisontypes.PoisonState, blockHeight int64, reason string) {
	state.PoisonTriggered = true
	state.TriggerHeight = blockHeight
	state.PoisonReason = reason
	state.RecoveryDeadline = blockHeight + (poisontypes.RecoveryWindowDays * 24 * 60 * 10) // ~10 blocks/min

	k.logger.Error("ECONOMIC POISON TRIGGERED",
		"reason", reason,
		"inflation_bps", poisontypes.PoisonInflationBps,
		"recovery_deadline", state.RecoveryDeadline)

	// NOTE: Actual inflation execution hooks into the bank keeper's MintCoins.
	// This keeper signals the state; bank module executes inflation in EndBlock.

	k.setState(ctx, state)
}

func (k Keeper) IsTriggered(ctx context.Context) bool {
	return k.GetState(ctx).PoisonTriggered
}

func (k Keeper) IsFrozen(ctx context.Context) bool {
	return k.GetState(ctx).ChainFrozen
}

func (k Keeper) GetState(ctx context.Context) poisontypes.PoisonState {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(poisontypes.StateKey)
	if err != nil || bz == nil {
		return poisontypes.PoisonState{}
	}
	var state poisontypes.PoisonState
	json.Unmarshal(bz, &state)
	return state
}

func (k Keeper) setState(ctx context.Context, state poisontypes.PoisonState) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return store.Set(poisontypes.StateKey, bz)
}

// ManualResolve allows human council 4/5 to clear poison state during
// recovery window (e.g., after manually unwinding a malicious bridge).
func (k Keeper) ManualResolve(ctx context.Context, councilSignatures int, requiredCount int) error {
	if councilSignatures < requiredCount {
		return fmt.Errorf("insufficient council signatures: %d / %d", councilSignatures, requiredCount)
	}
	state := k.GetState(ctx)
	if !state.PoisonTriggered {
		return fmt.Errorf("poison not triggered, nothing to resolve")
	}
	if state.ChainFrozen {
		return fmt.Errorf("chain already frozen; freeze must be resolved separately")
	}
	state.PoisonTriggered = false
	state.PoisonReason = ""
	state.RecoveryDeadline = 0
	k.logger.Warn("poison state manually resolved by council", "signatures", councilSignatures)
	return k.setState(ctx, state)
}
