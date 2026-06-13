package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	xaitypes "github.com/xiaoran/aichain/x/xai/types"
)

type Keeper struct {
	storeService store.KVStoreService
	logger       log.Logger
	reports      map[string][]byte // in-memory cache; production reads from IPFS gateway
}

func NewKeeper(storeService store.KVStoreService, logger log.Logger) Keeper {
	return Keeper{
		storeService: storeService,
		logger:       logger,
		reports:      make(map[string][]byte),
	}
}

// SubmitReport stores a reasoning report and returns its content hash.
// In production, AI agents upload to IPFS first, then submit the URI.
func (k Keeper) SubmitReport(ctx context.Context, reportJSON []byte) (string, error) {
	if len(reportJSON) > xaitypes.MaxReportBytes {
		return "", fmt.Errorf("report too large: %d bytes (max %d)", len(reportJSON), xaitypes.MaxReportBytes)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(reportJSON, &report); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	if err := validateReportContent(report); err != nil {
		return "", err
	}

	// Generate URI (in production: IPFS CID; here: simple hash)
	uri := fmt.Sprintf("xai://%x", simpleHash(reportJSON))
	k.reports[uri] = reportJSON
	k.logger.Info("XAI report submitted", "uri", uri, "size", len(reportJSON))
	return uri, nil
}

// ValidateReport checks if a report URI passes XAI requirements.
// Called by Timelock keeper before queuing an L1+ decision.
func (k Keeper) ValidateReport(reasoningURI string) error {
	if reasoningURI == "" {
		return fmt.Errorf("empty reasoning URI")
	}

	bz, ok := k.reports[reasoningURI]
	if !ok {
		// In production: fetch from IPFS
		return fmt.Errorf("report not found at URI: %s", reasoningURI)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(bz, &report); err != nil {
		return fmt.Errorf("report not valid JSON: %w", err)
	}

	return validateReportContent(report)
}

func validateReportContent(report map[string]interface{}) error {
	// 1. Required fields must exist and not be empty
	for _, field := range xaitypes.RequiredFields {
		val, ok := report[field]
		if !ok {
			return fmt.Errorf("missing required field: %s", field)
		}
		str, isStr := val.(string)
		if !isStr || strings.TrimSpace(str) == "" {
			return fmt.Errorf("field %s must be a non-empty string", field)
		}
	}

	// 2. Black-box markers anywhere in the report → reject
	flat := flattenReport(report)
	lowerFlat := strings.ToLower(flat)
	for _, marker := range xaitypes.BlackBoxMarkers {
		if strings.Contains(lowerFlat, marker) {
			return fmt.Errorf("report contains black-box marker: %q (XAI requires explainable reasoning)", marker)
		}
	}

	// 3. human_readable_summary length cap
	summary, _ := report["human_readable_summary"].(string)
	if len(summary) > xaitypes.MaxSummaryBytes {
		return fmt.Errorf("human_readable_summary too long: %d bytes (max %d)", len(summary), xaitypes.MaxSummaryBytes)
	}

	// 4. affected_modules must list at least one module
	mods, _ := report["affected_modules"].(string)
	if strings.TrimSpace(mods) == "" {
		return fmt.Errorf("affected_modules must list at least one module")
	}

	return nil
}

func flattenReport(report map[string]interface{}) string {
	var sb strings.Builder
	for k, v := range report {
		sb.WriteString(k)
		sb.WriteString(": ")
		if str, ok := v.(string); ok {
			sb.WriteString(str)
		} else {
			bz, _ := json.Marshal(v)
			sb.Write(bz)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// simpleHash for MVP. Production uses IPFS CIDs.
func simpleHash(data []byte) uint64 {
	var h uint64 = 14695981039346656037
	for _, b := range data {
		h ^= uint64(b)
		h *= 1099511628211
	}
	return h
}
