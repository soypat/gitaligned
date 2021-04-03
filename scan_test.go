package gitscan

import (
	"testing"
)

func TestScan(t *testing.T) {
	commits, err := ScanCWD()
	if len(commits) == 0 || err != nil {
		t.Errorf("zero commits or error:%v", err)
	}
	if commits[len(commits)-1].messsage != "Initial commit" {
		t.Error("Expected initial commit at the end")
	}
}
