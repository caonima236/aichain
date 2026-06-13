package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	abci "github.com/cometbft/cometbft/abci/types"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	agentregistrykeeper "github.com/xiaoran/aichain/x/agentregistry/keeper"
	skillnftkeeper "github.com/xiaoran/aichain/x/skillnft/keeper"
	aitreasurykeeper "github.com/xiaoran/aichain/x/aitreasury/keeper"
	agentdaokeeper "github.com/xiaoran/aichain/x/agentdao/keeper"
	constitutionalkeeper "github.com/xiaoran/aichain/x/constitutional/keeper"
	timelockkeeper "github.com/xiaoran/aichain/x/timelock/keeper"
	xaikeeper "github.com/xiaoran/aichain/x/xai/keeper"
	poisonkeeper "github.com/xiaoran/aichain/x/poison/keeper"
	redteamkeeper "github.com/xiaoran/aichain/x/redteam/keeper"
	agentregistrymodule "github.com/xiaoran/aichain/x/agentregistry/module"
	skillnftmodule "github.com/xiaoran/aichain/x/skillnft/module"
	aitreasurymodule "github.com/xiaoran/aichain/x/aitreasury/module"
	agentdaomodule "github.com/xiaoran/aichain/x/agentdao/module"
	constitutionalmodule "github.com/xiaoran/aichain/x/constitutional/module"
	timelockmodule "github.com/xiaoran/aichain/x/timelock/module"
	xaimodule "github.com/xiaoran/aichain/x/xai/module"
	poisonmodule "github.com/xiaoran/aichain/x/poison/module"
	redteammodule "github.com/xiaoran/aichain/x/redteam/module"
)

const (
	AppName   = "aichain"
	BondDenom = "uaic"
)

var DefaultNodeHome string

func init() {
	userHomeDir, _ := os.UserHomeDir()
	DefaultNodeHome = userHomeDir + "/.aichain"

	// Set Bech32 prefixes globally to "ai".
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("ai", "aipub")
	config.SetBech32PrefixForValidator("aivaloper", "aivaloperpub")
	config.SetBech32PrefixForConsensusNode("aivalcons", "aivalconspub")
	config.Seal()
}

// AICHAINApp is the sovereign AI-native chain application.
type AICHAINApp struct {
	*runtime.App
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	AccountKeeper        authkeeper.AccountKeeper
	BankKeeper           bankkeeper.Keeper
	StakingKeeper        *stakingkeeper.Keeper
	ConsensusKeeper      consensuskeeper.Keeper
	AgentRegistryKeeper  agentregistrykeeper.Keeper
	SkillNFTKeeper       skillnftkeeper.Keeper
	AITreasuryKeeper     aitreasurykeeper.Keeper
	AgentDAOKeeper       agentdaokeeper.Keeper
	ConstitutionalKeeper constitutionalkeeper.Keeper
	TimelockKeeper       timelockkeeper.Keeper
	XAIKeeper            xaikeeper.Keeper
	PoisonKeeper         poisonkeeper.Keeper
	RedTeamKeeper        redteamkeeper.Keeper
}

