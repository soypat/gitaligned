package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

var (
	username = ""
)

func run() error {

	pflag.StringVarP(&username, "user", "u", "", "git username. see `git config --get user.name`")
	pflag.Parse()

	commits, _, err := ScanCWD()
	if err != nil {
		return err
	}
	return Display(commits)

}

func main() {
	if err := run(); err != nil {
		fmt.Printf("\nError in run: %s", err)
	}
}
