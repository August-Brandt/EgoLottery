package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a directory to config",
	Long:  "Add the given path as a directory to the 'directories' varibale in the config",
	Run: func(cmd *cobra.Command, args []string) {
		absPaths := getAbsolutePaths(args)
		for _, path := range absPaths {
			Cfg.Directories = append(Cfg.Directories, path)
		}
		fmt.Println("Adding ", strings.Join(absPaths, ", "), "to config")
		err := overwriteConfig()
		if err != nil {
			panic(err)
		}
	},
}

func getAbsolutePaths(paths []string) []string {
	absPaths := []string{}
	for _, argsPath := range paths {
		absPath, err := filepath.Abs(argsPath)
		if err != nil {
			panic(err)
		}
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
