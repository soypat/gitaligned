package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

type initiator func()

func setMaxCommits(i int) initiator { return func() { maxCommits = i } }

func setMaxAuthors(i int) initiator { return func() { maxAuthors = i } }

func setBranch(b string) initiator { return func() { branch = b } }

func initTest(is ...func()) {
	maxCommits = 200
	maxAuthors = 20
	branch = ""
	for f := range is {
		is[f]()
	}
}
func TestScan(t *testing.T) {
	initTest()
	commits, authors, err := ScanCWD(branch)
	if len(commits) == 0 || err != nil {
		t.Errorf("zero commits or error:%v", err)
		t.FailNow()
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
		{maxCommits: 315, maxAuthors: 20, file: "testdata/deal.log"},
		{maxCommits: 2, maxAuthors: 20, file: "testdata/deal.log"},
		{maxCommits: 20, maxAuthors: 20, file: "testdata/deal.log"},
		{maxCommits: 80, maxAuthors: 20, file: "testdata/deal.log"},
	}
	for i := range tests {
		initTest(setMaxCommits(tests[i].maxCommits), setMaxAuthors(tests[i].maxAuthors))
		commits, authors := scanFromLogFile(t, tests[i].file)
		if len(commits) == 0 || len(commits) < maxCommits || len(authors) == 0 {
			t.Error("bad length of results " + tests[i].file)
		}
		err := SetAuthorAlignments(commits, authors)
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
		{maxCommits: 200, maxAuthors: 100, file: "testdata/awesomego.log"},
		{maxCommits: 315, maxAuthors: 20, file: "testdata/deal.log"},
		{maxCommits: 2, maxAuthors: 20, file: "testdata/deal.log"},
		{maxCommits: 20, maxAuthors: 20, file: "testdata/deal.log"},
		{maxCommits: 80, maxAuthors: 20, file: "testdata/deal.log"},
	}
	for i := range tests {
		initTest(setMaxCommits(tests[i].maxCommits), setMaxAuthors(tests[i].maxAuthors))
		commits, _ := scanFromLogFile(t, tests[i].file)
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

func TestDisplay(t *testing.T) {
	var tests = []struct {
		file       string
		maxCommits int
		maxAuthors int
	}{
		{maxCommits: 200, maxAuthors: 100, file: "testdata/awesomego.log"},
		{maxCommits: 315, maxAuthors: 20, file: "testdata/deal.log"},
		{maxCommits: 2, maxAuthors: 20, file: "testdata/deal.log"},
		{maxCommits: 20, maxAuthors: 20, file: "testdata/deal.log"},
		{maxCommits: 80, maxAuthors: 20, file: "testdata/deal.log"},
	}
	for i := range tests {
		initTest(setMaxCommits(tests[i].maxCommits), setMaxAuthors(tests[i].maxAuthors))
		commits, authors := scanFromLogFile(t, tests[i].file)
		err := WriteAuthorAlignments(io.Discard, authors)
		if err != nil {
			t.Error(err)
		}
		err = SetCommitAlignments(commits, authors)
		if err != nil {
			t.Error(err)
		}
		err = WriteCommitAlignments(io.Discard, commits)
		if err != nil {
			t.Error(err)
		}
		err = WriteNLPTags(io.Discard, commits)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestNoFolderError(t *testing.T) {
	err := os.Chdir("../")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	commits, _, err := ScanCWD("")
	if err == nil || len(commits) > 0 {
		t.Error("expected error being in non git dir")
	}
}
func scanFromLogFile(t *testing.T, filename string) ([]commit, []author) {
	f, err := os.Open(filename)
	if err != nil {
		t.FailNow()
	}
	t.Cleanup(func() { f.Close() })

	commits, authors, err := GitLogScan(f)
	if err != nil && err != io.EOF {
		t.Error(err)
		t.FailNow()
	}
	return commits, authors
}
