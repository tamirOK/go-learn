package simpledd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

const ChunkSize = 4096

func prepareSourceFile(sourcePath string, offset int, limit *int) (*os.File, error) {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file with path %s: %w", sourcePath, err)
	}

	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stat: %w", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return nil, fmt.Errorf("source file must be regular file")
	}

	sourceFileSize := int(fileInfo.Size())

	if offset > sourceFileSize {
		return nil, fmt.Errorf("offset %d is greater than source file size %d", offset, sourceFileSize)
	}

	_, err = sourceFile.Seek(int64(offset), 0)

	if err != nil {
		return nil, fmt.Errorf("could not seek source file to offset %d: %w", offset, err)
	}

	if *limit == 0 || *limit > sourceFileSize {
		*limit = sourceFileSize - offset
	}

	return sourceFile, nil
}

func copyFile(sourceFile *os.File, destFile *os.File, limit int) error {
	reader := bufio.NewReader(sourceFile)
	writer := bufio.NewWriter(destFile)
	buffer := make([]byte, ChunkSize)
	bar := progressbar.DefaultBytes(int64(limit), "copying")

	for {
		bytesRead, readErr := reader.Read(buffer)

		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return fmt.Errorf("failed during reading source file: %w", readErr)
		}

		if bytesRead > limit {
			bytesRead = limit
		}

		writtenBytes, writeErr := writer.Write(buffer[:bytesRead])

		if writeErr != nil {
			return fmt.Errorf("failed during writing to result file: %w", readErr)
		}

		bar.Add(writtenBytes)

		limit -= writtenBytes

		if limit <= 0 {
			break
		}
	}

	writer.Flush()
	return nil
}

func Copy(from string, to string, offset int, limit int) error {
	sourceFile, err := prepareSourceFile(from, offset, &limit)
	if err != nil {
		return err
	}

	defer sourceFile.Close()

	destFile, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("failed to open destination file with path %s: %w", to, err)
	}

	defer destFile.Close()

	return copyFile(sourceFile, destFile, limit)
}
