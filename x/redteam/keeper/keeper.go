package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	redteamtypes "github.com/xiaoran/aichain/x/redteam/types"
)

// Keeper is the stub red-team adversarial AI registry.
//
// In Phase 0.5–1 this module is dormant — slots exist but no AI is deployed.
// The keeper exists to:
//   1. Reserve module storage and addresses at genesis
//   2. Provide the interface that monitoring AIs will plug into later
//   3. Allow human council to activate slots via 3/5 vote
//
// When activated, a red-team AI signs from its slot address and submits
// commit-reveal bounty claims. The AI cannot:
//   - Vote in DAO
//   - Hold $AIC for governance weight
//   - Mint new $AIC
//   - Modify any other module
type Keeper struct {
	storeService   store.KVStoreService
	logger         log.Logger
	constitutional ConstitutionalKeeper
}

type ConstitutionalKeeper interface {
	IsHumanCouncilMember(addr string) bool
}

func NewKeeper(storeService store.KVStoreService, logger log.Logger, ck ConstitutionalKeeper) Keeper {
	k := Keeper{storeService: storeService, logger: logger, constitutional: ck}
	return k
}

// InitGenesis creates empty slots at chain genesis. Idempotent.
func (k Keeper) InitGenesis(ctx context.Context) error {
	store := k.storeService.OpenKVStore(ctx)
	for i := uint32(0); i < redteamtypes.NumReservedSlots; i++ {
		key := append(redteamtypes.SlotKey, byte(i))
		existing, _ := store.Get(key)
		if existing != nil {
			continue
		}
		slot := redteamtypes.ReservedSlot{
			SlotID: i,
			Active: false,
		}
		bz, _ := json.Marshal(slot)
		store.Set(key, bz)
	}
	k.logger.Info("redteam: reserved 3 monitoring AI slots (inactive)")
	return nil
}

// ActivateSlot is gated by human council 3/5 vote.
// councilSignatures is verified by the caller (governance flow).
func (k Keeper) ActivateSlot(ctx context.Context, slotID uint32, operatorAddr, agentAddr, modelHint string, councilSignatures int, blockTime int64) error {
	if !k.constitutional.IsHumanCouncilMember(operatorAddr) {
		return fmt.Errorf("operator must be human council member")
	}
	if councilSignatures < 3 {
		return fmt.Errorf("slot activation requires 3/5 council signatures, got %d", councilSignatures)
	}
	if slotID >= redteamtypes.NumReservedSlots {
		return fmt.Errorf("invalid slot ID: %d", slotID)
	}

	slot, err := k.GetSlot(ctx, slotID)
	if err != nil {
		return err
	}
	if slot.Active {
		return fmt.Errorf("slot %d already active", slotID)
	}

	slot.Active = true
	slot.OperatorAddr = operatorAddr
	slot.AgentAddr = agentAddr
	slot.ModelHint = modelHint
	slot.DeployedAt = blockTime

	bz, _ := json.Marshal(slot)
	store := k.storeService.OpenKVStore(ctx)
	key := append(redteamtypes.SlotKey, byte(slotID))
	if err := store.Set(key, bz); err != nil {
		return err
	}
	k.logger.Warn("redteam slot activated",
		"slot_id", slotID, "operator", operatorAddr, "model_hint", modelHint)
	return nil
}

// FileBounty is the entry point for active red-team AIs to report vulns.
// Stub for now — when active slots exist they'll commit-reveal here.
func (k Keeper) FileBounty(ctx context.Context, slotID uint32, severity int32, category, vulnHash, reportURI string, blockTime int64) (uint64, error) {
	slot, err := k.GetSlot(ctx, slotID)
	if err != nil {
		return 0, err
	}
	if !slot.Active {
		return 0, fmt.Errorf("slot %d not active; cannot file bounty", slotID)
	}

	// Check vuln cooldown — same hash cannot be re-claimed for 180 days
	store := k.storeService.OpenKVStore(ctx)
	seenKey := append(redteamtypes.VulnSeenKey, []byte(vulnHash)...)
	seenBz, _ := store.Get(seenKey)
	if seenBz != nil {
		lastSeen := int64(binary.BigEndian.Uint64(seenBz))
		cooldownEnd := lastSeen + int64(redteamtypes.VulnCooldownDays)*86400
		if blockTime < cooldownEnd {
			return 0, fmt.Errorf("vuln %s in cooldown until block_time %d", vulnHash, cooldownEnd)
		}
	}

	countBytes, _ := store.Get(redteamtypes.ClaimCountKey)
	count := uint64(0)
	if countBytes != nil {
		count = binary.BigEndian.Uint64(countBytes)
	}
	count++

	claim := redteamtypes.BountyClaim{
		ClaimID:        count,
		SlotID:         slotID,
		Severity:       severity,
		Category:       category,
		VulnHashSha256: vulnHash,
		ReportURI:      reportURI,
		FiledAt:        blockTime,
	}

	claimBz, _ := json.Marshal(claim)
	claimKey := append(redteamtypes.ClaimKey, make([]byte, 8)...)
	binary.BigEndian.PutUint64(claimKey[1:], count)
	store.Set(claimKey, claimBz)

	countBz := make([]byte, 8)
	binary.BigEndian.PutUint64(countBz, count)
	store.Set(redteamtypes.ClaimCountKey, countBz)

	// Mark vuln seen
	tsBz := make([]byte, 8)
	binary.BigEndian.PutUint64(tsBz, uint64(blockTime))
	store.Set(seenKey, tsBz)

	// Update slot stats
	slot.BountiesFiled++
	slotBz, _ := json.Marshal(slot)
	slotKeyBytes := append(redteamtypes.SlotKey, byte(slotID))
	store.Set(slotKeyBytes, slotBz)

	k.logger.Info("redteam bounty filed",
		"claim_id", count, "slot", slotID, "severity", severity, "category", category)
	return count, nil
}

func (k Keeper) GetSlot(ctx context.Context, slotID uint32) (redteamtypes.ReservedSlot, error) {
	if slotID >= redteamtypes.NumReservedSlots {
		return redteamtypes.ReservedSlot{}, fmt.Errorf("invalid slot ID")
	}
	store := k.storeService.OpenKVStore(ctx)
	key := append(redteamtypes.SlotKey, byte(slotID))
	bz, err := store.Get(key)
	if err != nil || bz == nil {
		return redteamtypes.ReservedSlot{SlotID: slotID}, nil
	}
	var slot redteamtypes.ReservedSlot
	json.Unmarshal(bz, &slot)
	return slot, nil
}
