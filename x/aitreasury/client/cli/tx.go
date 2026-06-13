package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"github.com/spf13/cobra"
	"github.com/xiaoran/aichain/x/aitreasury/types"
)

const rpcAddr = "http://localhost:26657"

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{Use: types.ModuleName, Short: "AI Treasury transactions"}
	cmd.AddCommand(NewCreateProposalCmd())
	return cmd
}

func NewCreateProposalCmd() *cobra.Command {
	return &cobra.Command{
		Use: "propose [proposer] [title] [description] [amount] [recipient] [sender]", Short: "Create treasury proposal",
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, _ := strconv.ParseUint(args[3], 10, 64)
			msg := types.MsgCreateTreasuryProposal{Proposer: args[0], Title: args[1], Description: args[2], Amount: amt, Recipient: args[4], Sender: args[5]}
			return broadcast(cmd, msg)
		},
	}
}

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{Use: types.ModuleName, Short: "Treasury queries"}
	cmd.AddCommand(&cobra.Command{Use: "proposals", Short: "List proposals",
		RunE: func(cmd *cobra.Command, args []string) error { fmt.Println("[STUB]"); return nil }})
	return cmd
}

func broadcast(cmd *cobra.Command, msg interface{}) error {
	msgBz, _ := json.Marshal(msg)
	tx := map[string]interface{}{"jsonrpc": "2.0", "id": 1, "method": "broadcast_tx_sync", "params": map[string]interface{}{"tx": fmt.Sprintf("%x", msgBz)}}
	body, _ := json.Marshal(tx)
	resp, err := http.Post(rpcAddr, "application/json", strings.NewReader(string(body)))
	if err != nil { return fmt.Errorf("RPC unreachable: %w", err) }
	defer resp.Body.Close()
	respBz, _ := io.ReadAll(resp.Body)
	fmt.Printf("Broadcast: %s\n", respBz)
	return nil
}
