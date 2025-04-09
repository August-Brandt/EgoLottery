package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/August-Brandt/EgoLottery/gitfinder"
	"github.com/August-Brandt/EgoLottery/gitstats"
	"github.com/August-Brandt/EgoLottery/termprinter"
)

func main() {
	var foldersFile string
	var commitsGroupType string
	var directoriesInput string
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	flag.StringVar(&foldersFile, "file", path.Join(configDir, "egolottery", "directories"), "path to file containing directories to look for .git directory in")
	flag.StringVar(&commitsGroupType, "group", "days", "Group commits by [days|weeks]")
	flag.StringVar(&directoriesInput, "dirs", "", "Commaseperated list of directories. Will override the file flag")
	depth := flag.Int("depth", 0, "The depth to recursively search for .git directories")
	flag.Parse()
	
	var directories []string
	if directoriesInput == "" {
		file, err := os.Open(foldersFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}
		directories = strings.Split(string(data), "\n")
	} else {
		directories = strings.Split(directoriesInput, ",")
	}

	fmt.Println(".git directories found:")
	dirs := gitfinder.FindGitRepos(directories, *depth)
	for _, dir := range dirs {
		fmt.Println(dir)
	}

	repos := gitstats.GetStats(dirs, "augustbrandt170@gmail.com", commitsGroupType)
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
