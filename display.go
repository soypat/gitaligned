package main

import (
	"fmt"

	"github.com/jdkato/prose/v2"
)

func DisplayCommitAlignments(commits []commit) error {
	for i := range commits {
		fmt.Println("Commit " + commits[i].alignment.Format())
		fmt.Printf("%+0.3g\n", commits[i].alignment)
		fmt.Printf("\t%v\n\n", commits[i].message)
	}

	return nil
}

func DisplayAuthorAlignments(authors []author) error {
	for i := range authors {
		fmt.Printf("%v", authors[i].Stats())
		if i != len(authors)-1 {
			fmt.Println()
		}
	}
	return nil
}

func DisplayNLPTags(commits []commit) (err error) {
	return walkCommits(commits, dispNLP)
}

func dispNLP(c *commit, tokens []prose.Token) {
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
