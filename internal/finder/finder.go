package finder

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/hexqueller/Splitter/internal/base64names"
)

func FindMissingParts(filePath string) ([]string, error) {
	dir := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)

	decodedName, _, totalParts, err := base64names.DecodeBase64(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file name: %w", err)
	}

	foundFiles := []string{}
	missingParts := []int{}
	for i := 1; i <= totalParts; i++ {
		expectedName := base64names.EncodeBase64(decodedName, i, totalParts)
		expectedPath := filepath.Join(dir, expectedName)

		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			missingParts = append(missingParts, i)
		} else {
			foundFiles = append(foundFiles, expectedPath)
		}
	}

	if len(missingParts) > 0 {
		return foundFiles, fmt.Errorf("missing parts: %v", missingParts)
	}

	sortedFiles, err := SortByPartNumber(foundFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to sort files: %w", err)
	}

	return sortedFiles, nil
}

func SortByPartNumber(filePaths []string) ([]string, error) {
	sort.Slice(filePaths, func(i, j int) bool {
		_, partNumberI, _, errI := base64names.DecodeBase64(filepath.Base(filePaths[i]))
		_, partNumberJ, _, errJ := base64names.DecodeBase64(filepath.Base(filePaths[j]))
		if errI != nil || errJ != nil {
			return false
		}
		return partNumberI < partNumberJ
	})
	return filePaths, nil
}
