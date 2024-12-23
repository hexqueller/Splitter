package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

		if input == "" || input == "y" {
			return true
		} else if input == "n" {
			return false
		} else {
			fmt.Println("Invalid choice, please enter 'y' or 'n'.")
		}
	}
}

func handleFileFlag() (string, os.FileInfo, error) {
	filePath := flag.String("f", "", "Path to file")
	flag.Parse()

	if *filePath == "" {
		return "", nil, fmt.Errorf("usage: ./main.go -f file")
	}

	info, err := os.Stat(*filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil, fmt.Errorf("file not found")
		}
		return "", nil, fmt.Errorf("error: %v", err)
	}

	if info.IsDir() {
		return "", nil, fmt.Errorf("file cannot be a directory")
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

	buffer := make([]byte, partSize)

	for i := int64(0); i < numParts; i++ {
		currentPartSize := partSize
		if i == numParts-1 {
			currentPartSize += remainder
		}

		bytesRead, err := io.ReadFull(file, buffer[:currentPartSize])
		if err != nil && err != io.EOF {
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

func main() {
	filePath, info, err := handleFileFlag()
	if err != nil {
		fmt.Println(err)
		return
	}

	filename := filepath.Base(filePath)
	size := info.Size()

	fmt.Printf("File: %s\n", filename)
	fmt.Printf("Size: %d byte\n", size)

	parts := SelectFromArray(Divisors(size))

	if UserConfirm(fmt.Sprintf("Split file into %d chunks of %d bytes", parts, size/int64(parts))) {
		if err := splitFileByParts(filePath, parts); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("File successfully split!")
		}
	}
}
