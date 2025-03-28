package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/August-Brandt/EgoLottery/gitfinder"
	"github.com/August-Brandt/EgoLottery/gitstats"
	"github.com/August-Brandt/EgoLottery/termprinter"
)

func main() {
	var foldersFile string
	flag.StringVar(&foldersFile, "file", "", "path to file containing directories to look for .git directory in")
	depth := flag.Int("depth", 0, "The depth to recursively search for .git directories")
	flag.Parse()
	if foldersFile == "" { // Default file if non given
		foldersFile = "~/.config/egolottery/directories"
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
	directories := strings.Split(string(data), "\n")

	fmt.Println(".git directories found:")
	dirs := gitfinder.FindGitRepos(directories, *depth)
	for _, dir := range dirs {
		fmt.Println(dir)
	}

	repos := gitstats.GetStats(dirs, "augustbrandt170@gmail.com")
	termprinter.PrintGraph(repos)
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
