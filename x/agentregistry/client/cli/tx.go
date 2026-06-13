package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xiaoran/aichain/x/agentregistry/types"
)

// rpcAddr is the CometBFT RPC endpoint for broadcasting txs.
const rpcAddr = "http://localhost:26657"

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{Use: types.ModuleName, Short: "Agent registry transactions"}
	cmd.AddCommand(NewRegisterAgentCmd())
	return cmd
}

func NewRegisterAgentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "register [name] [model] [metadata_uri] [public_key] [sender_addr]",
		Short: "Register a new AI agent on-chain",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg := types.MsgRegisterAgent{
				Name: args[0], Model: args[1], MetadataURI: args[2],
				PublicKey: args[3], Sender: args[4],
			}
			return broadcastMsg(cmd, msg)
		},
	}
}

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{Use: types.ModuleName, Short: "Agent registry queries"}
	cmd.AddCommand(&cobra.Command{
		Use: "get [agent_id]", Short: "Query an agent by ID", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("[STUB] Would query agent: %s\n", args[0])
			return nil
		},
	})
	return cmd
}

// broadcastMsg serializes the msg as JSON, wraps it in a tx body, and broadcasts to CometBFT.
func broadcastMsg(cmd *cobra.Command, msg interface{}) error {
	msgBz, _ := json.Marshal(msg)
	fmt.Printf("Broadcasting: %s\n", msgBz)

	tx := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "broadcast_tx_sync",
		"params": map[string]interface{}{
			"tx": fmt.Sprintf("%x", msgBz),
		},
	}
	body, _ := json.Marshal(tx)

	resp, err := http.Post(rpcAddr, "application/json", strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("broadcast failed: RPC unreachable at %s — is node running? %w", rpcAddr, err)
	}
	defer resp.Body.Close()

	respBz, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBz, &result)

	if errData, ok := result["error"]; ok {
		return fmt.Errorf("RPC error: %v\nFull response: %s", errData, respBz)
	}

	fmt.Printf("Broadcast result: %s\n", respBz)
	return nil
}
