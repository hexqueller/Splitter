package base64names

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func EncodeBase64(fileName string, partNumber int, totalParts int) string {
	data := fmt.Sprintf("%s|%03d|%03d", url.PathEscape(fileName), partNumber, totalParts)
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func DecodeBase64(encoded string) (string, int, int, error) {
	// fmt.Println("Encoded string:", encoded) // Debug
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to decode base64: %w", err)
	}

	decodedString := string(decodedBytes)
	parts := strings.Split(decodedString, "|")
	if len(parts) != 3 {
		return "", 0, 0, fmt.Errorf("invalid encoded format")
	}

	fileName, err := url.PathUnescape(parts[0])
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to unescape file name: %w", err)
	}

	partNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to parse part number: %w", err)
	}

	totalParts, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to parse total parts: %w", err)
	}

	return fileName, partNumber, totalParts, nil
}

func IsEncoded(encoded string) bool {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return false
	}
	decodedString := string(decodedBytes)

	parts := strings.Split(decodedString, "|")

	if len(parts) != 3 {
		return false
	}

	if _, err := strconv.Atoi(parts[1]); err != nil {
		return false
	}
	if _, err := strconv.Atoi(parts[2]); err != nil {
		return false
	}

	return true
}
