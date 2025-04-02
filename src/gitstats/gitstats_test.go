package gitstats

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/August-Brandt/EgoLottery/testutils"
)



func TestGetStats(t *testing.T) {
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

	repos = GetStats([]string{filepath.Join(repo, ".git")}, "john@doe.org", "days")
	
	if len(repos) != 1 {
		t.Errorf("Incorrect number of repos\n\tExpected: %d\n\tActual: %d", 1, len(repos))
	}
	if repos[0].Commits[0] != 1 {
		t.Errorf("Incorrect number of commits in repo\n\tExpected: %d\n\tActual: %d", 1, len(repos))
	}
}