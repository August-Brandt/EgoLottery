package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	var foldersFile string
	flag.StringVar(&foldersFile, "file", "", "path to file containing folder to look for .git folders in")
	flag.Parse()
	if foldersFile == "" { // Default file if non given
		foldersFile = "~/.config/egolottery/folders"
	}

	file, err := os.Open(foldersFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	folders := strings.Split(string(data), "\n")
	for _, value := range folders {
		PrintPath(value, log.Default())
	}
}

func PrintPath(path string, output *log.Logger) {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	output.Printf("%s:\n", path)
	for _, entry := range entries {
		output.Printf("\t%s\n", entry.Name())
	}
}
