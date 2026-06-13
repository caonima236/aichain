package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"github.com/spf13/cobra"
	"github.com/xiaoran/aichain/x/agentdao/types"
)

const rpcAddr = "http://localhost:26657"

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{Use: types.ModuleName, Short: "Agent DAO transactions"}
	cmd.AddCommand(NewCreateProposalCmd())
	return cmd
}

func NewCreateProposalCmd() *cobra.Command {
	return &cobra.Command{
		Use: "propose [proposer] [content] [proposal_type] [quorum] [sender]", Short: "Create governance proposal",
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			pt, _ := strconv.Atoi(args[2]); q, _ := strconv.ParseUint(args[3], 10, 64)
			msg := types.MsgCreateDAOProposal{Proposer: args[0], Content: args[1], ProposalType: int32(pt), Quorum: q, Sender: args[4]}
			return broadcast(cmd, msg)
		},
	}
}

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{Use: types.ModuleName, Short: "DAO queries"}
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
	var result map[string]interface{}
	json.Unmarshal(respBz, &result)
	if e, ok := result["error"]; ok { return fmt.Errorf("RPC error: %v\n%s", e, respBz) }
	fmt.Printf("Broadcast: %s\n", respBz)
	return nil
}
