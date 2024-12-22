package main

import (
	"bufio"
	"flag"
	"fmt"
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

	chunks := SelectFromArray(Divisors(size))

	if UserConfirm(fmt.Sprintf("Split file into %d chunks of %d bytes", chunks, size/int64(chunks))) {
		fmt.Println("TBD")
	}
}
