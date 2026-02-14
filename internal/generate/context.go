package generate

import (
	"strings"

	algokit "github.com/kylebeee/algokit-utils-go"
)

// GeneratorContext holds the processed data needed for template rendering.
type GeneratorContext struct {
	PackageName   string
	ContractName  string
	AppSpecJSON   string
	Methods       []MethodData
	Structs       []StructData
	State         StateData
	BareConfig    BareCallConfig
	HasFactory    bool
	Imports       map[string]bool
	PreserveNames bool
	usedNames     map[string]bool // tracks all type names to avoid collisions
}

// MethodData holds processed data for a single method.
type MethodData struct {
	Name          string // PascalCase Go name
	OriginalName  string // Original method name
	Signature     string // e.g. "hello(string)string"
	Args          []ArgData
	ReturnType    TypeMapping
	CallConfig    MethodCallConfig
	Desc          string
	ABIArgTypes   []string // ABI types for non-transaction args
}

// ArgData holds processed data for a single method argument.
type ArgData struct {
	Name         string // PascalCase Go name
	OriginalName string // Original arg name
	GoType       string // Go type string
	ABIType      string // Original ABI type string
	IsTransaction bool
	IsReference   bool
	StructName   string // If the arg is a struct
}

// StructData holds processed data for a generated struct type.
type StructData struct {
	Name   string // PascalCase Go name
	Fields []StructFieldData
}

// StructFieldData holds processed data for a struct field.
type StructFieldData struct {
	Name     string // PascalCase Go name
	GoType   string // Go type string
	ABIType  string // Original ABI type string
	JSONTag  string // JSON tag for serialization
}

// StateData holds processed data for the contract's state.
type StateData struct {
	HasGlobal bool
	HasLocal  bool
	HasBox    bool
	Global    []StateKeyData
	Local     []StateKeyData
	Box       []StateKeyData
	BoxMaps   []StateMapData
}

// StateKeyData holds processed data for a state key.
type StateKeyData struct {
	Name      string // PascalCase Go name
	Key       string // Base64 encoded key
	ValueType string // Go type for the value
	ABIType   string // ABI type string
	Desc      string
}

// StateMapData holds processed data for a state map.
type StateMapData struct {
	Name      string // PascalCase Go name
	KeyType   string // Go type for the map key
	ValueType string // Go type for the map value
	Prefix    string // Base64 encoded prefix
	Desc      string
}

// BuildContext creates a GeneratorContext from an ARC-56 contract.
func BuildContext(contract *algokit.Arc56Contract, packageName string, mode string, preserveNames bool) *GeneratorContext {
	ctx := &GeneratorContext{
		PackageName:   packageName,
		ContractName:  ToPascalCase(contract.Name),
		Imports:       make(map[string]bool),
		PreserveNames: preserveNames,
		usedNames:     make(map[string]bool),
	}

	ctx.BareConfig = AnalyzeBareConfig(contract.BareActions)
	ctx.HasFactory = mode == "full"

	// Process structs - register names first
	for name, fields := range contract.Structs {
		goName := ToPascalCase(name)
		ctx.usedNames[goName] = true
		sd := StructData{
			Name: goName,
		}
		for _, f := range fields {
			tm := mapType(f.Type, contract.Structs)
			for _, imp := range tm.Imports {
				ctx.Imports[imp] = true
			}
			sd.Fields = append(sd.Fields, StructFieldData{
				Name:    ToPascalCase(f.Name),
				GoType:  tm.GoType,
				ABIType: f.Type,
				JSONTag: f.Name,
			})
		}
		ctx.Structs = append(ctx.Structs, sd)
	}

	// Process methods
	for _, m := range contract.Methods {
		md := MethodData{
			Name:         ToPascalCase(m.Name),
			OriginalName: m.Name,
			Signature:    m.GetSignature(),
			CallConfig:   AnalyzeCallConfig(m),
			Desc:         m.Desc,
		}

		if preserveNames {
			md.Name = ToPascalCase(m.Name)
		}

		// Process args
		for _, arg := range m.Args {
			tm := MapABITypeToGo(arg.Type, contract.Structs, arg.Struct)
			for _, imp := range tm.Imports {
				ctx.Imports[imp] = true
			}

			isTransaction := isTransactionType(arg.Type)
			isReference := isReferenceType(arg.Type)

			ad := ArgData{
				Name:          ToPascalCase(arg.Name),
				OriginalName:  arg.Name,
				GoType:        tm.GoType,
				ABIType:       arg.Type,
				IsTransaction: isTransaction,
				IsReference:   isReference,
				StructName:    tm.StructName,
			}

			md.Args = append(md.Args, ad)

			if !isTransaction {
				md.ABIArgTypes = append(md.ABIArgTypes, arg.Type)
			}
		}

		// Process return type
		md.ReturnType = MapABITypeToGo(m.Returns.Type, contract.Structs, m.Returns.Struct)
		for _, imp := range md.ReturnType.Imports {
			ctx.Imports[imp] = true
		}

		ctx.Methods = append(ctx.Methods, md)
	}

	// Process state
	ctx.State = buildStateData(contract, ctx)

	return ctx
}

