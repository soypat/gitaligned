package main

import (
	"io"
	"os"
	"testing"
)

func initTest() {
	maxCommits = 200
}
func TestScan(t *testing.T) {
	initTest()
	commits, authors, err := ScanCWD()
	if len(commits) == 0 || err != nil {
		t.Errorf("zero commits or error:%v", err)
	}
	if commits[len(commits)-1].message != "Initial commit" {
		t.Error("Expected initial commit at the end")
	}
	_ = authors
}

func TestAlignments(t *testing.T) {
	initTest()
	commits, authors, _ := ScanCWD()
	SetAuthorAlignments(commits, authors)
	for i := range authors {
		if authors[i].alignment.Morality == 0 {
			t.Error("alignment not set")
		}
	}
}

func TestFile(t *testing.T) {
	initTest()
	var tests = []struct {
		maxCommits int
	}{
		{maxCommits: 1},
		{maxCommits: 20},
		{maxCommits: 80},
	}
	for i := range tests {
		f, err := os.Open("testdata/deal.log")
		if err != nil {
			t.FailNow()
		}
		defer f.Close()
		maxCommits = tests[i].maxCommits
		commits, authors, err := GitLogScan(f)
		if err != nil && err != io.EOF {
			t.Error(err)
		}
		if len(commits) == 0 || len(commits) < maxCommits || len(authors) == 0 {
			t.Error("bad length of results")
		}
		err = SetAuthorAlignments(commits, authors)
		if err != nil {
			t.Error(err)
		}
		emptyalignment := alignment{}
		for i := range authors {
			if authors[i].alignment == emptyalignment {
				t.Errorf("alignment not set for %v", authors[i].Stats())
			}
		}
	}

}
