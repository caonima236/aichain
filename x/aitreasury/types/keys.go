package types

import "encoding/binary"

const (
	ModuleName   = "aitreasury"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	ProposalKey      = []byte{0x01}
	ProposalCountKey = []byte{0x02}
	BalanceKey       = []byte{0x03}
)

func GetProposalKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append(ProposalKey, bz...)
}