func buildStateData(contract *algokit.Arc56Contract, ctx *GeneratorContext) StateData {
	sd := StateData{}

	for name, key := range contract.State.Keys.Global {
		tm := mapType(key.ValueType, contract.Structs)
		for _, imp := range tm.Imports {
			ctx.Imports[imp] = true
		}
		sd.Global = append(sd.Global, StateKeyData{
			Name:      ToPascalCase(name),
			Key:       key.Key,
			ValueType: tm.GoType,
			ABIType:   key.ValueType,
			Desc:      key.Desc,
		})
		sd.HasGlobal = true
	}

	for name, key := range contract.State.Keys.Local {
		tm := mapType(key.ValueType, contract.Structs)
		for _, imp := range tm.Imports {
			ctx.Imports[imp] = true
		}
		sd.Local = append(sd.Local, StateKeyData{
			Name:      ToPascalCase(name),
			Key:       key.Key,
			ValueType: tm.GoType,
			ABIType:   key.ValueType,
			Desc:      key.Desc,
		})
		sd.HasLocal = true
	}

	for name, key := range contract.State.Keys.Box {
		tm := mapType(key.ValueType, contract.Structs)
		for _, imp := range tm.Imports {
			ctx.Imports[imp] = true
		}
		sd.Box = append(sd.Box, StateKeyData{
			Name:      ToPascalCase(name),
			Key:       key.Key,
			ValueType: tm.GoType,
			ABIType:   key.ValueType,
			Desc:      key.Desc,
		})
		sd.HasBox = true
	}

	for name, mapDef := range contract.State.Maps.Box {
		keyTm := mapType(mapDef.KeyType, contract.Structs)
		valTm := mapType(mapDef.ValueType, contract.Structs)
		for _, imp := range keyTm.Imports {
			ctx.Imports[imp] = true
		}
		for _, imp := range valTm.Imports {
			ctx.Imports[imp] = true
		}
		sd.BoxMaps = append(sd.BoxMaps, StateMapData{
			Name:      ToPascalCase(name),
			KeyType:   keyTm.GoType,
			ValueType: valTm.GoType,
			Prefix:    mapDef.Prefix,
			Desc:      mapDef.Desc,
		})
		sd.HasBox = true
	}

	return sd
}

func isTransactionType(t string) bool {
	txnTypes := map[string]bool{
		"pay": true, "txn": true, "appl": true,
		"axfer": true, "acfg": true, "afrz": true, "keyreg": true,
	}
	return txnTypes[t]
}

func isReferenceType(t string) bool {
	return t == "account" || t == "application" || t == "asset"
}

// GetNonTransactionArgs returns only the non-transaction args.
func (m *MethodData) GetNonTransactionArgs() []ArgData {
	var result []ArgData
	for _, a := range m.Args {
		if !a.IsTransaction {
			result = append(result, a)
		}
	}
	return result
}

// HasNonVoidReturn returns true if the method has a non-void return.
func (m *MethodData) HasNonVoidReturn() bool {
	return !m.ReturnType.IsVoid
}

// HasArgs returns true if the method has any args.
func (m *MethodData) HasArgs() bool {
	return len(m.Args) > 0
}

// GetArgsStructName returns the name for the args struct.
func (m *MethodData) GetArgsStructName() string {
	return m.Name + "Args"
}

// GetResultStructName returns the name for the result struct.
func (m *MethodData) GetResultStructName() string {
	return m.Name + "MethodResult"
}

// SortedImports returns sorted imports for the context.
func (ctx *GeneratorContext) SortedImports() []string {
	var imports []string
	for imp := range ctx.Imports {
		imports = append(imports, imp)
	}

	// Sort: stdlib first, then external
	var stdlib, external []string
	for _, imp := range imports {
		if !strings.Contains(imp, ".") {
			stdlib = append(stdlib, imp)
		} else {
			external = append(external, imp)
		}
	}

	var result []string
	result = append(result, stdlib...)
	result = append(result, external...)
	return result
}
