package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

const defaultMaxCommits = 1200
const defaultMaxAuthors = 30

var maxAuthors, maxCommits int

func run() (err error) {
	var (
		username    string
		why         bool
		showNLPTags bool
		branch      string
		noMerges    bool
	)
	pflag.StringVarP(&username, "user", "u", "", "git username. recieves `<pattern>`")
	pflag.IntVarP(&maxCommits, "max-commits", "n", defaultMaxCommits, "max amount of commits to process")
	pflag.IntVarP(&maxAuthors, "max-authors", "a", defaultMaxAuthors, "max amount of authors to process")
	pflag.BoolVarP(&why, "why", "y", false, "print alignments and message for each commit")

	pflag.BoolVarP(&noMerges, "no-merges", "", false, "do not process commits with more than one parent")
	pflag.BoolVarP(&showNLPTags, "show-nlp", "k", false, "shows natural language processing tags detected for each commit")
	pflag.StringVarP(&branch, "branch", "b", "", "git branch to scan")
	pflag.Parse()

	options := []gitOption{
		optionNoMerges(noMerges),
		optionAuthorPattern(username),
		optionBranch(branch),
	}
	// if username is specified we know all commits returned by git log will be processed, else undefined
	if username != "" {
		options = append(options, optionMaxCommits(maxCommits))
	}
	var authors []author
	var commits []commit
	commits, authors, err = ScanCWD(options...)
	if err != nil {
		return err
	}
	if len(commits) == 0 {
		return errors.New("no commits found. Are you sure username is correct? Run `git log` to see usernames")
	}

	if why {
		SetCommitAlignments(commits, authors)
		return WriteCommitAlignments(os.Stdout, commits)
	}
	if showNLPTags {
		return WriteNLPTags(os.Stdout, commits)
	}
	SetAuthorAlignments(commits, authors)
	return WriteAuthorAlignments(os.Stdout, authors)
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error in run: %s\n", err)
	}
}
