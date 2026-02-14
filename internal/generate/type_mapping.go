package generate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	algokit "github.com/kylebeee/algokit-utils-go"
)

var (
	staticArrayRegex  = regexp.MustCompile(`^(.+)\[(\d+)\]$`)
	dynamicArrayRegex = regexp.MustCompile(`^(.+)\[\]$`)
	tupleRegex        = regexp.MustCompile(`^\((.+)\)$`)
	ufixedRegex       = regexp.MustCompile(`^ufixed(\d+)x(\d+)$`)
)

// TypeMapping holds the result of mapping an ABI type to a Go type.
type TypeMapping struct {
	GoType     string   // The Go type string
	Imports    []string // Required imports for this type
	IsVoid     bool     // True if the type is void
	IsStruct   bool     // True if this maps to a generated struct
	StructName string   // Name of the struct if IsStruct
}

// MapABITypeToGo converts an ABI type string to a Go type string.
func MapABITypeToGo(abiType string, structs map[string][]algokit.StructField, structName string) TypeMapping {
	// If there's a struct reference, use the struct name
	if structName != "" {
		if _, ok := structs[structName]; ok {
			goName := ToPascalCase(structName)
			return TypeMapping{
				GoType:     goName,
				IsStruct:   true,
				StructName: goName,
			}
		}
	}

	return mapType(abiType, structs)
}

func mapType(abiType string, structs map[string][]algokit.StructField) TypeMapping {
	// Check for struct reference first
	if _, ok := structs[abiType]; ok {
		goName := ToPascalCase(abiType)
		return TypeMapping{
			GoType:     goName,
			IsStruct:   true,
			StructName: goName,
		}
	}

	switch abiType {
	case "void":
		return TypeMapping{IsVoid: true}
	case "bool":
		return TypeMapping{GoType: "bool"}
	case "byte":
		return TypeMapping{GoType: "byte"}
	case "string":
		return TypeMapping{GoType: "string"}
	case "address":
		return TypeMapping{GoType: "types.Address", Imports: []string{"github.com/algorand/go-algorand-sdk/v2/types"}}
	case "bytes", "byte[]":
		return TypeMapping{GoType: "[]byte"}
	case "uint8":
		return TypeMapping{GoType: "uint8"}
	case "uint16":
		return TypeMapping{GoType: "uint16"}
	case "uint32":
		return TypeMapping{GoType: "uint32"}
	case "uint64":
		return TypeMapping{GoType: "uint64"}

	// AVM types
	case "AVMBytes":
		return TypeMapping{GoType: "[]byte"}
	case "AVMString":
		return TypeMapping{GoType: "string"}
	case "AVMUint64":
		return TypeMapping{GoType: "uint64"}

	// Transaction reference types - these must be passed as TransactionWithSigner
	case "pay", "txn", "appl", "axfer", "acfg", "afrz", "keyreg":
		return TypeMapping{GoType: "transaction.TransactionWithSigner", Imports: []string{"github.com/algorand/go-algorand-sdk/v2/transaction"}}

	// Reference types
	case "account":
		return TypeMapping{GoType: "string"}
	case "application":
		return TypeMapping{GoType: "uint64"}
	case "asset":
		return TypeMapping{GoType: "uint64"}
	}

	// Check for uint128+ (big integers)
	if strings.HasPrefix(abiType, "uint") {
		bits, err := strconv.Atoi(strings.TrimPrefix(abiType, "uint"))
		if err == nil && bits > 64 {
			return TypeMapping{GoType: "*big.Int", Imports: []string{"math/big"}}
		}
	}

	// ufixed types
	if ufixedRegex.MatchString(abiType) {
		return TypeMapping{GoType: "*big.Rat", Imports: []string{"math/big"}}
	}

	// Static array: T[N]
	if m := staticArrayRegex.FindStringSubmatch(abiType); m != nil {
		elemType := mapType(m[1], structs)
		return TypeMapping{
			GoType:  fmt.Sprintf("[%s]%s", m[2], elemType.GoType),
			Imports: elemType.Imports,
		}
	}

	// Dynamic array: T[]
	if m := dynamicArrayRegex.FindStringSubmatch(abiType); m != nil {
		elemType := mapType(m[1], structs)
		return TypeMapping{
			GoType:  "[]" + elemType.GoType,
			Imports: elemType.Imports,
		}
	}

	// Tuple: (T1,T2,...)
	if m := tupleRegex.FindStringSubmatch(abiType); m != nil {
		// For unnamed tuples, we generate a struct
		// For now, just use the raw tuple notation
		return TypeMapping{GoType: "[]interface{}"}
	}

	// Fallback for unknown types
	return TypeMapping{GoType: "interface{}"}
}

// SplitTupleTypes splits a tuple type string into its component types.
// e.g., "(uint64,address,uint64[])" -> ["uint64", "address", "uint64[]"]
func SplitTupleTypes(tupleType string) []string {
	// Remove outer parens
	inner := strings.TrimPrefix(strings.TrimSuffix(tupleType, ")"), "(")
	if inner == "" {
		return nil
	}

	var result []string
	depth := 0
	var current strings.Builder

	for _, r := range inner {
		switch r {
		case '(':
			depth++
			current.WriteRune(r)
		case ')':
			depth--
			current.WriteRune(r)
		case ',':
			if depth == 0 {
				result = append(result, current.String())
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
