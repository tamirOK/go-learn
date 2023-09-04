package simpledd

import (
	"crypto/rand"
	"log"
	"os"
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

//nolint:unparam
func createFile(dir string, pattern string) *os.File {
	tempfile, err := os.CreateTemp(dir, pattern)
	if err != nil {
		log.Fatalf("Error during creating file %v", err)
	}

	return tempfile
}

func TestCopyWhenOffsetTooBig(t *testing.T) {
	tempFile := createFile("/tmp", "temp")
	defer os.Remove(tempFile.Name())

	err := Copy(tempFile.Name(), "/tmp/copied", 10000, 0)

	require.ErrorContains(t, err, "offset 10000 is greater than source file size 0")
}

func getFileContents(path string) []byte {
	contents, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error during reading file %v", err)
	}

	return contents
}

func writeFile(path string, contents []byte) {
	if err := os.WriteFile(path, contents, 0o644); err != nil {
		log.Fatalf("Error during writing to file %v", err)
	}
}

//nolint:unparam
func getRandomContent(size int) []byte {
	content := make([]byte, size)

	if _, err := rand.Read(content); err != nil {
		log.Fatalf("Error during generating random data %v", err)
	}

	return content
}

func TestCopy(t *testing.T) {
	tempSourceFile := createFile("/tmp", "temp_source")
	defer os.Remove(tempSourceFile.Name())

	tempDestinationFile := createFile("/tmp", "temp_destination")
	defer os.Remove(tempDestinationFile.Name())

	content := getRandomContent(1 << 20)

	writeFile(tempSourceFile.Name(), content)

	err := Copy(tempSourceFile.Name(), tempDestinationFile.Name(), 0, 0)

	require.Nil(t, err)
	require.Equal(
		t,
		getFileContents(tempSourceFile.Name()),
		getFileContents(tempDestinationFile.Name()),
	)
}

func TestCopyWithOffset(t *testing.T) {
	tempSourceFile := createFile("/tmp", "temp_source")
	defer os.Remove(tempSourceFile.Name())

	tempDestinationFile := createFile("/tmp", "temp_destination")
	defer os.Remove(tempDestinationFile.Name())

	content := getRandomContent(1 << 20)

	writeFile(tempSourceFile.Name(), content)

	offset := 1000
	err := Copy(tempSourceFile.Name(), tempDestinationFile.Name(), offset, 0)

	require.Nil(t, err)
	require.Equal(
		t,
		getFileContents(tempDestinationFile.Name()),
		content[offset:],
	)
}

func TestCopyWithLimit(t *testing.T) {
	tempSourceFile := createFile("/tmp", "temp_source")
	defer os.Remove(tempSourceFile.Name())

	tempDestinationFile := createFile("/tmp", "temp_destination")
	defer os.Remove(tempDestinationFile.Name())

	content := getRandomContent(1 << 20)

	writeFile(tempSourceFile.Name(), content)

	limit := 1 << 15
	err := Copy(tempSourceFile.Name(), tempDestinationFile.Name(), 0, limit)

	require.Nil(t, err)
	require.Equal(
		t,
		getFileContents(tempDestinationFile.Name()),
		content[:limit],
	)
}

func TestCopyWithOffsetAndLimit(t *testing.T) {
	tempSourceFile := createFile("/tmp", "temp_source")
	defer os.Remove(tempSourceFile.Name())

	tempDestinationFile := createFile("/tmp", "temp_destination")
	defer os.Remove(tempDestinationFile.Name())

	content := getRandomContent(1 << 20)

	writeFile(tempSourceFile.Name(), content)

	offset := 2742
	limit := 247812
	err := Copy(tempSourceFile.Name(), tempDestinationFile.Name(), offset, limit)

	require.Nil(t, err)
	require.Equal(
		t,
		getFileContents(tempDestinationFile.Name()),
		content[offset:limit+offset],
	)
}
