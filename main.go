package main

import (
	"fmt"
	"os"

	"github.com/kylebeee/algokit-client-generator-go/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "algokit-client-generator-go",
	Short: "Generate typed Go clients for Algorand smart contracts",
	Long: `AlgoKit Client Generator for Go generates type-safe Go client code
from ARC-56/ARC-32 application specifications.

Generated code provides a Client struct with methods for each ABI method,
a Composer for building atomic transaction groups, and optionally a Factory
for deploying new contract instances.`,
}

func init() {
	rootCmd.AddCommand(cmd.GetGenerateCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
