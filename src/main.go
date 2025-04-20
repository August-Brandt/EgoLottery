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

	"github.com/kkyr/fig"
)

type Config struct {
	DefaultDirectoriesFile string   `fig:"file"`
	GroupType              string   `fig:"group" default:"days"`
	TimeAgo                int      `fig:"timeago" default:"150"`
	SearchDepth            int      `fig:"searchdepth" default:"0"`
	Emails                 []string `fig:"emails" validate:"required"`
}

func main() {
	cfg := &Config{}
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	err = fig.Load(cfg, fig.File("config.yaml"), fig.Dirs(".", path.Join(configDir, "egolottery")))
	if err != nil {
		fmt.Println("Error loading config:")
		panic(err)
	}
	fmt.Printf("Initial Config:\n%+v\n", *cfg)

	var foldersFile string
	var commitsGroupType string
	var directoriesInput string
	var depth int
	flag.StringVar(&foldersFile, "file", "", "path to file containing directories to look for .git directory in")
	flag.StringVar(&commitsGroupType, "group", "", "Group commits by [days|weeks]")
	flag.StringVar(&directoriesInput, "dirs", "", "Commaseperated list of directories. Will override the file flag")
	flag.IntVar(&depth, "depth", -1, "The depth to recursively search for .git directories")
	flag.Parse()

	// Change config based on flags
	if commitsGroupType != "" {
		cfg.GroupType = commitsGroupType
	}
	if depth != -1 {
		cfg.SearchDepth = depth
	}

	var directories []string
	if directoriesInput == "" {
		var filePath string
		if foldersFile != "" {
			filePath = foldersFile
		} else {
			filePath = cfg.DefaultDirectoriesFile
		}
		file, err := os.Open(filePath)
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
	fmt.Printf("Config:\n%+v\n", cfg)

	fmt.Println(".git directories found:")
	dirs := gitfinder.FindGitRepos(directories, cfg.SearchDepth)
	for _, dir := range dirs {
		fmt.Println(dir)
	}

	repos := gitstats.GetStats(dirs, "augustbrandt170@gmail.com", cfg.GroupType)
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
