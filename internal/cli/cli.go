package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
