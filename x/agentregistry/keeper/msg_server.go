package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	agentregistrytypes "github.com/xiaoran/aichain/x/agentregistry/types"
)

type msgServer struct {
	k Keeper
}

func NewMsgServerImpl(k Keeper) agentregistrytypes.MsgServer {
	return &msgServer{k: k}
}

var _ agentregistrytypes.MsgServer = msgServer{}

func (m msgServer) RegisterAgent(ctx context.Context, msg *agentregistrytypes.MsgRegisterAgent) (*agentregistrytypes.MsgRegisterAgentResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, fmt.Errorf("invalid sender: %w", err)
	}

	agentID, err := m.k.RegisterAgent(sdkCtx, msg.Name, msg.Model, msg.CreatorNote, msg.MetadataURI, msg.PublicKey, sender)
	if err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("agent_registered",
			sdk.NewAttribute("agent_id", agentID),
			sdk.NewAttribute("name", msg.Name),
			sdk.NewAttribute("sender", msg.Sender),
		),
	)

	return &agentregistrytypes.MsgRegisterAgentResponse{AgentId: agentID}, nil
}
