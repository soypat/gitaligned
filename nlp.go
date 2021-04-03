package main

import (
	"fmt"
	"strings"

	"github.com/jdkato/prose/v2"
)

func Display(commits []commit) (err error) {
	return walkCommits(commits, disp)
}

func disp(c *commit, tokens []prose.Token) {
	for i := range tokens {
		fmt.Print(tokens[i].Text, " ")
	}
	fmt.Println()
	for cursor := range tokens {
		taglen := len(tokens[cursor].Tag) + 1
		txtlen := len(tokens[cursor].Text) + 1
		tag := tokens[cursor].Tag + spaces(max(0, txtlen-taglen))
		fmt.Print(tag, " ")
	}
	fmt.Println()
}
func walkCommits(commits []commit, f func(*commit, []prose.Token)) error {
	var err error
	if len(commits) == 0 {
		panic("expected non-nil/non-zero number of commits")
	}
	var doc *prose.Document
	var allCommits = &strings.Builder{}
	cap := allCommits.Cap()
	if cap < len(commits)*20 {
		allCommits.Grow(len(commits)*20 - cap)
	}

	for i := range commits {
		allCommits.WriteString(strings.ReplaceAll(commits[i].messsage, ".", ",") + ". ")
	}
	doc, err = prose.NewDocument(allCommits.String())
	if err != nil {
		return err
	}
	tokens := doc.Tokens()
	atCommit := 0
	last := 0
	for i := range tokens {
		if tokens[i].Tag == "." {
			f(&commits[atCommit], tokens[last+1:i])
			last = i
		}
	}
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func spaces(n int) string {
	const spaces32 = "                                "
	if n < 32 {
		return spaces32[:n]
	}
	var res string
	for i := 0; i < n/32; i++ {
		res += spaces32
	}
	return res + spaces32[:n%32]
}