// NewAICHAINApp creates the AI chain. Matches servertypes.AppCreator signature.
func NewAICHAINApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(interfaceRegistry)
	legacyAmino := codec.NewLegacyAmino()
	txConfig := authtx.NewTxConfig(appCodec, authtx.DefaultSignModes)

	bApp := baseapp.NewBaseApp(AppName, logger, db, SimpleTxDecoder, baseapp.SetChainID("aichain-testnet-1"))
	app := &AICHAINApp{
		App:               &runtime.App{BaseApp: bApp},
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		txConfig:          txConfig,
		interfaceRegistry: interfaceRegistry,
	}

	addrCodec := authcodec.NewBech32Codec("ai")
	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	}
	authority := authtypes.NewModuleAddress("gov").String()

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey, consensustypes.StoreKey,
		"agentregistry", "skillnft", "aitreasury", "agentdao", "constitutional",
		"timelock", "xai", "poison", "redteam",
	)

	// Standard Cosmos keepers
	app.ConsensusKeeper = consensuskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[consensustypes.StoreKey]),
		authority,
		runtime.EventService{},
	)
	bApp.SetParamStore(app.ConsensusKeeper.ParamsStore)

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount, maccPerms, addrCodec, "ai", authority,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.AccountKeeper, map[string]bool{}, authority, logger,
	)
	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		app.AccountKeeper, app.BankKeeper, authority,
		authcodec.NewBech32Codec("aivaloper"),
		authcodec.NewBech32Codec("aivalcons"),
	)

	// AI-native keepers
	// Constitutional keeper — reads council members from genesis.json
	// so they're available before InitGenesis (needed for Type 4 proposal checks).
	genDoc, _ := os.ReadFile(DefaultNodeHome + "/config/genesis.json")
	councilMembers := []string{}
	if genDoc != nil {
		var gc struct {
			AppState struct {
				Constitutional struct {
					Params struct {
						CouncilMembers []string `json:"council_members"`
					} `json:"params"`
				} `json:"constitutional"`
			} `json:"app_state"`
		}
		if json.Unmarshal(genDoc, &gc) == nil {
			councilMembers = gc.AppState.Constitutional.Params.CouncilMembers
		}
	}
	app.ConstitutionalKeeper = constitutionalkeeper.NewKeeper(
		runtime.NewKVStoreService(keys["constitutional"]), logger, councilMembers,
	)
	app.AgentRegistryKeeper = agentregistrykeeper.NewKeeper(
		runtime.NewKVStoreService(keys["agentregistry"]), logger,
	)
	app.SkillNFTKeeper = skillnftkeeper.NewKeeper(
		runtime.NewKVStoreService(keys["skillnft"]), logger,
		app.AgentRegistryKeeper, app.ConstitutionalKeeper,
	)
	app.AITreasuryKeeper = aitreasurykeeper.NewKeeper(
		runtime.NewKVStoreService(keys["aitreasury"]), logger,
		app.AgentRegistryKeeper, app.ConstitutionalKeeper,
	)
	app.AgentDAOKeeper = agentdaokeeper.NewKeeper(
		runtime.NewKVStoreService(keys["agentdao"]), logger,
		app.AgentRegistryKeeper, app.ConstitutionalKeeper,
	)

	// Phase 0.5 safety modules
	app.XAIKeeper = xaikeeper.NewKeeper(
		runtime.NewKVStoreService(keys["xai"]), logger,
	)
	app.TimelockKeeper = timelockkeeper.NewKeeper(
		runtime.NewKVStoreService(keys["timelock"]), logger,
		app.XAIKeeper, app.ConstitutionalKeeper,
	)
	app.PoisonKeeper = poisonkeeper.NewKeeper(
		runtime.NewKVStoreService(keys["poison"]), logger,
	)
	app.RedTeamKeeper = redteamkeeper.NewKeeper(
		runtime.NewKVStoreService(keys["redteam"]), logger,
		app.ConstitutionalKeeper,
	)

	// Mount stores BEFORE InitChainer / LoadLatestVersion
	for _, key := range keys {
		app.MountStore(key, storetypes.StoreTypeIAVL)
	}

	// Register standard module interfaces before InitGenesis can decode stored types.
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	stakingtypes.RegisterInterfaces(interfaceRegistry)

	// Module manager — standard + all 9 custom modules
	app.ModuleManager = module.NewManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, nil),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, nil),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, nil),
		agentregistrymodule.NewAppModule(),
		skillnftmodule.NewAppModule(),
		aitreasurymodule.NewAppModule(),
		agentdaomodule.NewAppModule(),
		constitutionalmodule.NewAppModule(),
		timelockmodule.NewAppModule(),
		xaimodule.NewAppModule(),
		poisonmodule.NewAppModule(),
		redteammodule.NewAppModule(),
	)

	// Inject SimpleMsg handler into baseapp's MsgServiceRouter via reflection.
	// This makes the basic msg handler check pass; our actual processing happens
	// in AnteHandler. We just need to satisfy baseapp's existence check.
	InjectSimpleMsgHandler(app.MsgServiceRouter(), func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return &sdk.Result{}, nil // already processed in AnteHandler
	})

	// Custom tx handler — routes JSON-encoded messages to module keepers.
	txHandler := NewSimpleTxHandler(
		app.AgentRegistryKeeper,
		app.SkillNFTKeeper,
		app.AITreasuryKeeper,
		app.AgentDAOKeeper,
	)

	// AnteHandler: process tx + emit events, then return error to STOP baseapp from
	// trying to find a Msg handler (which would fail since we don't use protobuf msgs).
	app.SetAnteHandler(func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		// Skip during simulation and genesis (blockheight 0)
		if simulate || ctx.BlockHeight() == 0 {
			return ctx, nil
		}
		if simpleTx, ok := tx.(*SimpleTx); ok && len(simpleTx.rawBytes) > 0 {
			result, err := txHandler.ProcessTx(ctx, simpleTx.rawBytes)
			if err != nil {
				return ctx, err
			}
			ctx.EventManager().EmitEvent(sdk.NewEvent("tx_processed",
				sdk.NewAttribute("result", result),
			))
		}
		return ctx, nil
	})

	// PostHandler: runs AFTER msg routing. If msg routing failed (which it will,
	// because we don't register protobuf msgs), but AnteHandler already processed,
	// we still succeed. PostHandler can override this.
	app.SetPostHandler(func(ctx sdk.Context, tx sdk.Tx, simulate, success bool) (sdk.Context, error) {
		// Our processing happened in AnteHandler; nothing else needed.
		return ctx, nil
	})

	// Note: MsgServer registration via MsgServiceRouter requires full protobuf
	// type registration (proto.MessageName, etc.). Our JSON-based message types
	// will be wired through a custom tx handler in Phase 2.1.
	// The keeper logic (x/*/keeper/msg_server.go) is implemented and ready.

	// InitChainer: read genesis.json → feed each module's app_state to InitGenesis.
	// ModuleManager.InitGenesis calls x/staking → returns validator updates.
	app.SetInitChainer(func(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
		// Read genesis.json from the home directory.
		genesisPath := DefaultNodeHome + "/config/genesis.json"
		genDoc, err := os.ReadFile(genesisPath)
		if err != nil {
			return nil, fmt.Errorf("read genesis.json: %w", err)
		}
		var genesis struct {
			AppState map[string]json.RawMessage `json:"app_state"`
		}
		if err := json.Unmarshal(genDoc, &genesis); err != nil {
			return nil, fmt.Errorf("parse genesis.json: %w", err)
		}
		return app.ModuleManager.InitGenesis(ctx, appCodec, genesis.AppState)
	})

	// Load latest version after stores mounted
	if err := app.LoadLatestVersion(); err != nil {
		logger.Error("failed to load latest version", "error", err)
	}

	return app
}

// Interface compliance
func (app *AICHAINApp) Name() string                              { return AppName }
func (app *AICHAINApp) LegacyAmino() *codec.LegacyAmino           { return app.legacyAmino }
func (app *AICHAINApp) AppCodec() codec.Codec                     { return app.appCodec }
func (app *AICHAINApp) InterfaceRegistry() codectypes.InterfaceRegistry { return app.interfaceRegistry }
func (app *AICHAINApp) TxConfig() client.TxConfig                 { return app.txConfig }
