package gitfinder

import (
	"os"
	"path/filepath"
	"testing"
	"EgoLottery/testutils"
)

func TestFindGitFolders0Depth(t *testing.T) {
	name, err := testutils.CreateTestFolder()
	if err != nil {
		t.Fatalf("Error while creating test folder: %v\n", err)
	}
	defer os.RemoveAll(name) // Clean up
	absName, err := filepath.Abs(name)
	
	dirs := FindGitFolders([]string{absName}, 0)
	if err != nil {
		t.Fatalf("Error while finding absolute path to test folder: %v\n", err)
	}
	correctValue := []string{filepath.Join(absName, ".git")}
	if len(correctValue) != len(dirs) {
		t.Errorf("Length of returned values not the same as correct\n\tExpected: %d\n\tActual: %d\n", len(correctValue), len(dirs))
	}
	if !t.Failed() { // Only run if both slices are same length
		for i := 0; i < len(correctValue); i++ {
			if correctValue[i] != dirs[i] {
				t.Errorf("Actual value not the same as expected:\n\tExpected: %s\n\tActual: %s", correctValue[i], dirs[i])
			}
		}
	}
}

func TestPathIsDir(t *testing.T) {
	if !pathIsDir(".") {
		t.Error("Gave false when path to directory")
	}
	if pathIsDir(filepath.Join(".", "gitfinder_test.go")) {
		t.Error("Gave true when path is file")
	}
}
