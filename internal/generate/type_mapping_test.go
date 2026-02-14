package generate

import (
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"
)

func TestMapABITypeToGo(t *testing.T) {
	structs := map[string][]algokit.StructField{
		"MyStruct": {
			{Name: "a", Type: "uint64"},
			{Name: "b", Type: "address"},
		},
	}

	tests := []struct {
		abiType    string
		structName string
		expected   string
	}{
		{"void", "", ""},
		{"bool", "", "bool"},
		{"byte", "", "byte"},
		{"string", "", "string"},
		{"uint8", "", "uint8"},
		{"uint16", "", "uint16"},
		{"uint32", "", "uint32"},
		{"uint64", "", "uint64"},
		{"uint128", "", "*big.Int"},
		{"uint256", "", "*big.Int"},
		{"address", "", "types.Address"},
		{"bytes", "", "[]byte"},
		{"byte[]", "", "[]byte"},
		{"uint64[]", "", "[]uint64"},
		{"uint64[5]", "", "[5]uint64"},
		{"byte[32]", "", "[32]byte"},
		{"account", "", "string"},
		{"application", "", "uint64"},
		{"asset", "", "uint64"},
		{"pay", "", "transaction.TransactionWithSigner"},
		{"AVMBytes", "", "[]byte"},
		{"AVMString", "", "string"},
		{"AVMUint64", "", "uint64"},
	}

	for _, tt := range tests {
		t.Run(tt.abiType, func(t *testing.T) {
			result := MapABITypeToGo(tt.abiType, structs, tt.structName)
			if tt.expected == "" {
				if !result.IsVoid {
					t.Errorf("expected void for %s", tt.abiType)
				}
			} else if result.GoType != tt.expected {
				t.Errorf("MapABITypeToGo(%q) = %q, want %q", tt.abiType, result.GoType, tt.expected)
			}
		})
	}
}

func TestMapABITypeToGoStruct(t *testing.T) {
	structs := map[string][]algokit.StructField{
		"MyStruct": {
			{Name: "a", Type: "uint64"},
			{Name: "b", Type: "address"},
		},
	}

	result := MapABITypeToGo("(uint64,address)", structs, "MyStruct")
	if !result.IsStruct {
		t.Error("expected IsStruct to be true")
	}
	if result.GoType != "MyStruct" {
		t.Errorf("expected MyStruct, got %s", result.GoType)
	}
}

func TestSplitTupleTypes(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"(uint64,address)", []string{"uint64", "address"}},
		{"(uint64,address,uint64[])", []string{"uint64", "address", "uint64[]"}},
		{"((uint64,address),uint64)", []string{"(uint64,address)", "uint64"}},
		{"()", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := SplitTupleTypes(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("SplitTupleTypes(%q) = %v (len %d), want %v (len %d)", tt.input, result, len(result), tt.expected, len(tt.expected))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("SplitTupleTypes(%q)[%d] = %q, want %q", tt.input, i, v, tt.expected[i])
				}
			}
		})
	}
}
