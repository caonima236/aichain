package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xiaoran/aichain/x/aichain/types"
	skillnfttypes "github.com/xiaoran/aichain/x/skillnft/types"
)

type msgServer struct{ k Keeper }

func NewMsgServerImpl(k Keeper) skillnfttypes.MsgServer { return &msgServer{k} }

var _ skillnfttypes.MsgServer = msgServer{}

func (m msgServer) MintSkill(ctx context.Context, msg *skillnfttypes.MsgMintSkill) (*skillnfttypes.MsgMintSkillResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	skillID, err := m.k.MintSkill(sdkCtx, msg.Creator, msg.Name,
		types.SkillType(msg.SkillType), msg.Version, msg.MetadataURI,
		msg.Price, types.LicenseMode(msg.License), msg.RoyaltyBps)
	if err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("skill_minted",
			sdk.NewAttribute("skill_id", skillID),
			sdk.NewAttribute("name", msg.Name),
		),
	)
	return &skillnfttypes.MsgMintSkillResponse{SkillId: skillID}, nil
}
