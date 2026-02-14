package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kylebeee/algokit-client-generator-go/internal/schema"
)

func TestGenerateApplicationEquality(t *testing.T) {
	testGenerate(t, "../../testdata/ApplicationEquality.arc56.json", "applicationequality", "full")
}

func TestGenerateStateDecoding(t *testing.T) {
	testGenerate(t, "../../testdata/StateDecoding.arc56.json", "statedecoding", "full")
}

func TestGenerateXGovRegistry(t *testing.T) {
	testGenerate(t, "../../testdata/XGovRegistry.arc56.json", "xgovregistry", "full")
}

func TestGenerateMinimalMode(t *testing.T) {
	testGenerate(t, "../../testdata/ApplicationEquality.arc56.json", "appequality_minimal", "minimal")
}

// TestGenerateAllAkitaSpecs tests that the generator can produce valid Go code
// for all 79 ARC-56 specs from the akita-sc project.
func TestGenerateAllAkitaSpecs(t *testing.T) {
	akitaDir := "../../testdata/akita"
	entries, err := os.ReadDir(akitaDir)
	if err != nil {
		t.Skipf("akita testdata not available: %v", err)
	}

	var specs []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".arc56.json") {
			specs = append(specs, e.Name())
		}
	}

	if len(specs) == 0 {
		t.Skip("no ARC-56 specs found in testdata/akita/")
	}

	t.Logf("Found %d ARC-56 specs to test", len(specs))

	passed := 0
	failed := 0
	for _, specFile := range specs {
		specPath := filepath.Join(akitaDir, specFile)
		name := strings.TrimSuffix(specFile, ".arc56.json")
		pkgName := ToPackageName(name)

		t.Run(name, func(t *testing.T) {
			contract, err := schema.LoadAppSpec(specPath)
			if err != nil {
				t.Fatalf("failed to load spec: %v", err)
			}

			outputDir := t.TempDir()
			opts := Options{
				AppSpecPath: specPath,
				OutputDir:   outputDir,
				PackageName: pkgName,
				Mode:        "full",
			}

			if err := Generate(contract, opts); err != nil {
				t.Fatalf("generation failed: %v", err)
			}

			// Verify all expected files exist and have valid Go content
			expectedFiles := []string{"appspec.go", "types.go", "client.go", "composer.go", "factory.go"}
			for _, f := range expectedFiles {
				path := filepath.Join(outputDir, f)
				data, err := os.ReadFile(path)
				if err != nil {
					t.Errorf("failed to read %s: %v", f, err)
					continue
				}
				if len(data) == 0 {
					t.Errorf("file %s is empty", f)
				}
			}

			t.Logf("OK: %s (%d methods, %d structs)", contract.Name, len(contract.Methods), len(contract.Structs))
		})
	}

	t.Logf("Results: %d passed, %d failed out of %d specs", passed, failed, len(specs))
}

func testGenerate(t *testing.T, specPath string, pkgName string, mode string) {
	t.Helper()

	contract, err := schema.LoadAppSpec(specPath)
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	outputDir := t.TempDir()

	opts := Options{
		AppSpecPath: specPath,
		OutputDir:   outputDir,
		PackageName: pkgName,
		Mode:        mode,
	}

	if err := Generate(contract, opts); err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	// Verify expected files exist
	expectedFiles := []string{"appspec.go", "types.go", "client.go", "composer.go"}
	if mode == "full" {
		expectedFiles = append(expectedFiles, "factory.go")
	}

	for _, f := range expectedFiles {
		path := filepath.Join(outputDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", f)
		}
	}

	// Verify no factory.go in minimal mode
	if mode == "minimal" {
		factoryPath := filepath.Join(outputDir, "factory.go")
		if _, err := os.Stat(factoryPath); !os.IsNotExist(err) {
			t.Error("factory.go should not exist in minimal mode")
		}
	}

	// Verify generated files are valid Go (formatting succeeds in renderTemplate)
	for _, f := range expectedFiles {
		path := filepath.Join(outputDir, f)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("failed to read %s: %v", f, err)
			continue
		}
		if len(data) == 0 {
			t.Errorf("file %s is empty", f)
		}
		content := string(data)
		if len(content) < 20 {
			t.Errorf("file %s seems too small: %d bytes", f, len(data))
		}
	}

	// Log the method count for visibility
	t.Logf("Generated %d methods for %s (%s mode)", len(contract.Methods), contract.Name, mode)
}
