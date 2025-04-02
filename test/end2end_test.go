package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

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

	currentDir, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("Error while getting current directory: %v\n", err)
	}
	err = os.WriteFile("./testDirs", []byte(currentDir), 0644)
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

	splitOut := strings.Split(outString, "\n")
	if splitOut[0] != ".git directories found:" {
		t.Errorf("Incorrect output compared to expected\n\tActual:\n%s\n\tExpected:\n.git directories found:\n", outString)
	}
	if splitOut[1] != filepath.Join(currentDir, ".git") {
		t.Errorf("Incorrect output compared to expected\n\tActual:\n%s\n\tExpected:\n%s\n", outString, filepath.Join(currentDir, ".git"))
	}
	if len(splitOut) < 2 {
		t.Errorf("Output too short")
	}
}
