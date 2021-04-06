package main

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
)

var (
	username    string
	maxCommits  int
	maxAuthors  int
	why         bool
	showNLPTags bool
	branch      string
)

func run() (err error) {

	pflag.StringVarP(&username, "user", "u", "", "git username. see 'git config --get user.name'")
	pflag.IntVarP(&maxCommits, "max-commits", "n", 200, "max amount of commits to process")
	pflag.IntVarP(&maxAuthors, "max-authors", "a", 20, "max amount of authors to process")
	pflag.BoolVarP(&why, "why", "y", false, "print alignments and message for each commit")
	pflag.BoolVarP(&showNLPTags, "show-nlp", "k", false, "shows natural language processing tags detected for each commit")
	pflag.StringVarP(&branch, "branch", "b", "", "git branch to scan")
	pflag.Parse()
	var authors []author
	var commits []commit
	commits, authors, err = ScanCWD(branch)
	if err != nil {
		return err
	}
	if len(commits) == 0 {
		return errors.New("no commits found. Are you sure username is correct? Run `git log` to see usernames")
	}

	if why {
		SetCommitAlignments(commits, authors)
		return DisplayCommitAlignments(commits)
	}
	if showNLPTags {
		return DisplayNLPTags(commits)
	}
	SetAuthorAlignments(commits, authors)
	return DisplayAuthorAlignments(authors)
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error in run: %s\n", err)
	}
}
