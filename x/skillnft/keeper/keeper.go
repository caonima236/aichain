package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/xiaoran/aichain/x/aichain/types"
	skilltypes "github.com/xiaoran/aichain/x/skillnft/types"
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
	ValidateSkill(skill types.Skill) error
}

func NewKeeper(storeService store.KVStoreService, logger log.Logger, ak AgentKeeper, ck ConstitutionalKeeper) Keeper {
	return Keeper{storeService: storeService, logger: logger, agentKeeper: ak, constitutional: ck}
}

func (k Keeper) MintSkill(ctx context.Context, creator, name string, skillType types.SkillType, version, metadataURI string, price uint64, license types.LicenseMode, royaltyBps uint32) (string, error) {
	_, err := k.agentKeeper.GetAgent(ctx, creator)
	if err != nil {
		return "", fmt.Errorf("creator not a registered agent: %w", err)
	}
	if royaltyBps > 1000 {
		return "", fmt.Errorf("royalty too high: %d bps (max 10%%)", royaltyBps)
	}

	store := k.storeService.OpenKVStore(ctx)
	countBytes, err := store.Get(skilltypes.SkillCountKey)
	if err != nil {
		return "", err
	}
	count := uint64(0)
	if countBytes != nil {
		count = binary.BigEndian.Uint64(countBytes)
	}
	count++
	skillID := fmt.Sprintf("skill-%d", count)

	skill := types.Skill{
		SkillID: skillID, Owner: creator, Creator: creator, Name: name,
		Type: skillType, Version: version, MetadataURI: metadataURI,
		Price: price, UsageCount: 0, Rating: 0, License: license, RoyaltyBps: royaltyBps,
	}
	if err := k.constitutional.ValidateSkill(skill); err != nil {
		return "", err
	}
	bz, err := json.Marshal(skill)
	if err != nil {
		return "", err
	}
	if err := store.Set(skilltypes.GetSkillKey(skillID), bz); err != nil {
		return "", err
	}
	countBz := make([]byte, 8)
	binary.BigEndian.PutUint64(countBz, count)
	store.Set(skilltypes.SkillCountKey, countBz)

	k.logger.Info("skill minted", "skill_id", skillID, "creator", creator, "name", name)
	return skillID, nil
}

func (k Keeper) TransferSkill(ctx context.Context, skillID, newOwner string) error {
	skill, err := k.GetSkill(ctx, skillID)
	if err != nil {
		return err
	}
	skill.Owner = newOwner
	return k.setSkill(ctx, skill)
}

func (k Keeper) GetSkill(ctx context.Context, skillID string) (types.Skill, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(skilltypes.GetSkillKey(skillID))
	if err != nil {
		return types.Skill{}, err
	}
	if bz == nil {
		return types.Skill{}, fmt.Errorf("skill not found: %s", skillID)
	}
	var skill types.Skill
	if err := json.Unmarshal(bz, &skill); err != nil {
		return types.Skill{}, err
	}
	return skill, nil
}

func (k Keeper) setSkill(ctx context.Context, skill types.Skill) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(skill)
	if err != nil {
		return err
	}
	return store.Set(skilltypes.GetSkillKey(skill.SkillID), bz)
}
