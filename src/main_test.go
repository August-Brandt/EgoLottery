package main

import (
	"github.com/August-Brandt/EgoLottery/testutils"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestPrintPath(t *testing.T) {
	name, err := testutils.CreateTestDirectory()
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
		t.Errorf("Incorrect data logged:\nExpected:\n%s\nActual:\n%s\n", correctValue, string(data))
	}
}
