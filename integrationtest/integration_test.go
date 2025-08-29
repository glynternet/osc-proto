package integrationtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	tempDir := t.TempDir()

	yamlPath, err := filepath.Abs("test.yaml")
	require.NoError(t, err)

	// Change to temp directory so generated files go there
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err)
	}()
	require.NoError(t, os.Chdir(tempDir))

	// Build osc-proto binary first
	binaryPath := filepath.Join(tempDir, "osc-proto")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/osc-proto")
	buildCmd.Dir = filepath.Join(originalDir, "..")
	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed with output: %s", string(buildOutput))

	// Run osc-proto generate command
	cmd := exec.Command(binaryPath, "generate", yamlPath)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command failed with output: %s", string(output))

	// Check that both files were generated
	goFile := filepath.Join(tempDir, "types.go")
	csFile := filepath.Join(tempDir, "types.cs")

	require.FileExists(t, goFile, "Go file should be generated")
	require.FileExists(t, csFile, "C# file should be generated")

	goContent, err := os.ReadFile(goFile)
	require.NoError(t, err)

	csContent, err := os.ReadFile(csFile)
	require.NoError(t, err)

	expectedGoPath := filepath.Join(originalDir, "testdata", "expected_types.go")
	expectedGoContent, err := os.ReadFile(expectedGoPath)
	require.NoError(t, err)

	expectedCsPath := filepath.Join(originalDir, "testdata", "expected_types.cs")
	expectedCsContent, err := os.ReadFile(expectedCsPath)
	require.NoError(t, err)

	t.Run("Go file content", func(t *testing.T) {
		assert.Equal(t, string(expectedGoContent), string(goContent), "Generated Go file should match expected content")
	})

	t.Run("C# file content", func(t *testing.T) {
		assert.Equal(t, string(expectedCsContent), string(csContent), "Generated C# file should match expected content")
	})
}
