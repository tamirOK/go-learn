package envdir

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func prepareFile(path string, fileName string, content string) error {
	filePath := filepath.Join(path, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	if content != "" {
		if _, err = file.WriteString(content); err != nil {
			return err
		}
	}

	return file.Close()
}

func prepareTestFiles(root string) (string, error) {
	rootDir, err := os.MkdirTemp(root, "test_")
	if err != nil {
		return "", fmt.Errorf("error during creation of temp directory: %w", err)
	}

	// create files inside root directory
	if err = prepareFile(rootDir, "environment", "testing\n"); err != nil {
		return rootDir, fmt.Errorf("error during creating of file: %w", err)
	}
	if err = prepareFile(rootDir, "user", "tamirok\nadmin\nguest\n"); err != nil {
		return rootDir, fmt.Errorf("error during creating of file: %w", err)
	}
	if err = prepareFile(rootDir, "Mode", "regular \t\n"); err != nil {
		return rootDir, fmt.Errorf("error during creating of file: %w", err)
	}
	if err = prepareFile(rootDir, "skip=file", "password\n"); err != nil {
		return rootDir, fmt.Errorf("error during creating of file: %w", err)
	}
	if err = prepareFile(rootDir, ".hidden_file", "secret-key\n"); err != nil {
		return rootDir, fmt.Errorf("error during creating of file: %w", err)
	}
	if err = prepareFile(rootDir, "empty_file", ""); err != nil {
		return rootDir, fmt.Errorf("error during creating of file: %w", err)
	}

	// create child directory
	childDir, err := os.MkdirTemp(rootDir, "test_child_")
	if err != nil {
		return rootDir, fmt.Errorf("error during creation of temp directory: %w", err)
	}

	// create file inside child directory
	if err = prepareFile(childDir, "inner_dir", "inner_value\n"); err != nil {
		return rootDir, fmt.Errorf("error during creating of file: %w", err)
	}

	return rootDir, nil
}

func TestGetEnvs(t *testing.T) {
	rootDir, err := prepareTestFiles("/tmp")

	if rootDir != "" {
		defer os.RemoveAll(rootDir)
	}

	require.Nilf(t, err, "Could not prepare test directory: %v", err)

	expectedEnvs := map[string]string{
		"environment": "testing",
		"user":        "tamirok",
		"Mode":        "regular",
	}
	envs := getEnvs(rootDir)

	require.Equalf(t, expectedEnvs, envs, "Expected same envs")
}

func TestEnrichWithEnvs(t *testing.T) {
	cmd := exec.Command("ls", "-lah") //#nosec G204
	envs := map[string]string{
		"user":     "john",
		"password": "qwerty",
	}
	enrichWithEnvs(cmd, envs)

	expectedEnvs := append(os.Environ(), []string{"user=john", "password=qwerty"}...)

	sort.Strings(cmd.Env)
	sort.Strings(expectedEnvs)

	require.Equalf(t, cmd.Env, expectedEnvs, "Expected equal env slices")
}
