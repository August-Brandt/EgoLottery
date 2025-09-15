package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add /path/to/directory [/path/to/directory] ...",
	Short: "Add a directory to config",
	Long:  "Add the given path as a directory to the 'directories' varibale in the config",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}

		absPaths := getAbsolutePaths(args)
		for _, path := range absPaths {
			Cfg.Directories = append(Cfg.Directories, path)
		}

		err := overwriteConfig()
		if err != nil {
			panic(err)
		}
	},
}

func getAbsolutePaths(paths []string) []string {
	absPaths := []string{}
	for _, argsPath := range paths {
		if _, err := os.Stat(argsPath); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("'%s' is not an existing path\n", argsPath)
			os.Exit(1)
		}
		absPath, err := filepath.Abs(argsPath)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Adding %s to config\n", absPath)
		absPaths = append(absPaths, absPath)
	}

	return absPaths
}

func overwriteConfig() error {
	newValues, err := yaml.Marshal(Cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(cfgFile, newValues, 0644)
	if err != nil {
		return err
	}
	return nil
}
