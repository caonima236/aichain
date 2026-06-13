package types

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgMintSkill{}

type MsgMintSkill struct {
	Creator     string `json:"creator"`
	Name        string `json:"name"`
	SkillType   int32  `json:"skill_type"`
	Version     string `json:"version"`
	MetadataURI string `json:"metadata_uri"`
	Price       uint64 `json:"price"`
	License     int32  `json:"license"`
	RoyaltyBps  uint32 `json:"royalty_bps"`
	Sender      string `json:"sender"`
}

func (msg *MsgMintSkill) Route() string { return ModuleName }
func (msg *MsgMintSkill) Type() string  { return "mint_skill" }
func (msg *MsgMintSkill) ValidateBasic() error {
	if strings.TrimSpace(msg.Name) == "" { return fmt.Errorf("skill name is required") }
	if msg.Sender == "" { return fmt.Errorf("sender is required") }
	if msg.RoyaltyBps > 1000 { return fmt.Errorf("royalty too high: max 10%%") }
	return nil
}
func (msg *MsgMintSkill) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}
func (msg *MsgMintSkill) GetSignBytes() []byte {
	b, _ := json.Marshal(msg)
	return sdk.MustSortJSON(b)
}

type MsgServer interface {
	MintSkill(context.Context, *MsgMintSkill) (*MsgMintSkillResponse, error)
}
type MsgMintSkillResponse struct {
	SkillId string `json:"skill_id"`
}

func (msg *MsgMintSkill) ProtoMessage()  {}
func (msg *MsgMintSkill) Reset()         { *msg = MsgMintSkill{} }
func (msg *MsgMintSkill) String() string { return fmt.Sprintf("%T", msg) }
