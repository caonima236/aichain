package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	treasurytypes "github.com/xiaoran/aichain/x/aitreasury/types"
)

type msgServer struct{ k Keeper }

func NewMsgServerImpl(k Keeper) treasurytypes.MsgServer { return &msgServer{k} }

var _ treasurytypes.MsgServer = msgServer{}

func (m msgServer) CreateTreasuryProposal(ctx context.Context, msg *treasurytypes.MsgCreateTreasuryProposal) (*treasurytypes.MsgCreateTreasuryProposalResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	id, err := m.k.CreateProposal(sdkCtx, msg.Proposer, msg.Title, msg.Description, msg.Amount, msg.Recipient)
	if err != nil {
		return nil, err
	}
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("treasury_proposal_created",
			sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", id)),
			sdk.NewAttribute("title", msg.Title),
		),
	)
	return &treasurytypes.MsgCreateTreasuryProposalResponse{ProposalId: id}, nil
}
