package types

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateTreasuryProposal{}

type MsgCreateTreasuryProposal struct {
	Proposer    string `json:"proposer"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Amount      uint64 `json:"amount"`
	Recipient   string `json:"recipient"`
	Sender      string `json:"sender"`
}

func (msg *MsgCreateTreasuryProposal) Route() string { return ModuleName }
func (msg *MsgCreateTreasuryProposal) Type() string  { return "create_treasury_proposal" }
func (msg *MsgCreateTreasuryProposal) ValidateBasic() error {
	if strings.TrimSpace(msg.Title) == "" { return fmt.Errorf("title is required") }
	if msg.Sender == "" { return fmt.Errorf("sender is required") }
	return nil
}
func (msg *MsgCreateTreasuryProposal) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}
func (msg *MsgCreateTreasuryProposal) GetSignBytes() []byte {
	b, _ := json.Marshal(msg)
	return sdk.MustSortJSON(b)
}

type MsgServer interface {
	CreateTreasuryProposal(context.Context, *MsgCreateTreasuryProposal) (*MsgCreateTreasuryProposalResponse, error)
}
type MsgCreateTreasuryProposalResponse struct {
	ProposalId uint64 `json:"proposal_id"`
}

func (msg *MsgCreateTreasuryProposal) ProtoMessage()  {}
func (msg *MsgCreateTreasuryProposal) Reset()         { *msg = MsgCreateTreasuryProposal{} }
func (msg *MsgCreateTreasuryProposal) String() string { return fmt.Sprintf("%T", msg) }
