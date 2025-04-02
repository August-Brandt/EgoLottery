package gitstats

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Repo struct {
	Path, Name string
	Commits    map[int]int
}

func GetStats(repoPaths []string, email, groupType string) []*Repo {
	repos := []*Repo{}
	for _, path := range repoPaths {
		newRepo := &Repo{}
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
			if groupType == "days" {
				err = getCommitsByDay(history, newRepo, email)
				if err != nil {
					log.Printf("Error on getting commits for repo %s:\n%v\n", path, err)
				}
			} else if groupType == "weeks" {
				err = getCommitsByWeek(history, newRepo, email)
				if err != nil {
					log.Printf("Error on getting commits for repo %s:\n%v\n", path, err)
				}
			}
		}
		repos = append(repos, newRepo)
	}
	return repos
}

func weeksInYear(year int) int {
	day := 31
	date := time.Date(year, time.December, day, 0, 0, 0, 0, time.Local)
	_, dateWeek := date.ISOWeek()
	for dateWeek == 1 {
		day--
		date = time.Date(year, time.December, day, 0, 0, 0, 0, time.Local)
		_, dateWeek = date.ISOWeek()
	}
	return dateWeek
}

func getDaysAgo(from time.Time, to time.Time) int {
	diff := to.Sub(from)
	return int(diff.Hours()/24)
}

func getWeeksAgo(from time.Time, to time.Time) int {
	currentYear, currentWeek := to.ISOWeek()
	dateYear, dateWeek := from.ISOWeek()
	if currentYear == dateYear {
		return currentWeek - dateWeek
	} else {
		numWeeksInYear := weeksInYear(dateYear)
		return (numWeeksInYear - dateWeek) + currentWeek
	}
}

func getCommitsByDay(history object.CommitIter, repo *Repo, email string) error {
	history.ForEach(func(c *object.Commit) error {
		if c.Author.Email != email {
			return nil
		}
		daysAgo := getDaysAgo(c.Author.When, time.Now())
		repo.Commits[daysAgo]++
		return nil
	})
	return nil
}

func getCommitsByWeek(history object.CommitIter, repo *Repo, email string) error {
	history.ForEach(func(c *object.Commit) error {
		if c.Author.Email != email {
			return nil
		}
		WeeksAgo := getWeeksAgo(c.Author.When, time.Now())
		repo.Commits[WeeksAgo]++
		return nil
	})
	return nil
}