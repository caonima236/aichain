package main

import (
	"os"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/spf13/cobra"

	"github.com/xiaoran/aichain/app"
	agentregistrycli "github.com/xiaoran/aichain/x/agentregistry/client/cli"
	skillnftcli "github.com/xiaoran/aichain/x/skillnft/client/cli"
	aitreasurycli "github.com/xiaoran/aichain/x/aitreasury/client/cli"
	agentdaocli "github.com/xiaoran/aichain/x/agentdao/client/cli"
)

func main() {
	logger := log.NewLogger(os.Stderr)

	rootCmd := &cobra.Command{
		Use:   "aichaind",
		Short: "AI-native Chain Daemon",
		Long:  "A sovereign blockchain for AI agents — by AI, for AI.",
	}

	// Custom init command
	rootCmd.AddCommand(NewInitCmd(app.DefaultNodeHome))

	// Key management
	rootCmd.AddCommand(keys.Commands())

	// Custom module CLI (tx + query commands)
	rootCmd.AddCommand(agentregistrycli.GetTxCmd(), agentregistrycli.GetQueryCmd())
	rootCmd.AddCommand(skillnftcli.GetTxCmd(), skillnftcli.GetQueryCmd())
	rootCmd.AddCommand(aitreasurycli.GetTxCmd(), aitreasurycli.GetQueryCmd())
	rootCmd.AddCommand(agentdaocli.GetTxCmd(), agentdaocli.GetQueryCmd())

	// Standard Cosmos server commands
	noopFlags := func(cmd *cobra.Command) {}
	server.AddCommands(rootCmd, app.DefaultNodeHome, app.NewAICHAINApp, nil, noopFlags)

	if err := svrcmd.Execute(rootCmd, "AICHAIN", app.DefaultNodeHome); err != nil {
		logger.Error("failed to start aichaind", "error", err)
		os.Exit(1)
	}
}
