package simpledd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyWhenSouceFileNotFound(t *testing.T) {
	err := Copy("/nonexistent_source_file", "/tmp/simpledd_result", 0, 0)

	require.ErrorContains(t, err, "/nonexistent_source_file: no such file or directory")
}

func TestCopyWhenSouceFileIsNotRegularFile(t *testing.T) {
	err := Copy("/dev/null", "/tmp/simpledd_result", 0, 0)

	require.ErrorContains(t, err, "source file must be regular file")
}
