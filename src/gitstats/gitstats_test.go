package gitstats

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/August-Brandt/EgoLottery/testutils"
)



func TestGetStatsForDays(t *testing.T) {
	// Setup

	repo, err := testutils.CreateTestDirectory()
	if err != nil {
		t.Fatalf("Error while creating repo: %v", err)
	}
	defer os.RemoveAll(repo)

	err = testutils.AddCommitInRepo(repo)
	if err != nil {
		t.Fatalf("Error while adding commit to repo: %v", err)
	}

	// Test
	repos := GetStats([]string{filepath.Join(repo, ".git")}, "john@doe.org", "days")
	if len(repos) != 1 {
		t.Errorf("Incorrect number of repos\n\tExpected: %d\n\tActual: %d", 1, len(repos))
	}
	if repos[0].Commits[0] != 1 {
		t.Errorf("Incorrect number of commits in repo\n\tExpected: %d\n\tActual: %d", 1, len(repos))
	}
}

func TestGetStatsForWeeks(t *testing.T) {
	// Setup

	repo, err := testutils.CreateTestDirectory()
	if err != nil {
		t.Fatalf("Error while creating repo: %v", err)
	}
	defer os.RemoveAll(repo)

	err = testutils.AddCommitInRepo(repo)
	if err != nil {
		t.Fatalf("Error while adding commit to repo: %v", err)
	}

	// Test	

	repos := GetStats([]string{filepath.Join(repo, ".git")}, "john@doe.org", "days")
	if len(repos) != 1 {
		t.Errorf("Incorrect number of repos\n\tExpected: %d\n\tActual: %d", 1, len(repos))
	}
	if repos[0].Commits[0] != 1 {
		t.Errorf("Incorrect number of commits in repo\n\tExpected: %d\n\tActual: %d", 1, len(repos))
	}
}

func TestGetWeeksAgo(t *testing.T) {
	// Test same week
	if val := getWeeksAgo(time.Now(), time.Now()); val != 0 {
		t.Errorf("Incorrect number of weeks ago\n\tExpected: 0\n\tActual: %d", val)
	}

	// Test last week
	if val := getWeeksAgo(time.Date(2025, time.Month(1), 2, 0, 0, 0, 0, time.Local), time.Date(2025, time.Month(1), 9, 0, 0, 0, 0, time.Local)); val != 1 {
		t.Errorf("Incorrect number of weeks ago\n\tExpected: 1\n\tActual: %d", val)
	}

	// Test accross year
	if val := getWeeksAgo(time.Date(2024, time.Month(12), 19, 0, 0, 0, 0, time.Local), time.Date(2025, time.Month(1), 9, 0, 0, 0, 0, time.Local)); val != 3 {
		t.Errorf("Incorrect number of weeks ago\n\tExpected: 3\n\tActual: %d", val)
	}
}

func TestGetDaysAgo(t *testing.T) {
	// Test same day
	if val := getDaysAgo(time.Now(), time.Now()); val != 0 {
		t.Errorf("Incorrect number of days ago\n\tExpected: 0\n\tActual: %d", val)
	}

	// Test 1 day ago
	if val := getDaysAgo(time.Date(2025, time.Month(1), 2, 0, 0, 0, 0, time.Local), time.Date(2025, time.Month(1), 3, 0, 0, 0, 0, time.Local)); val != 1 {
		t.Errorf("Incorrect number of days ago\n\tExpected: 1\n\tActual: %d", val)
	}

	// Test accross year
	if val := getDaysAgo(time.Date(2024, time.December, 27, 0, 0, 0, 0, time.Local), time.Date(2025, time.Month(1), 2, 0, 0, 0, 0, time.Local)); val != 6 {
		t.Errorf("Incorrect number of days ago\n\tExpected: 6\n\tActual: %d", val)
	}
}
