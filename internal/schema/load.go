package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	algokit "github.com/kylebeee/algokit-utils-go"
)

// LoadAppSpec loads an ARC-56 or ARC-32 application specification from a JSON file.
// It auto-detects the format and converts ARC-32 to ARC-56 if needed.
func LoadAppSpec(path string) (*algokit.Arc56Contract, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read app spec file: %w", err)
	}

	// Try to detect format
	if isArc56(data) {
		return algokit.ParseArc56Contract(data)
	}

	// Try ARC-32 conversion
	return algokit.Arc32ToArc56(data)
}

// isArc56 detects whether the JSON data is an ARC-56 spec.
// ARC-56 specs have top-level "methods" with "actions" fields and "bareActions".
// ARC-32 specs have "contract", "hints", "source", "schema" at top level.
func isArc56(data []byte) bool {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return false
	}

	// ARC-56 has "methods" and "bareActions" at top level
	_, hasMethods := raw["methods"]
	_, hasBareActions := raw["bareActions"]
	_, hasContract := raw["contract"]

	// ARC-32 has "contract" at top level; ARC-56 does not
	if hasContract {
		return false
	}

	return hasMethods && hasBareActions
}

// DetectFormat returns "arc56" or "arc32" based on the file content.
func DetectFormat(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	if isArc56(data) {
		return "arc56", nil
	}

	// Check if it could be ARC-32
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if _, hasContract := raw["contract"]; hasContract {
		return "arc32", nil
	}

	// If the filename contains arc56, trust it
	if strings.Contains(strings.ToLower(path), "arc56") {
		return "arc56", nil
	}

	return "unknown", nil
}
