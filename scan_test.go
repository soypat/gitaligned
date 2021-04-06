package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func initTest() {
	maxCommits = 200
	maxAuthors = 20
	branch = "main"
}
func TestScan(t *testing.T) {
	initTest()
	commits, authors, err := ScanCWD(branch)
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
	commits, authors, _ := ScanCWD(branch)
	SetAuthorAlignments(commits, authors)
	emptyalignment := alignment{}
	for i := range authors {
		if authors[i].alignment == emptyalignment {
			t.Error("alignment not set for", authors[i].name)
		}
	}
}

func TestFile(t *testing.T) {
	var tests = []struct {
		file       string
		maxCommits int
		maxAuthors int
	}{
		{maxCommits: 200, maxAuthors: 100, file: "testdata/awesomego.log"},
		{maxCommits: 315, file: "testdata/deal.log"},
		{maxCommits: 2, file: "testdata/deal.log"},
		{maxCommits: 20, file: "testdata/deal.log"},
		{maxCommits: 80, file: "testdata/deal.log"},
	}
	for i := range tests {
		initTest()
		f, err := os.Open(tests[i].file)
		if err != nil {
			t.FailNow()
		}
		defer f.Close()
		if tests[i].maxCommits != 0 {
			maxCommits = tests[i].maxCommits
		}
		if tests[i].maxAuthors != 0 {
			maxAuthors = tests[i].maxAuthors
		}
		commits, authors, err := GitLogScan(f)
		if err != nil && err != io.EOF {
			t.Error(err)
		}
		if len(commits) == 0 || len(commits) < maxCommits || len(authors) == 0 {
			t.Error("bad length of results " + tests[i].file)
		}
		err = SetAuthorAlignments(commits, authors)
		if err != nil {
			t.Error(err)
		}
		emptyalignment := alignment{}
		for i := range authors {
			if authors[i].alignment == emptyalignment {
				t.Errorf("%d:alignment not set for %v", maxCommits, authors[i].Stats())
			}
		}
	}
}

func TestTokenize(t *testing.T) {
	var tests = []struct {
		file       string
		maxCommits int
		maxAuthors int
	}{
		{maxCommits: 300, maxAuthors: 100, file: "testdata/awesomego.log"},
		{maxCommits: 315, file: "testdata/deal.log"},
		{maxCommits: 2, file: "testdata/deal.log"},
		{maxCommits: 20, file: "testdata/deal.log"},
		{maxCommits: 80, file: "testdata/deal.log"},
	}
	for i := range tests {
		initTest()
		f, err := os.Open(tests[i].file)
		if err != nil {
			t.FailNow()
		}
		defer f.Close()
		if tests[i].maxCommits != 0 {
			maxCommits = tests[i].maxCommits
		}
		if tests[i].maxAuthors != 0 {
			maxAuthors = tests[i].maxAuthors
		}
		commits, _, err := GitLogScan(f)
		if err != nil && err != io.EOF {
			t.Error(err)
		}
		tokens, err := tokenizeCommits(commits)
		if err != nil {
			t.Error(err)
		}
		atCommit := 0
		last := -1
		for i := range tokens {
			if tokens[i].Tag == "." {
				// f(&commits[atCommit], tokens[last+1:i])
				msg := replacecommits.Replace(commits[atCommit].message)
				splitsies := strings.Fields(msg)
				for j := range splitsies {
					if splitsies[j] != "" && splitsies[j] != tokens[last+1+j].Text {
						t.Errorf("expected %q=%q\n in %v=%v\n\n", splitsies[j], tokens[last+1+j].Text, splitsies, tokens[last+1:i])
					}
				}
				last = i
				atCommit++
			}
		}
	}
}
