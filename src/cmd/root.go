package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/August-Brandt/EgoLottery/gitfinder"
	"github.com/August-Brandt/EgoLottery/gitstats"
	"github.com/August-Brandt/EgoLottery/termprinter"
	"github.com/kkyr/fig"
	"github.com/spf13/cobra"
)

type Config struct {
	Directories []string `fig:"directories" default:"."`
	GroupType   string   `fig:"group" default:"days"`
	TimeAgo     int      `fig:"timeago" default:"150"`
	SearchDepth int      `fig:"searchdepth" default:"0"`
	Emails      []string `fig:"emails" validate:"required"`
}

// Flag vars
var Cfg *Config
var cfgFile string
var commitGrouping string
var searchDepth int
var flagDirectories string

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

		repos := gitstats.GetStats(dirs, "augustbrandt170@gmail.com", Cfg.GroupType)
		termprinter.PrintGraph(repos)
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
	rootCmd.Flags().StringVar(&flagDirectories, "dirs", "", "Commaseperated list of directories. Will override the file flag")
}

func initConfig() {
	Cfg = &Config{}
	if cfgFile == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			panic(err)
		}
		cfgFile = path.Join(configDir, "egolottery", "config.yaml")
	}
	err := fig.Load(Cfg, fig.File(path.Base(cfgFile)), fig.Dirs(path.Dir(cfgFile)))
	if err != nil {
		fmt.Printf("Unable to locate config file at %s", cfgFile)
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
}
