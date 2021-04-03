package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

var (
	username   string
	maxCommits int
)

func run() error {

	pflag.StringVarP(&username, "user", "u", "", "git username. see `git config --get user.name`")
	pflag.IntVarP(&maxCommits, "max-commits", "n", 200, "max amount of commits to process")
	pflag.Parse()

	commits, _, err := ScanCWD()
	if err != nil {
		return err
	}
	return DisplayNLPTags(commits[:min(len(commits), maxCommits)])

}

func main() {
	if err := run(); err != nil {
		fmt.Printf("\nError in run: %s", err)
	}
}
