package splitter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func SplitFileByParts(filePath string, numParts int64) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}
	fileSize := fileInfo.Size()

	partSize := fileSize / numParts
	remainder := fileSize % numParts

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

		partFileName := fmt.Sprintf("%s.part%d", filePath, i+1)
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

func MergeFileParts(partFilePath string) error {
	dir := filepath.Dir(partFilePath)
	baseName := strings.TrimSuffix(filepath.Base(partFilePath), filepath.Ext(partFilePath))
	outputFileName := filepath.Join(dir, baseName)

	partFiles, err := filepath.Glob(filepath.Join(dir, baseName+".part*"))
	if err != nil {
		return fmt.Errorf("error finding part files: %w", err)
	}

	sort.Slice(partFiles, func(i, j int) bool {
		return extractPartNumber(partFiles[i]) < extractPartNumber(partFiles[j])
	})

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()

	for _, partFileName := range partFiles {
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

	fmt.Println("Files successfully merged into", outputFileName)
	return nil
}

func extractPartNumber(fileName string) int {
	re := regexp.MustCompile(`\.part(\d+)$`)
	matches := re.FindStringSubmatch(fileName)
	if len(matches) == 2 {
		number, _ := strconv.Atoi(matches[1])
		return number
	}
	return 0
}
