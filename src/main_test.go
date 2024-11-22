package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func createTestFolder() (string, error) {
	name, err := os.MkdirTemp(".", "")
	if err != nil {
		return name, err
	}
	// Filling test folder with stuff so there is something to look at
	file1, err := os.Create(filepath.Join(".", name, ".gitignore"))
	if err != nil {
		return name, err
	}
	file1.Close()
	file2, err := os.Create(filepath.Join(".", name, "README.md"))
	if err != nil {
		return name, err
	}
	file2.Close()

	err = os.Mkdir(filepath.Join(".", name, ".git"), 0644)
	if err != nil {
		return name, err
	}
	err = os.Mkdir(filepath.Join(".", name, "src"), 0644)
	if err != nil {
		return name, err
	}
	return name, nil
}

func TestPrintPath(t *testing.T) {
	name, err := createTestFolder()
	if err != nil {
		t.Fatalf("Failed to create test folder: %v\n", err)
	}
	// Clean up
	defer os.RemoveAll(filepath.Join(".", name))

	
	var logbuf bytes.Buffer
	PrintPath(name, log.New(&logbuf, "", log.Lmsgprefix))

	data, err := io.ReadAll(&logbuf)
	if err != nil {
		t.Fatalf("Could not read log file: %v", err)
	}

	correctValue := fmt.Sprintf("%s:\n\t.git\n\t.gitignore\n\tREADME.md\n\tsrc\n", name)
	if string(data) != correctValue {
		t.Errorf("Incorrect data logged:\nCorrect:\n%s\nActual:\n%s\n", correctValue, string(data))
	}
}
