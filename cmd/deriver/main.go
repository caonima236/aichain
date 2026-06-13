package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("ai", "aipub")
	config.Seal()

	raw, _ := os.ReadFile(os.Args[1])
	var pv struct {
		PubKey struct {
			Value string `json:"value"`
		} `json:"pub_key"`
	}
	json.Unmarshal(raw, &pv)
	pubBytes, _ := base64.StdEncoding.DecodeString(pv.PubKey.Value)

	pk := &ed25519.PubKey{Key: pubBytes}
	addr := sdk.AccAddress(pk.Address().Bytes())
	fmt.Println(addr.String())
}
