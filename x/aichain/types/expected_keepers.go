package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeperI defines the interface for account management
type AccountKeeperI interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx sdk.Context, acc sdk.AccountI)
	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
}

// BankKeeper defines the interface for bank operations
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}
