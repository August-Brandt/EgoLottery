package gitstats

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type repo struct {
	Path, Name string
	Commits    map[int]int
}

func GetStats(repoPaths []string, email string) []*repo {
	repos := []*repo{}
	for _, path := range repoPaths {
		newRepo := &repo{}
		newRepo.Path = path
		newRepo.Name = filepath.Base(strings.Replace(path, "/.git", "", 1))
		newRepo.Commits = make(map[int]int)

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
				daysAgo := getDaysAgo(c.Author.When)
				newRepo.Commits[daysAgo]++
				return nil
			})
		}
		repos = append(repos, newRepo)
	}
	return repos
}

func getDaysAgo(date time.Time) int {
	duration := time.Since(date)
	days := int(duration.Hours() / 24)
	return days
}
