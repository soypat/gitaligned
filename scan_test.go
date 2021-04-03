package main

import (
	"testing"
)

func TestScan(t *testing.T) {
	commits, _, err := ScanCWD()
	if len(commits) == 0 || err != nil {
		t.Errorf("zero commits or error:%v", err)
	}
	if commits[len(commits)-1].messsage != "Initial commit" {
		t.Error("Expected initial commit at the end")
	}
}

func TestAlignments(t *testing.T) {
	commits, authors, _ := ScanCWD()
	SetAlignments(commits, authors)
	t.Error("ERR")
}
