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

func main() {
	filePath := flag.String("f", "", "Path to file")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Usage: ./main.go -f file")
		return
	}

	info, err := os.Stat(*filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File not found.")
		} else {
			fmt.Println("Error:", err)
		}
		return
	}

	if info.IsDir() {
		fmt.Println("File cant be directory.")
		return
	}

	filename := filepath.Base(*filePath)
	size := info.Size()

	fmt.Printf("File: %s\n", filename)
	fmt.Printf("Size: %d byte\n", size)

	choice := SelectFromArray(Divisors(size))
	fmt.Println(choice)
}
