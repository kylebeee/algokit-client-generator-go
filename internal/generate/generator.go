package generate

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	algokit "github.com/kylebeee/algokit-utils-go"
)

//go:embed *.go.tmpl
var templateFS embed.FS

// Options configures the code generator.
type Options struct {
	AppSpecPath   string
	OutputDir     string
	PackageName   string
	Mode          string // "full" or "minimal"
	PreserveNames bool
}

// Generate generates typed Go client code from an ARC-56 contract specification.
func Generate(contract *algokit.Arc56Contract, opts Options) error {
	// Determine package name
	packageName := opts.PackageName
	if packageName == "" {
		packageName = ToPackageName(contract.Name)
	}

	// Build generator context
	ctx := BuildContext(contract, packageName, opts.Mode, opts.PreserveNames)

	// Serialize app spec JSON and quote it as a Go string literal
	specJSON, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("failed to marshal app spec: %w", err)
	}
	ctx.AppSpecJSON = strconv.Quote(string(specJSON))

	// Create output directory
	if err := os.MkdirAll(opts.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Parse all templates
	funcMap := template.FuncMap{
		"join":    strings.Join,
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
		"comment": func(s string) string {
			lines := strings.Split(s, "\n")
			for i, line := range lines {
				lines[i] = "// " + line
			}
			return strings.Join(lines, "\n")
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templateFS, "*.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	// Build template data
	data := buildTemplateData(ctx, contract)

	// Generate each file
	files := map[string]string{
		"appspec.go":  "appspec.go.tmpl",
		"types.go":    "types.go.tmpl",
		"client.go":   "client.go.tmpl",
		"composer.go": "composer.go.tmpl",
	}

	if ctx.HasFactory {
		files["factory.go"] = "factory.go.tmpl"
	}

	for filename, tmplName := range files {
		if err := renderTemplate(tmpl, tmplName, data, filepath.Join(opts.OutputDir, filename)); err != nil {
			return fmt.Errorf("failed to generate %s: %w", filename, err)
		}
	}

	return nil
}

// templateData is the data passed to templates.
type templateData struct {
	PackageName              string
	ContractName             string
	AppSpecJSON              string
	Methods                  []MethodData
	Structs                  []StructData
	State                    StateData
	BareConfig               BareCallConfig
	HasFactory               bool
	HasMethodCreateWithArgs  bool
	HasMethodCreateNoArgs    bool
	CreateMethodOriginalName string
	CreateMethodGoName       string
	TypesImports             []string
	ClientImports            []string
}

func buildTemplateData(ctx *GeneratorContext, contract *algokit.Arc56Contract) *templateData {
	data := &templateData{
		PackageName:  ctx.PackageName,
		ContractName: ctx.ContractName,
		AppSpecJSON:  ctx.AppSpecJSON,
		Methods:      ctx.Methods,
		Structs:      ctx.Structs,
		State:        ctx.State,
		BareConfig:   ctx.BareConfig,
		HasFactory:   ctx.HasFactory,
	}

	// Compute create method metadata for factory template
	for _, m := range ctx.Methods {
		if m.CallConfig.CanCreate {
			data.CreateMethodOriginalName = m.OriginalName
			data.CreateMethodGoName = m.Name
			if m.HasArgs() {
				data.HasMethodCreateWithArgs = true
			} else {
				data.HasMethodCreateNoArgs = true
			}
			break
		}
	}

	// Compute imports for types.go
	typesImports := make(map[string]bool)
	// Only need algokit import if there are non-void return methods (for SendAppTransactionResult embed)
	hasNonVoidReturn := false
	for _, m := range ctx.Methods {
		if !m.ReturnType.IsVoid {
			hasNonVoidReturn = true
			break
		}
	}
	if hasNonVoidReturn || data.HasMethodCreateWithArgs {
		typesImports["github.com/kylebeee/algokit-utils-go"] = true
	}
	for _, s := range ctx.Structs {
		for _, f := range s.Fields {
			tm := mapType(f.ABIType, contract.Structs)
			for _, imp := range tm.Imports {
				typesImports[imp] = true
			}
		}
	}
	for _, m := range ctx.Methods {
		for _, a := range m.Args {
			tm := mapType(a.ABIType, contract.Structs)
			for _, imp := range tm.Imports {
				typesImports[imp] = true
			}
		}
		if !m.ReturnType.IsVoid {
			for _, imp := range m.ReturnType.Imports {
				typesImports[imp] = true
			}
		}
	}

	data.TypesImports = sortImports(typesImports)

	return data
}

func sortImports(imports map[string]bool) []string {
	var stdlib, external []string
	for imp := range imports {
		if !strings.Contains(imp, ".") {
			stdlib = append(stdlib, imp)
		} else {
			external = append(external, imp)
		}
	}
	sort.Strings(stdlib)
	sort.Strings(external)

	var result []string
	result = append(result, stdlib...)
	result = append(result, external...)
	return result
}

func renderTemplate(tmpl *template.Template, name string, data interface{}, outputPath string) error {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	// Try to format the Go code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// If formatting fails, write the unformatted code for debugging
		if writeErr := os.WriteFile(outputPath, buf.Bytes(), 0o644); writeErr != nil {
			return fmt.Errorf("format error: %w, write error: %v", err, writeErr)
		}
		return fmt.Errorf("generated code for %s has syntax errors: %w", name, err)
	}

	return os.WriteFile(outputPath, formatted, 0o644)
}
