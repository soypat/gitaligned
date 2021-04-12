package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type initiator func()

func defaultOptions() []gitOption {
	maxAuthors = defaultMaxAuthors
	maxCommits = defaultMaxCommits
	var opts = []gitOption{
		optionMaxCommits(defaultMaxCommits),
		optionBranch(""),
	}
	return opts

}
func TestScan(t *testing.T) {
	opts := defaultOptions()
	commits, authors, err := ScanCWD(opts...)
	if len(commits) == 0 || err != nil {
		t.Errorf("zero commits or error:%v", err)
		t.FailNow()
	}
	if commits[len(commits)-1].message != "initial commit" {
		t.Error("Expected initial commit at the end")
	}
	_ = authors
}

func TestAlignments(t *testing.T) {
	opts := defaultOptions()
	commits, authors, _ := ScanCWD(opts...)
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
		maxAuthors, maxCommits = tests[i].maxAuthors, tests[i].maxCommits
		commits, authors := scanFromLogFile(t, tests[i].file)
		if len(commits) == 0 || len(commits) < tests[i].maxCommits || len(authors) == 0 {
			t.Error("bad length of results " + tests[i].file)
		}
		err := SetAuthorAlignments(commits, authors)
		if err != nil {
			t.Error(err)
		}
		emptyalignment := alignment{}
		for i := range authors {
			if authors[i].alignment == emptyalignment {
				t.Errorf("%d:alignment not set for %v", tests[i].maxCommits, authors[i].Stats())
			}
		}
	}
}

func TestTokenize(t *testing.T) {
	var tests = []struct {
		file                   string
		maxAuthors, maxCommits int
	}{
		{maxCommits: 200, maxAuthors: 100, file: "testdata/awesomego.log"},
		{maxCommits: 315, maxAuthors: 20, file: "testdata/deal.log"},
		{maxCommits: 2, maxAuthors: 20, file: "testdata/deal.log"},
		{maxCommits: 20, maxAuthors: 20, file: "testdata/deal.log"},
		{maxCommits: 80, maxAuthors: 20, file: "testdata/deal.log"},
	}
	for i := range tests {
		maxAuthors, maxCommits = tests[i].maxAuthors, tests[i].maxCommits
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
		maxAuthors = tests[i].maxAuthors
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

func TestNoFolderError(t *testing.T) {
	opts := defaultOptions()
	dir, err := filepath.Abs(".")
	if err != nil {
		t.FailNow()
	}
	base := filepath.Base(dir)
	err = os.Chdir("../")
	t.Cleanup(func() {
		os.Chdir(base)
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	commits, _, err := ScanCWD(opts...)
	if err == nil || len(commits) > 0 {
		t.Error("expected error and no commits detected being in non git dir")
	}
}

func TestFindAuthor(t *testing.T) {
	firstCommiter := "Patricio Whittingslow"
	opts := []gitOption{
		optionAuthorPattern(firstCommiter),
		optionMaxCommits(4),
	}
	maxAuthors = 1

	commits, authors, err := ScanCWD(opts...)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(authors) != 1 || len(commits) == 0 {
		t.Error("got more than 1 author or no commits")
	}
	if authors[0].name != firstCommiter {
		t.Error("could not find", firstCommiter, "author")
	}
}

func TestRun(t *testing.T) {
	err := run()
	if err != nil {
		t.Error(err)
	}
}
