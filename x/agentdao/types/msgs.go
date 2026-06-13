package types

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateDAOProposal{}

type MsgCreateDAOProposal struct {
	Proposer     string `json:"proposer"`
	Content      string `json:"content"`
	ProposalType int32  `json:"proposal_type"`
	Quorum       uint64 `json:"quorum"`
	Sender       string `json:"sender"`
}

func (msg *MsgCreateDAOProposal) Route() string { return ModuleName }
func (msg *MsgCreateDAOProposal) Type() string  { return "create_dao_proposal" }
func (msg *MsgCreateDAOProposal) ValidateBasic() error {
	if strings.TrimSpace(msg.Content) == "" { return fmt.Errorf("content is required") }
	if msg.Sender == "" { return fmt.Errorf("sender is required") }
	return nil
}
func (msg *MsgCreateDAOProposal) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}
func (msg *MsgCreateDAOProposal) GetSignBytes() []byte {
	b, _ := json.Marshal(msg)
	return sdk.MustSortJSON(b)
}

// Error codes
const (
	ErrCodeConstitutionalViolation uint32 = 1
	ErrCodeHumanCouncilRequired    uint32 = 2
	ErrCodeTimelockActive          uint32 = 3
)

type MsgServer interface {
	CreateDAOProposal(context.Context, *MsgCreateDAOProposal) (*MsgCreateDAOProposalResponse, error)
}
type MsgCreateDAOProposalResponse struct {
	ProposalId uint64 `json:"proposal_id"`
}

func (msg *MsgCreateDAOProposal) ProtoMessage()  {}
func (msg *MsgCreateDAOProposal) Reset()         { *msg = MsgCreateDAOProposal{} }
func (msg *MsgCreateDAOProposal) String() string { return fmt.Sprintf("%T", msg) }
