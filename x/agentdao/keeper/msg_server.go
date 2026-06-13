package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	daotypes "github.com/xiaoran/aichain/x/agentdao/types"
)

type msgServer struct{ k Keeper }

func NewMsgServerImpl(k Keeper) daotypes.MsgServer { return &msgServer{k} }

var _ daotypes.MsgServer = msgServer{}

func (m msgServer) CreateDAOProposal(ctx context.Context, msg *daotypes.MsgCreateDAOProposal) (*daotypes.MsgCreateDAOProposalResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	id, err := m.k.CreateProposal(sdkCtx, msg.Proposer, msg.Content, msg.ProposalType, msg.Quorum)
	if err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("dao_proposal_created",
			sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", id)),
			sdk.NewAttribute("type", fmt.Sprintf("%d", msg.ProposalType)),
		),
	)
	return &daotypes.MsgCreateDAOProposalResponse{ProposalId: id}, nil
}
