package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

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
}
