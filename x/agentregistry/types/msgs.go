package types

import (
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgRegisterAgent = "register_agent"
)

var _ sdk.Msg = &MsgRegisterAgent{}

// MsgRegisterAgent registers a new AI agent on-chain.
type MsgRegisterAgent struct {
	Name        string `json:"name"`
	Model       string `json:"model"`
	CreatorNote string `json:"creator_note"`
	MetadataURI string `json:"metadata_uri"`
	PublicKey   string `json:"public_key"`
	Sender      string `json:"sender"`
}

func NewMsgRegisterAgent(name, model, creatorNote, metadataURI, publicKey, sender string) *MsgRegisterAgent {
	return &MsgRegisterAgent{
		Name: name, Model: model, CreatorNote: creatorNote,
		MetadataURI: metadataURI, PublicKey: publicKey, Sender: sender,
	}
}

func (msg *MsgRegisterAgent) Route() string { return ModuleName }
func (msg *MsgRegisterAgent) Type() string  { return TypeMsgRegisterAgent }
func (msg *MsgRegisterAgent) ValidateBasic() error {
	if strings.TrimSpace(msg.Name) == "" {
		return fmt.Errorf("agent name is required")
	}
	if msg.Sender == "" {
		return fmt.Errorf("sender is required")
	}
	return nil
}
func (msg *MsgRegisterAgent) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}
func (msg *MsgRegisterAgent) GetSignBytes() []byte {
	bz, _ := json.Marshal(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterAgent) ProtoMessage()  {}
func (msg *MsgRegisterAgent) Reset()         { *msg = MsgRegisterAgent{} }
func (msg *MsgRegisterAgent) String() string { return fmt.Sprintf("%T", msg) }
