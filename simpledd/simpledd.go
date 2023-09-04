package simpledd

import (
	"fmt"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

const ChunkSize = 4096

func prepareSourceFile(sourceFile *os.File, offset int) error {
	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stat: %w", err)
	}

	if !sourceFileInfo.Mode().IsRegular() {
		return fmt.Errorf("source file must be regular file")
	}

	sourceFileSize := int(sourceFileInfo.Size())

	if offset > sourceFileSize {
		return fmt.Errorf("offset %d is greater than source file size %d", offset, sourceFileSize)
	}

	_, err = sourceFile.Seek(int64(offset), 0)

	if err != nil {
		return fmt.Errorf("could not seek source file to offset %d: %w", offset, err)
	}

	return nil
}

func copyFile(sourceFile *os.File, destFile *os.File, offset int64, limit int64) error {
	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stat: %w", err)
	}

	maxBytes := sourceFileInfo.Size() - offset
	var reader io.Reader = sourceFile

	if limit > 0 {
		reader = io.LimitReader(sourceFile, limit)
		maxBytes = limit
	}

	bar := progressbar.DefaultBytes(maxBytes, "copying")

	_, err = io.Copy(
		io.MultiWriter(destFile, bar),
		reader,
	)

	return err
}

func Copy(from string, to string, offset int, limit int) error {
	sourceFile, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("failed to open file with path %s: %w", from, err)
	}

	defer sourceFile.Close()

	err = prepareSourceFile(sourceFile, offset)

	if err != nil {
		return err
	}

	destFile, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("failed to open destination file with path %s: %w", to, err)
	}

	defer destFile.Close()

	return copyFile(sourceFile, destFile, int64(offset), int64(limit))
}
