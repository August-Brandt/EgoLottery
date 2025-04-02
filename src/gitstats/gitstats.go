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
	// Check if it's a leap year
	if time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay() == 366 {
		return 52 + 1 // Leap year has 53 weeks in some cases
	}
	return 52 // Normal years usually have 52 full weeks
}

func getDaysAgo(date time.Time) int {
	duration := time.Since(date)
	days := int(duration.Hours() / 24)
	return days
}

func getWeeksAgo(date time.Time) int {
	currentYear, currentWeek := time.Now().ISOWeek()
	fmt.Printf("CurrentYear, CurrentWeek: %d, %d\n", currentYear, currentWeek)
	dateYear, dateWeek := date.ISOWeek()
	fmt.Printf("date: %v | DateYear, DateWeek: %d, %d\n", date, dateYear, dateWeek)
	if currentYear == dateYear {
		fmt.Printf("Weeks ago: %d\n", currentWeek - dateWeek)
		return currentWeek - dateWeek
	} else {
		numWeeksInYear := weeksInYear(dateYear)
		fmt.Printf("NumWeeksInYear: %d | Weeks ago: %d\n", numWeeksInYear, (numWeeksInYear - dateWeek) + currentWeek)
		return (numWeeksInYear - dateWeek) + currentWeek
	}
}

func getCommitsByDay(history object.CommitIter, repo *Repo, email string) error {
	history.ForEach(func(c *object.Commit) error {
		if c.Author.Email != email {
			return nil
		}
		daysAgo := getDaysAgo(c.Author.When)
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
		WeeksAgo := getWeeksAgo(c.Author.When)
		repo.Commits[WeeksAgo]++
		return nil
	})
	return nil
}