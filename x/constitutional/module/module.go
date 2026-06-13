package module

import (
	"encoding/json"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/types/module"
)

type Module struct {
	module.GenesisOnlyAppModule
}

func NewAppModule() Module {
	return Module{
		GenesisOnlyAppModule: module.NewGenesisOnlyAppModule(basic{}),
	}
}

type basic struct{}

func (basic) Name() string { return "constitutional" }
func (basic) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}
func (basic) RegisterInterfaces(codectypes.InterfaceRegistry) {}
func (basic) RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux) {}
func (basic) RegisterTxService(client.Context) {}
func (basic) RegisterTendermintService(client.Context) {}
func (basic) RegisterNodeService(client.Context, config.Config) {}
func (basic) DefaultGenesis(codec.JSONCodec) json.RawMessage { return json.RawMessage(`{}`) }
func (basic) ValidateGenesis(codec.JSONCodec, client.TxEncodingConfig, json.RawMessage) error { return nil }
func (basic) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate { return nil }
func (basic) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage { return json.RawMessage(`{}`) }
