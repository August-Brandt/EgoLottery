package testutils

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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

func AddCommitInRepo(path string) error {
	num := []byte(fmt.Sprintf("%d", rand.IntN(9999999)))
	err := os.WriteFile(filepath.Join(path, "README.md"), num, 0644)
	if err != nil {
		return err
	}
	
	repo, err := git.PlainOpenWithOptions(filepath.Join(path, ".git"), &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = worktree.Add(".")
	if err != nil {
		return err
	}
	
	_, err = worktree.Commit(fmt.Sprint(num), &git.CommitOptions{
		Author: &object.Signature{
			Name: "John Doe",
			Email: "john@doe.org",
			When: time.Now(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}