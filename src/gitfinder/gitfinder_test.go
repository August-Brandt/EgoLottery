package gitfinder

import (
	"EgoLottery/testutils"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestFindGitReposNDepth1Path(t *testing.T) {
	// Setup
	root, err := testutils.CreateTestDirectory()
	if err != nil {
		t.Fatalf("Error while creating root directory: %v\n", err)
	}
	defer os.RemoveAll(root) // Clean up
	absRoot, err := filepath.Abs(root)
	if err != nil {
		t.Fatalf("Error while getting absolute path: %v\n", err)
	}
	subrepo1, err := testutils.CreateTestDirectoryFromPath(absRoot)
	if err != nil {
		t.Fatalf("Error while creating subrepo1 directory: %v\n", err)
	}
	subrepo2, err := testutils.CreateTestDirectoryFromPath(absRoot)
	if err != nil {
		t.Fatalf("Error while creating subrepo1 directory: %v\n", err)
	}
	subDir, err := os.MkdirTemp(absRoot, "test_")
	if err != nil {
		t.Fatalf("Error while creating subDir directory: %v\n", err)
	}
	subsubrepo1, err := testutils.CreateTestDirectoryFromPath(subDir)
	if err != nil {
		t.Fatalf("Error while creating subrepo1 directory: %v\n", err)
	}

	// Create tests
	var tests = []struct {
		path []string
		depth int
		correct []string
	}{
		{
			[]string{absRoot}, 
			0,
			[]string{
				filepath.Join(absRoot, ".git"),
			},
		},
		{
			[]string{absRoot}, 
			1,
			[]string{
				filepath.Join(absRoot, ".git"),
				filepath.Join(subrepo1, ".git"),
				filepath.Join(subrepo2, ".git"),
			},
		},
		{
			[]string{absRoot}, 
			2,
			[]string{
				filepath.Join(absRoot, ".git"),
				filepath.Join(subrepo1, ".git"),
				filepath.Join(subrepo2, ".git"),
				filepath.Join(subsubrepo1, ".git"),
			},
		},
	}

	for i, test := range tests {
		testname := fmt.Sprintf("test %d: depth: %d", i+1, test.depth)

		t.Run(testname, func(t *testing.T){
			dirs := FindGitRepos(test.path, test.depth)
			if len(test.correct) != len(dirs) {
				t.Errorf("Length of returned values not the same as correct\n\tExpected: %d\n\tActual: %d\n", len(test.correct), len(dirs))
			}
			sort.Strings(dirs)
			sort.Strings(test.correct)
			if !t.Failed() { // Only run if both slices are same length
				for i := 0; i < len(test.correct); i++ {
					if test.correct[i] != dirs[i] {
						t.Errorf("Actual value not the same as expected:\n\tExpected: %s\n\tActual: %s", test.correct[i], dirs[i])
					}
				}
			}
		})
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
