package gitfinder

import (
	"EgoLottery/testutils"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestFindGitRepos2Paths(t *testing.T) {
	// Setup

	// Directory 1
	root1, err := testutils.CreateTestDirectory()
	if err != nil {
		t.Fatalf("Error while creating root1 directory: %v\n", err)
	}
	defer os.RemoveAll(root1)
	absRoot1, err := filepath.Abs(root1)
	if err != nil {
		t.Fatalf("Error while getting absolute path: %v\n", err)
	}
	subrepo1, err := testutils.CreateTestDirectoryFromPath(absRoot1)
	if err != nil {
		t.Fatalf("Error while creating subrepo1 directory: %v\n", err)
	}
	subDir1, err := os.MkdirTemp(absRoot1, "test_")
	if err != nil {
		t.Fatalf("Error while creating subDir directory: %v\n", err)
	}
	subsubrepo1, err := testutils.CreateTestDirectoryFromPath(subDir1)
	if err != nil {
		t.Fatalf("Error while creating subrepo1 directory: %v\n", err)
	}

	// Directory 2
	root2, err := testutils.CreateTestDirectory()
	if err != nil {
		t.Fatalf("Error while creating root1 directory: %v\n", err)
	}
	defer os.RemoveAll(root2)
	absRoot2, err := filepath.Abs(root2)
	if err != nil {
		t.Fatalf("Error while getting absolute path: %v\n", err)
	}
	subrepo2, err := testutils.CreateTestDirectoryFromPath(absRoot2)
	if err != nil {
		t.Fatalf("Error while creating subrepo1 directory: %v\n", err)
	}
	subDir2, err := os.MkdirTemp(absRoot1, "test_")
	if err != nil {
		t.Fatalf("Error while creating subDir directory: %v\n", err)
	}
	subsubrepo2, err := testutils.CreateTestDirectoryFromPath(subDir2)
	if err != nil {
		t.Fatalf("Error while creating subrepo1 directory: %v\n", err)
	}	

	// Create sub-tests
	var tests = []struct {
		path []string
		depth int
		correct []string
	}{
		{
			[]string{absRoot1, absRoot2}, 
			0,
			[]string{
				filepath.Join(absRoot1, ".git"),
				filepath.Join(absRoot2, ".git"),
			},
		},
		{
			[]string{absRoot1, absRoot2}, 
			1,
			[]string{
				filepath.Join(absRoot1, ".git"),
				filepath.Join(subrepo1, ".git"),
				filepath.Join(absRoot2, ".git"),
				filepath.Join(subrepo2, ".git"),
			},
		},
		{
			[]string{absRoot1, absRoot2}, 
			2,
			[]string{
				filepath.Join(absRoot1, ".git"),
				filepath.Join(subrepo1, ".git"),
				filepath.Join(subsubrepo1, ".git"),
				filepath.Join(absRoot2, ".git"),
				filepath.Join(subrepo2, ".git"),
				filepath.Join(subsubrepo2, ".git"),
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

func TestFindGitRepos1Path(t *testing.T) {
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
