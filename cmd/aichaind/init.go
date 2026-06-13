package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
)

// genesisAppState is the genesis state for our AI chain modules.
var genesisAppState = map[string]interface{}{
	"constitutional": map[string]interface{}{
		"params": map[string]interface{}{
			"council_members":  []string{},
			"emergency_paused": false,
			"pause_level":      0,
		},
	},
	"agentregistry": map[string]interface{}{
		"agents": []interface{}{},
	},
	"skillnft": map[string]interface{}{
		"skills": []interface{}{},
	},
	"redteam": map[string]interface{}{
		"slots": []map[string]interface{}{
			{"slot_id": 0, "active": false},
			{"slot_id": 1, "active": false},
			{"slot_id": 2, "active": false},
		},
	},
	"poison": map[string]interface{}{
		"state": map[string]interface{}{
			"poison_triggered": false,
			"chain_frozen":     false,
		},
	},
}

// NewInitCmd creates our own init command that generates proper validator keys
// and writes an AICHAIN-specific genesis.json. Avoids the Cosmos SDK InitCmd's
// dependency on a fully-wired ModuleManager.
func NewInitCmd(defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize AI chain home: validator keys, config, genesis",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			moniker := args[0]
			home, _ := cmd.Flags().GetString("home")
			if home == "" {
				home = defaultHome
			}
			chainID, _ := cmd.Flags().GetString("chain-id")
			if chainID == "" {
				chainID = "aichain-testnet-1"
			}

			cfgDir := filepath.Join(home, "config")
			dataDir := filepath.Join(home, "data")
			if err := os.MkdirAll(cfgDir, 0755); err != nil {
				return err
			}
			if err := os.MkdirAll(dataDir, 0755); err != nil {
				return err
			}

			cmtConfig := cmtcfg.DefaultConfig()
			cmtConfig.SetRoot(home)
			cmtConfig.Moniker = moniker
			cmtConfig.P2P.ListenAddress = "tcp://0.0.0.0:26656"
			cmtConfig.RPC.ListenAddress = "tcp://127.0.0.1:26657"
			cmtConfig.Consensus.TimeoutCommit = 5 * time.Second
			cmtConfig.LogLevel = "info"
			cmtcfg.WriteConfigFile(filepath.Join(cfgDir, "config.toml"), cmtConfig)

			// Generate validator + node keys
			nodeID, _, err := genutil.InitializeNodeValidatorFiles(cmtConfig)
			if err != nil {
				return fmt.Errorf("init node validator files: %w", err)
			}

			// Write genesis.json
			genesis := map[string]interface{}{
				"genesis_time":  time.Now().UTC().Format(time.RFC3339),
				"chain_id":      chainID,
				"initial_height": "1",
				"consensus_params": map[string]interface{}{
					"block": map[string]interface{}{
						"max_bytes": "22020096",
						"max_gas":   "-1",
					},
					"evidence": map[string]interface{}{
						"max_age_num_blocks": "100000",
						"max_age_duration":   "172800000000000",
						"max_bytes":          "1048576",
					},
					"validator": map[string]interface{}{
						"pub_key_types": []string{"ed25519"},
					},
					"version": map[string]interface{}{
						"app": "0",
					},
				},
				"app_hash":  "",
				"app_state": genesisAppState,
			}
			genesisBz, _ := json.MarshalIndent(genesis, "", "  ")
			if err := os.WriteFile(filepath.Join(cfgDir, "genesis.json"), genesisBz, 0644); err != nil {
				return err
			}

			// Minimal app.toml
			appToml := `minimum-gas-prices = "0.025uaic"
[telemetry]
enabled = false
[api]
enable = true
address = "tcp://0.0.0.0:1317"
[grpc]
enable = true
address = "0.0.0.0:9090"
[state-sync]
snapshot-interval = 0
`
			if err := os.WriteFile(filepath.Join(cfgDir, "app.toml"), []byte(appToml), 0644); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✅ AICHAIN initialized\n")
			fmt.Fprintf(cmd.OutOrStdout(), "   home:     %s\n", home)
			fmt.Fprintf(cmd.OutOrStdout(), "   chain_id: %s\n", chainID)
			fmt.Fprintf(cmd.OutOrStdout(), "   moniker:  %s\n", moniker)
			fmt.Fprintf(cmd.OutOrStdout(), "   node_id:  %s\n", nodeID)
			return nil
		},
	}
	cmd.Flags().String("chain-id", "aichain-testnet-1", "chain id")
	return cmd
}
