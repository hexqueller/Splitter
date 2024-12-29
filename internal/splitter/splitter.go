package splitter

import (
	"fmt"
	"io"
	"os"

	"github.com/hexqueller/Splitter/internal/base64names"
)

func SplitFileByParts(filePath string, numParts int64, size int64, fileName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	partSize := size / numParts
	remainder := size % numParts

	for i := int64(0); i < numParts; i++ {
		currentPartSize := partSize
		if i == numParts-1 {
			currentPartSize += remainder
		}

		buffer := make([]byte, currentPartSize)
		bytesRead, err := io.ReadFull(file, buffer)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return fmt.Errorf("error reading file: %w", err)
		}

		partFileName := base64names.EncodeBase64(fileName, int(i+1), int(numParts))
		partFile, err := os.Create(partFileName)
		if err != nil {
			return fmt.Errorf("error creating part file: %w", err)
		}

		if _, err := partFile.Write(buffer[:bytesRead]); err != nil {
			return fmt.Errorf("error writing to part file: %w", err)
		}
		partFile.Close()
	}

	return nil
}

func MergeFileParts(sortedFilePaths []string, outputFilePath string) error {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()

	for _, partFileName := range sortedFilePaths {
		partFile, err := os.Open(partFileName)
		if err != nil {
			return fmt.Errorf("error opening part file %s: %w", partFileName, err)
		}

		if _, err := io.Copy(outputFile, partFile); err != nil {
			partFile.Close()
			return fmt.Errorf("error writing part file %s to output: %w", partFileName, err)
		}
		partFile.Close()
	}

	fmt.Println("Files successfully merged into", outputFilePath)
	return nil
}

func DeleteFile(filePath string) {
	var err = os.Remove(filePath)
	if err != nil {
		panic(err)
	}
}

func DeleteFileArray(parts []string) {
	for _, path := range parts {
		DeleteFile(path)
	}
}
