package main

import (
	"fmt"
	"strings"

	"github.com/jdkato/prose/v2"
)

func Display(commits []commit) (err error) {
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
	cursor := 0

	for i := range tokens {
		if tokens[i].Tag == "." {
			fmt.Println()
			for ; cursor < i; cursor++ {
				if tokens[cursor].Tag == "." {
					continue
				}
				taglen := len(tokens[cursor].Tag) + 1
				txtlen := len(tokens[cursor].Text) + 1
				tag := tokens[cursor].Tag + spaces(max(0, txtlen-taglen))
				fmt.Print(tag, " ")
			}
			fmt.Println()
			continue
		}
		fmt.Print(tokens[i].Text, " ")
	}
	return nil
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

const spaces32 = "                                "

func spaces(n int) string {
	if n < 32 {
		return spaces32[:n]
	}
	var res string
	for i := 0; i < n/32; i++ {
		res += spaces32
	}
	return res + spaces32[:n%32]
}
