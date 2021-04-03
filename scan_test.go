package main

import (
	"fmt"
	"testing"
)

func TestScan(t *testing.T) {
	commits, _, err := ScanCWD()
	if len(commits) == 0 || err != nil {
		t.Errorf("zero commits or error:%v", err)
	}
	if commits[len(commits)-1].message != "Initial commit" {
		t.Error("Expected initial commit at the end")
	}
}

func TestAlignments(t *testing.T) {
	commits, authors, _ := ScanCWD()
	SetAuthorAlignments(commits, authors)
	for i := range authors {
		fmt.Printf("%+v\n", authors[i])
	}
}
