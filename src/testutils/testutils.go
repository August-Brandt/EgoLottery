package testutils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CreateTestDirectory() (string, error) {
	return CreateTestDirectoryFromPath(".")
}

func CreateTestDirectoryFromPath(path string) (string, error) {
	name, err := os.MkdirTemp(path, "test_*")
	if err != nil {
		return name, err
	}
	temp := strings.Split(name, "/")
	folderName := temp[len(temp)-1]
	// Filling test folder with stuff so there is something to look at
	file1, err := os.Create(filepath.Join(path, folderName, ".gitignore"))
	if err != nil {
		return name, err
	}
	file1.Close()
	file2, err := os.Create(filepath.Join(path, folderName, "README.md"))
	if err != nil {
		return name, err
	}
	file2.Close()

	gitCmd := exec.Command("git", "init", name)
	_, err = gitCmd.Output()
	if err != nil {
		return name, err
	}
	// err = os.Mkdir(filepath.Join(path, folderName, ".git"), 0644)
	// if err != nil {
	// 	return name, err
	// }
	err = os.Mkdir(filepath.Join(path, folderName, "src"), 0644)
	if err != nil {
		return name, err
	}
	return name, nil
}
