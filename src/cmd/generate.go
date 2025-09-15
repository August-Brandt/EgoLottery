package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/August-Brandt/EgoLottery/gitfinder"
	"github.com/August-Brandt/EgoLottery/gitstats"
	"github.com/August-Brandt/EgoLottery/termprinter"
	"github.com/spf13/cobra"
)

var commitGrouping string
var searchDepth int
var flagDirectories []string
var timeAgo int
var groupingOptions []string = []string{"days", "weeks"}
var outputFile string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Pull data from git repositories and generate visualization",
	Long:  "Pull git commit data from the specified directories and generate and output a graphical visualization",
	Run: func(cmd *cobra.Command, args []string) {
		if outputFile != "" {
			_, err := os.Create(outputFile)
			if err != nil {
				fmt.Printf("Issue with creating or truncating output file")
				os.Exit(1)
			}
		}

		fmt.Println(".git directories found:")
		dirs := gitfinder.FindGitRepos(Cfg.Directories, Cfg.SearchDepth)
		for _, dir := range dirs {
			fmt.Println(dir)
		}

		repos := gitstats.GetStats(dirs, Cfg)
		termprinter.PrintGraph(repos, Cfg, outputFile)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	cobra.OnInitialize(initGenerateConfig)

	generateCmd.Flags().StringVarP(&commitGrouping, "group", "g", "", "Grouping commits by [days|weeks]")
	generateCmd.Flags().IntVar(&searchDepth, "depth", -1, "The depth to recursively search for .git directories")
	generateCmd.Flags().StringSliceVar(&flagDirectories, "dirs", []string{}, "Comma separated list of directories in which to look for git repositories. Replaces directories in config")
	generateCmd.Flags().IntVar(&timeAgo, "timeago", -1, "The amount of time to include in the final graph")
	generateCmd.Flags().StringVarP(&outputFile, "out", "o", "", "Specify an output file which will be used for outputting the visualization instead of printing to the terminal")
}

func initGenerateConfig() {
	if commitGrouping != "" {
		if slices.Contains(groupingOptions, commitGrouping) {
			Cfg.GroupType = commitGrouping
		} else {
			fmt.Printf("'%s' is an invalid group value. Valid values: %s\n", commitGrouping, strings.Join(groupingOptions, "|"))
			os.Exit(1)
		}
	}
	if searchDepth != -1 {
		if searchDepth > 0 {
			Cfg.SearchDepth = searchDepth
		} else {
			fmt.Printf("'%d' is an invalid depth value. Values must be positive\n", searchDepth)
			os.Exit(1)
		}
	}
	if len(flagDirectories) != 0 {
		absPaths := []string{}
		for _, dir := range flagDirectories {
			if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
				fmt.Printf("The path '%s' does not exist on the system. Please only include valid paths\n", dir)
				os.Exit(1)
			}
			absPath, err := filepath.Abs(dir)
			if err != nil {
				panic(err)
			}
			absPaths = append(absPaths, absPath)
		}
		Cfg.Directories = absPaths
	}
	if timeAgo != -1 {
		if timeAgo > 0 {
			Cfg.TimeAgo = timeAgo
		} else {
			fmt.Printf("'%d' is an invalid timeago value. Values must be larger than 0\n", timeAgo)
			os.Exit(1)
		}
	}
}
