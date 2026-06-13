package app

import (
	protov2 "google.golang.org/protobuf/proto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SimpleTx is our minimal transaction wrapper - just carries raw JSON bytes.
type SimpleTx struct {
	rawBytes []byte
	msg      *SimpleMsg
}

func (t *SimpleTx) GetMsgs() []sdk.Msg                  { return []sdk.Msg{t.msg} }
func (t *SimpleTx) GetMsgsV2() ([]protov2.Message, error) { return nil, nil }
func (t *SimpleTx) ValidateBasic() error                { return nil }

// SimpleMsg wraps the raw transaction bytes as a single sdk.Msg.
type SimpleMsg struct {
	RawTx []byte `json:"raw_tx"`
}

func (m *SimpleMsg) ProtoMessage()           {}
func (m *SimpleMsg) Reset()                  { *m = SimpleMsg{} }
func (m *SimpleMsg) String() string          { return string(m.RawTx) }
func (m *SimpleMsg) GetSigners() []sdk.AccAddress { return nil }
func (m *SimpleMsg) GetSignBytes() []byte    { return m.RawTx }
func (m *SimpleMsg) ValidateBasic() error    { return nil }

// SimpleTxDecoder converts raw bytes into a SimpleTx.
func SimpleTxDecoder(bz []byte) (sdk.Tx, error) {
	return &SimpleTx{
		rawBytes: bz,
		msg:      &SimpleMsg{RawTx: bz},
	}, nil
}
