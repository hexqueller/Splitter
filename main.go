package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func Divisors(num int64) []int64 {
	var divisors []int64
	for i := int64(1); i*i <= num; i++ {
		if num%i == 0 {
			divisors = append(divisors, i)
			if i != num/i {
				divisors = append(divisors, num/i)
			}
		}
	}
	sort.Slice(divisors, func(i, j int) bool {
		return divisors[i] < divisors[j]
	})
	return divisors
}

func SelectFromArray(options []int64) int64 {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Select a number from the list:")
		for i, option := range options {
			fmt.Printf("%d: %d\n", i+1, option)
		}
		fmt.Print("Choice: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input error, try again.")
			continue
		}
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(options) {
			fmt.Println("Wrong number, try again.")
			continue
		}
		return options[choice-1]
	}
}

func UserConfirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input error, try again.")
			continue
		}
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "" {
			return true
		} else if input == "n" {
			return false
		}
		fmt.Println("Invalid choice, please enter 'y' or 'n'.")
	}
}

func handleFileFlag() (string, os.FileInfo, error) {
	filePath := flag.String("f", "", "Path to file or directory")
	flag.Parse()
	if *filePath == "" {
		return "", nil, fmt.Errorf("usage: ./main.go -f file_or_directory")
	}
	info, err := os.Stat(*filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil, fmt.Errorf("path not found")
		}
		return "", nil, fmt.Errorf("error: %v", err)
	}
	return *filePath, info, nil
}

func splitFileByParts(filePath string, numParts int64) error {
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

func mergeFileParts(partFilePath string) error {
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

func main() {
	path, info, err := handleFileFlag()
	if err != nil {
		fmt.Println(err)
		return
	}

	if strings.HasSuffix(path, ".part1") {
		fmt.Printf("Detected part file: %s\n", path)
		if UserConfirm("Merge parts into the original file?") {
			if err := mergeFileParts(path); err != nil {
				fmt.Println("Error merging files:", err)
			}
		}
		return
	}

	if info.IsDir() {
		fmt.Println("Directory input not supported for splitting.")
		return
	}

	filename := filepath.Base(path)
	size := info.Size()

	fmt.Printf("File: %s\n", filename)
	fmt.Printf("Size: %d bytes\n", size)

	parts := SelectFromArray(Divisors(size))

	if UserConfirm(fmt.Sprintf("Split file into %d chunks of %d bytes", parts, size/int64(parts))) {
		if err := splitFileByParts(path, parts); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("File successfully split!")
		}
	}
}
