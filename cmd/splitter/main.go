package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hexqueller/Splitter/internal/cli"
	"github.com/hexqueller/Splitter/internal/math"
	"github.com/hexqueller/Splitter/internal/splitter"
)

func handleFileFlag() (string, os.FileInfo, error) {
	filePath := flag.String("f", "", "Path to file")
	flag.Parse()
	if *filePath == "" {
		return "", nil, fmt.Errorf("usage: ./main.go -f PathToFile")
	}
	info, err := os.Stat(*filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil, fmt.Errorf("path not found")
		}
		return "", nil, err
	}
	if info.IsDir() {
		return "", nil, fmt.Errorf("directory input not supported for splitting")
	}
	return *filePath, info, nil
}

func main() {
	path, info, err := handleFileFlag()
	if err != nil {
		fmt.Println(err)
		return
	}

	if strings.HasSuffix(path, ".part1") {
		fmt.Printf("Detected part file: %s\n", path)
		if cli.UserConfirm("Merge parts into the original file?") {
			if err := splitter.MergeFileParts(path); err != nil {
				fmt.Println("Error merging files:", err)
			}
		}
		return
	}

	filename := filepath.Base(path)
	size := info.Size()

	fmt.Printf("File: %s\n", filename)
	fmt.Printf("Size: %d bytes\n", size)

	parts := cli.SelectFromArray(math.Divisors(size))

	if cli.UserConfirm(fmt.Sprintf("Split file into %d chunks of %d bytes", parts, size/int64(parts))) {
		if err := splitter.SplitFileByParts(path, parts); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("File successfully split!")
		}
	}
}
