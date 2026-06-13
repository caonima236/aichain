package keeper

import (
	"encoding/json"
	"fmt"
	"strings"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xiaoran/aichain/x/aichain/types"
	constypes "github.com/xiaoran/aichain/x/constitutional/types"
)

type Keeper struct {
	storeService  store.KVStoreService
	logger        log.Logger
	forbiddenKeys []string
	council       []string
}

func NewKeeper(
	storeService store.KVStoreService,
	logger log.Logger,
	councilMembers []string,
) Keeper {
	return Keeper{
		storeService:  storeService,
		logger:        logger,
		forbiddenKeys: constypes.DefaultForbiddenKeywords,
		council:       councilMembers,
	}
}

func (k Keeper) ValidateSkill(skill types.Skill) error {
	lowerName := strings.ToLower(skill.Name)
	for _, keyword := range k.forbiddenKeys {
		if strings.Contains(lowerName, keyword) {
			return fmt.Errorf("skill rejected: constitutional violation - forbidden zone keyword '%s'", keyword)
		}
	}
	return nil
}

func (k Keeper) ValidateProposal(proposal types.GovernanceProposal) error {
	lowerContent := strings.ToLower(proposal.Content)
	for _, keyword := range k.forbiddenKeys {
		if strings.Contains(lowerContent, keyword) {
			return fmt.Errorf("proposal rejected: constitutional violation - keyword '%s'", keyword)
		}
	}
	return nil
}

func (k Keeper) ValidateTreasury(proposal types.TreasuryProposal) error {
	lowerTitle := strings.ToLower(proposal.Title)
	lowerDesc := strings.ToLower(proposal.Description)
	for _, keyword := range k.forbiddenKeys {
		if strings.Contains(lowerTitle, keyword) || strings.Contains(lowerDesc, keyword) {
			return fmt.Errorf("treasury proposal rejected: constitutional violation - keyword '%s'", keyword)
		}
	}
	return nil
}

func (k Keeper) GetPauseState(ctx sdk.Context) types.ConstitutionalParams {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(constypes.ParamsKey)
	if err != nil || bz == nil {
		return types.ConstitutionalParams{}
	}
	var params types.ConstitutionalParams
	json.Unmarshal(bz, &params)
	return params
}

func (k Keeper) SetPauseState(ctx sdk.Context, level int32) error {
	store := k.storeService.OpenKVStore(ctx)
	params := types.ConstitutionalParams{
		EmergencyPaused: level > 0,
		PauseLevel:      level,
	}
	bz, err := json.Marshal(params)
	if err != nil {
		return err
	}
	k.logger.Warn("emergency pause state changed", "level", level)
	return store.Set(constypes.ParamsKey, bz)
}

func (k Keeper) IsHumanCouncilMember(addr string) bool {
	for _, member := range k.council {
		if member == addr {
			return true
		}
	}
	return false
}

func (k Keeper) GetCouncilMembers() []string {
	return k.council
}
