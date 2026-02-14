package cmd

import (
	"fmt"
	"os"

	"github.com/kylebeee/algokit-client-generator-go/internal/generate"
	"github.com/kylebeee/algokit-client-generator-go/internal/schema"
	"github.com/spf13/cobra"
)

var (
	applicationPath string
	outputDir       string
	packageName     string
	mode            string
	preserveNames   bool
)

// generateCmd represents the generate command.
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate typed Go client from an ARC-56/ARC-32 app spec",
	Long: `Generate typed Go client code from an ARC-56 or ARC-32 application specification.

The generated code provides type-safe interaction with Algorand smart contracts
through a Client struct with methods for each ABI method call.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if applicationPath == "" {
			return fmt.Errorf("--application flag is required")
		}
		if outputDir == "" {
			return fmt.Errorf("--output flag is required")
		}

		// Load the app spec
		contract, err := schema.LoadAppSpec(applicationPath)
		if err != nil {
			return fmt.Errorf("failed to load app spec: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Generating Go client for %s...\n", contract.Name)

		// Generate code
		opts := generate.Options{
			AppSpecPath:   applicationPath,
			OutputDir:     outputDir,
			PackageName:   packageName,
			Mode:          mode,
			PreserveNames: preserveNames,
		}

		if err := generate.Generate(contract, opts); err != nil {
			return fmt.Errorf("generation failed: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully generated Go client in %s\n", outputDir)
		return nil
	},
}

func init() {
	generateCmd.Flags().StringVarP(&applicationPath, "application", "a", "", "Path to ARC-56/ARC-32 app spec JSON file")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for generated Go package")
	generateCmd.Flags().StringVarP(&packageName, "package", "p", "", "Go package name (default: derived from contract name)")
	generateCmd.Flags().StringVarP(&mode, "mode", "m", "full", "Generation mode: full or minimal")
	generateCmd.Flags().BoolVar(&preserveNames, "preserve-names", false, "Preserve original method names (don't sanitize)")
}

// GetGenerateCmd returns the generate command for registration.
func GetGenerateCmd() *cobra.Command {
	return generateCmd
}
