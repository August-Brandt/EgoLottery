package testutils

import (
	"os"
	"path/filepath"
)

func CreateTestFolder() (string, error) {
	name, err := os.MkdirTemp(".", "")
	if err != nil {
		return name, err
	}
	// Filling test folder with stuff so there is something to look at
	file1, err := os.Create(filepath.Join(".", name, ".gitignore"))
	if err != nil {
		return name, err
	}
	file1.Close()
	file2, err := os.Create(filepath.Join(".", name, "README.md"))
	if err != nil {
		return name, err
	}
	file2.Close()

	err = os.Mkdir(filepath.Join(".", name, ".git"), 0644)
	if err != nil {
		return name, err
	}
	err = os.Mkdir(filepath.Join(".", name, "src"), 0644)
	if err != nil {
		return name, err
	}
	return name, nil
}
