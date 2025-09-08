package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/August-Brandt/EgoLottery/config"
	"github.com/August-Brandt/EgoLottery/gitfinder"
	"github.com/August-Brandt/EgoLottery/gitstats"
	"github.com/August-Brandt/EgoLottery/termprinter"
	"github.com/kkyr/fig"
	"github.com/spf13/cobra"
)

// Flag vars
var Cfg *config.Config
var cfgFile string
var commitGrouping string
var searchDepth int
var flagDirectories string
var timeAgo int

var rootCmd = &cobra.Command{
	Use:   "egolottery",
	Short: "EgoLottery is a local git repository visualizer",
	Long: `This is the long 
discription of EgoLottery`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(".git directories found:")
		dirs := gitfinder.FindGitRepos(Cfg.Directories, Cfg.SearchDepth)
		for _, dir := range dirs {
			fmt.Println(dir)
		}

		repos := gitstats.GetStats(dirs, Cfg)
		termprinter.PrintGraph(repos, Cfg)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&cfgFile, "config", "", "Path to config file for EgoLottery. Default is ~/.config/egolottery/config.yaml")
	rootCmd.Flags().StringVarP(&commitGrouping, "group", "g", "", "Grouping commits by [days|weeks]")
	rootCmd.Flags().IntVar(&searchDepth, "depth", -1, "The depth to recursively search for .git directories")
	rootCmd.Flags().StringVar(&flagDirectories, "dirs", "", "Comma separated list of directories. Will override the file flag")
	rootCmd.Flags().IntVar(&timeAgo, "timeago", -1, "The amount of time to include in the final graph")
}

func initConfig() {
	Cfg = &config.Config{}
	var configDir string
	var err error
	if cfgFile == "" {
		configDir, err = os.UserConfigDir()
		if err != nil {
			panic(err)
		}
		cfgFile = path.Join(configDir, "egolottery", "config.yaml")
	}
	err = fig.Load(Cfg, fig.File(path.Base(cfgFile)), fig.Dirs(path.Dir(cfgFile)))
	if err != nil && strings.Contains(err.Error(), "file not found") {
		fmt.Printf("Unable to locate config file at %s\n", cfgFile)
		if cfgFile == path.Join(configDir, "egolottery", "config.yaml") {
			stdinReader := bufio.NewReader(os.Stdin)
			fmt.Print("\nNo config was found at default location. Would you like to create one?[y|n] ")
			answer, err := stdinReader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			answer = strings.TrimSpace(strings.ToLower(answer))

			for answer != "y" && answer != "n" {
				fmt.Printf("'%s' is an invalid answer. Please respond with either 'y' or 'n': ", answer)
				answer, err = stdinReader.ReadString('\n')
				answer = strings.TrimSpace(strings.ToLower(answer))
			}

			if answer == "y" {
				createConfig(stdinReader, cfgFile)
				err = fig.Load(Cfg, fig.File(path.Base(cfgFile)), fig.Dirs(path.Dir(cfgFile)))
				if err != nil {
					panic(err)
				}
			} else {
				os.Exit(1)
			}
		} else {
			os.Exit(1)
		}
	} else if err != nil {
		fmt.Printf("err: %v\n", err)
		fmt.Printf("err type: %T\n", err)

		panic(err)
	}

	if commitGrouping != "" {
		Cfg.GroupType = commitGrouping
	}
	if searchDepth != -1 {
		Cfg.SearchDepth = searchDepth
	}
	if flagDirectories != "" {
		fmt.Println("Hello")
		Cfg.Directories = strings.Split(flagDirectories, ",")
	}
	if timeAgo != -1 {
		if timeAgo < 1 {
			fmt.Println("'timeago' values should be above 1!")
			os.Exit(1)
		}
		Cfg.TimeAgo = timeAgo
	}
}

func createConfig(stdinReader *bufio.Reader, path string) {
	defaultConfig := `group: "days"
timeago: 60
searchdepth: 0

emails:
  - "<email>"

directories:
  - "<working dir>"
`
	fmt.Print("Please enter git email: ")
	email, err := stdinReader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	email = strings.TrimSpace(email)

	defaultConfig = strings.Replace(defaultConfig, "<email>", email, 1)

	fmt.Print("Please enter the path to a git repo: ")
	dir, err := stdinReader.ReadString('\n')
	dir = strings.TrimSpace(dir)
	dir, err = filepath.Abs(dir)
	if err != nil {
		panic(err)
	}
	defaultConfig = strings.Replace(defaultConfig, "<working dir>", dir, 1)

	err = os.WriteFile(path, []byte(defaultConfig), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Config created at %s\n", path)
}
