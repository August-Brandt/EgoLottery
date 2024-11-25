package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func createTestDirectory(path string) error {
	err := os.Mkdir(path, 0755)
	if err != nil {
		return err
	}
	// Filling test folder with stuff so there is something to look at
	file1, err := os.Create(filepath.Join(path, ".gitignore"))
	if err != nil {
		return err
	}
	file1.Close()

	gitCmd := exec.Command("git", "init", path)
	_, err = gitCmd.Output()
	if err != nil {
		return err
	}
	// err = os.Mkdir(filepath.Join(path, ".git"), 0644)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func buildProgram(t *testing.T) {
	cmd := exec.Command("go", "build", "-C", "../src", "-o", "../test/egolottery")
	_, err := cmd.Output()
	if err != nil {
		t.Fatalf("%v\n", err)
	}
}

func TestEndToEnd(t *testing.T) {
	buildProgram(t)
	defer os.Remove("egolottery")
	err := createTestDirectory("./testDir")
	if err != nil {
		t.Fatalf("Error while creating directory: %v\n", err)
	}
	defer os.RemoveAll("./testDir")
	absPath, err := filepath.Abs("./testDir")
	if err != nil {
		t.Fatalf("Error while getting absolute path: %v\n", err)
	}
	err = createTestDirectory(filepath.Join(absPath, "testSubDir"))
	if err != nil {
		t.Fatalf("Error while creating sub-directory: %v\n", err)
	}

	err = os.WriteFile("./testDirs", []byte(absPath), 0644)
	if err != nil {
		t.Fatalf("Error while creating testDirs file: %v\n", err)
	}
	defer os.Remove("./testDirs")

	cmd := exec.Command("./egolottery", "-file", "./testDirs", "-depth", "1")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("Error while running program: %v\n", err)
	}
	outString := string(out)
	currentDir, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("Error while getting current directory: %v\n", err)
	}
	correct := fmt.Sprintf(".git directories found:\n%s\n%s\n",
		filepath.Join(currentDir, "testDir", ".git"),
		filepath.Join(currentDir, "testDir", "testSubDir", ".git"))

	if outString != correct {
		t.Errorf("Incorrect output compared to expected\n\tActual:\n%s\n\tExpected:\n%s\n", outString, correct)
	}
}
