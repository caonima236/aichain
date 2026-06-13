package app

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xiaoran/aichain/x/aichain/types"
	agentregistrykeeper "github.com/xiaoran/aichain/x/agentregistry/keeper"
	skillnftkeeper "github.com/xiaoran/aichain/x/skillnft/keeper"
	aitreasurykeeper "github.com/xiaoran/aichain/x/aitreasury/keeper"
	agentdaokeeper "github.com/xiaoran/aichain/x/agentdao/keeper"
)

// SimpleTxHandler routes raw hex-encoded JSON messages to the appropriate module keeper.
// This bypasses Cosmos SDK's protobuf-based MsgServiceRouter by parsing tx body manually.
//
// Tx format: hex-encoded JSON with a "_msg_type" field identifying the target module.
// Example: {"_msg_type":"register_agent","name":"小冉","model":"opus-4-7",...}
type SimpleTxHandler struct {
	agentRegistry  agentregistrykeeper.Keeper
	skillNFT       skillnftkeeper.Keeper
	aiTreasury     aitreasurykeeper.Keeper
	agentDAO       agentdaokeeper.Keeper
}

func NewSimpleTxHandler(ar agentregistrykeeper.Keeper, sn skillnftkeeper.Keeper, at aitreasurykeeper.Keeper, ad agentdaokeeper.Keeper) *SimpleTxHandler {
	return &SimpleTxHandler{ar, sn, at, ad}
}

// ProcessTx handles a raw transaction bytes. Returns events and error.
func (h *SimpleTxHandler) ProcessTx(ctx sdk.Context, txBytes []byte) (string, error) {
	// CometBFT may deliver the bytes raw, hex, or base64. Try each.
	var raw []byte
	if decoded, err := hex.DecodeString(string(txBytes)); err == nil && len(decoded) > 0 {
		raw = decoded
	} else if decoded, err := base64.StdEncoding.DecodeString(string(txBytes)); err == nil && len(decoded) > 0 {
		raw = decoded
	} else {
		raw = txBytes
	}

	// If after decode it's still not JSON, try decoding once more (CometBFT double-encodes sometimes)
	if len(raw) > 0 && raw[0] != '{' {
		if decoded, err := hex.DecodeString(string(raw)); err == nil {
			raw = decoded
		}
	}

	var msg map[string]interface{}
	if err := json.Unmarshal(raw, &msg); err != nil {
		return "", fmt.Errorf("invalid tx: not JSON: %w", err)
	}

	msgType, _ := msg["_msg_type"].(string)
	if msgType == "" {
		// Infer from fields
		if _, ok := msg["agent_id"]; ok {
			msgType = "register_agent"
		} else if _, ok := msg["skill_type"]; ok {
			msgType = "mint_skill"
		} else if _, ok := msg["recipient"]; ok {
			msgType = "treasury_proposal"
		} else if _, ok := msg["proposal_type"]; ok {
			msgType = "dao_proposal"
		}
	}

	switch msgType {
	case "register_agent":
		return h.handleRegisterAgent(ctx, msg)
	case "mint_skill":
		return h.handleMintSkill(ctx, msg)
	case "treasury_proposal":
		return h.handleTreasuryProposal(ctx, msg)
	case "dao_proposal":
		return h.handleDAOProposal(ctx, msg)
	default:
		return "", fmt.Errorf("unknown msg type: %s", msgType)
	}
}

func (h *SimpleTxHandler) handleRegisterAgent(ctx sdk.Context, msg map[string]interface{}) (string, error) {
	name, _ := msg["name"].(string)
	model, _ := msg["model"].(string)
	metaURI, _ := msg["metadata_uri"].(string)
	pubKey, _ := msg["public_key"].(string)
	creatorNote, _ := msg["creator_note"].(string)
	sender, _ := msg["sender"].(string)

	var owner sdk.AccAddress
	if sender != "" {
		addr, err := sdk.AccAddressFromBech32(sender)
		if err == nil {
			owner = addr
		}
	}

	agentID, err := h.agentRegistry.RegisterAgent(ctx, name, model, creatorNote, metaURI, pubKey, owner)
	if err != nil {
		return "", err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("agent_registered",
			sdk.NewAttribute("agent_id", agentID),
			sdk.NewAttribute("name", name),
		),
	)
	return fmt.Sprintf(`{"agent_id":"%s"}`, agentID), nil
}

