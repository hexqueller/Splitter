package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hexqueller/Splitter/internal/base64names"
	"github.com/hexqueller/Splitter/internal/cli"
	"github.com/hexqueller/Splitter/internal/finder"
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
	fileName := filepath.Base(path)
	if base64names.IsEncoded(fileName) {
		name, partNumber, totalParts, err := base64names.DecodeBase64(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}

		if cli.UserConfirm(fmt.Sprintf("%s %d/%d  Merge parts?", name, partNumber, totalParts)) {
			parts, err := finder.FindMissingParts(path)
			if err != nil {
				fmt.Println(err)
				return
			}
			outputFilePath := filepath.Join(filepath.Dir(path), name)
			err = splitter.MergeFileParts(parts, outputFilePath)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	} else {
		size := info.Size()
		fmt.Printf("File: %s\n", fileName)
		fmt.Printf("Size: %d bytes\n", size)

		parts := cli.SelectFromArray(math.Divisors(size))

		if cli.UserConfirm(fmt.Sprintf("Split file into %d chunks of %d bytes", parts, size/int64(parts))) {
			if err := splitter.SplitFileByParts(path, parts, size, fileName); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("File successfully split!")
				splitter.DeleteFile(path)
			}
		}
	}
}
