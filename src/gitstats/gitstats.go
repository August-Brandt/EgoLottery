package gitstats

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5"
)

type repo struct {
	Path, Name string
	Commits    int
}

func GetStats(repoPaths []string, email string) []*repo {
	repos := []*repo{}
	for _, path := range repoPaths {
		newRepo := &repo{}
		newRepo.Path = path
		newRepo.Name = filepath.Base(strings.Replace(path, "/.git", "", 1))

		repo, err := git.PlainOpen(path)
		if err != nil {
			log.Printf("Error while opening repo: %v\n", err)
		}
		ref, err := repo.Head()
		if err != nil {
			log.Printf("Error while finding head of repo: %v\n", err)
		} else {
			history, err := repo.Log(&git.LogOptions{From: ref.Hash()})
			if err != nil {
				log.Printf("Error while reading log: %v\n", err)
			}
			history.ForEach(func(c *object.Commit) error {
				if c.Author.Email != email {
					return nil
				}
				newRepo.Commits++
				return nil
			})
		}
		repos = append(repos, newRepo)
	}
	return repos
}