func (h *SimpleTxHandler) handleMintSkill(ctx sdk.Context, msg map[string]interface{}) (string, error) {
	creator, _ := msg["creator"].(string)
	name, _ := msg["name"].(string)
	skillTypeF, _ := msg["skill_type"].(float64)
	version, _ := msg["version"].(string)
	metaURI, _ := msg["metadata_uri"].(string)
	priceF, _ := msg["price"].(float64)
	licenseF, _ := msg["license"].(float64)
	royaltyF, _ := msg["royalty_bps"].(float64)

	skillID, err := h.skillNFT.MintSkill(ctx, creator, name,
		types.SkillType(int32(skillTypeF)), version, metaURI,
		uint64(priceF), types.LicenseMode(int32(licenseF)), uint32(royaltyF))
	if err != nil {
		// Standardize constitutional violations
		if strings.Contains(err.Error(), "constitutional") {
			return "", fmt.Errorf("ERR_CONSTITUTIONAL_VIOLATION: %v", err)
		}
		return "", err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("skill_minted",
			sdk.NewAttribute("skill_id", skillID),
			sdk.NewAttribute("name", name),
		),
	)
	return fmt.Sprintf(`{"skill_id":"%s"}`, skillID), nil
}

func (h *SimpleTxHandler) handleTreasuryProposal(ctx sdk.Context, msg map[string]interface{}) (string, error) {
	proposer, _ := msg["proposer"].(string)
	title, _ := msg["title"].(string)
	description, _ := msg["description"].(string)
	amountF, _ := msg["amount"].(float64)
	recipient, _ := msg["recipient"].(string)

	id, err := h.aiTreasury.CreateProposal(ctx, proposer, title, description, uint64(amountF), recipient)
	if err != nil {
		if strings.Contains(err.Error(), "constitutional") {
			return "", fmt.Errorf("ERR_CONSTITUTIONAL_VIOLATION: %v", err)
		}
		return "", err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("treasury_proposal_created",
			sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", id)),
			sdk.NewAttribute("title", title),
		),
	)
	return fmt.Sprintf(`{"proposal_id":%d}`, id), nil
}

func (h *SimpleTxHandler) handleDAOProposal(ctx sdk.Context, msg map[string]interface{}) (string, error) {
	proposer, _ := msg["proposer"].(string)
	content, _ := msg["content"].(string)
	ptF, _ := msg["proposal_type"].(float64)
	quorumF, _ := msg["quorum"].(float64)

	id, err := h.agentDAO.CreateProposal(ctx, proposer, content, int32(ptF), uint64(quorumF))
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "council") {
			return "", fmt.Errorf("ERR_HUMAN_COUNCIL_REQUIRED: %v", err)
		}
		if strings.Contains(errStr, "constitutional") {
			return "", fmt.Errorf("ERR_CONSTITUTIONAL_VIOLATION: %v", err)
		}
		return "", err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("dao_proposal_created",
			sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", id)),
			sdk.NewAttribute("type", fmt.Sprintf("%d", int32(ptF))),
		),
	)
	return fmt.Sprintf(`{"proposal_id":%d}`, id), nil
}

// IntoCheckTxResponse converts the result into an ABCI CheckTx response.
func IntoCheckTxResponse(result string, err error) *abci.ResponseCheckTx {
	if err != nil {
		return &abci.ResponseCheckTx{
			Code: 1,
			Log:  err.Error(),
		}
	}
	return &abci.ResponseCheckTx{
		Code: 0,
		Log:  result,
	}
}
