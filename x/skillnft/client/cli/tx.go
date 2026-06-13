package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xiaoran/aichain/x/skillnft/types"
)

const rpcAddr = "http://localhost:26657"

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{Use: types.ModuleName, Short: "Skill NFT transactions"}
	cmd.AddCommand(NewMintSkillCmd())
	return cmd
}

func NewMintSkillCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mint [creator] [name] [skill_type] [version] [metadata_uri] [price] [license] [royalty_bps] [sender]",
		Short: "Mint a new Skill NFT",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, _ := strconv.Atoi(args[2])
			pr, _ := strconv.ParseUint(args[5], 10, 64)
			li, _ := strconv.Atoi(args[6])
			rp, _ := strconv.ParseUint(args[7], 10, 32)
			msg := types.MsgMintSkill{Creator: args[0], Name: args[1], SkillType: int32(st), Version: args[3], MetadataURI: args[4], Price: pr, License: int32(li), RoyaltyBps: uint32(rp), Sender: args[8]}
			return broadcast(cmd, msg)
		},
	}
}

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{Use: types.ModuleName, Short: "Skill NFT queries"}
	cmd.AddCommand(&cobra.Command{
		Use: "get [skill_id]", Short: "Query skill", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("[STUB] Would query: %s\n", args[0]); return nil
		},
	})
	return cmd
}

func broadcast(cmd *cobra.Command, msg interface{}) error {
	msgBz, _ := json.Marshal(msg)
	tx := map[string]interface{}{"jsonrpc": "2.0", "id": 1, "method": "broadcast_tx_sync", "params": map[string]interface{}{"tx": fmt.Sprintf("%x", msgBz)}}
	body, _ := json.Marshal(tx)
	resp, err := http.Post(rpcAddr, "application/json", strings.NewReader(string(body)))
	if err != nil { return fmt.Errorf("RPC unreachable at %s: %w", rpcAddr, err) }
	defer resp.Body.Close()
	respBz, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBz, &result)
	if errData, ok := result["error"]; ok { return fmt.Errorf("RPC error: %v\n%s", errData, respBz) }
	fmt.Printf("Broadcast: %s\n", respBz)
	return nil
}
