package main

import (
	"fmt"
	"io"

	"github.com/jdkato/prose/v2"
)

func WriteCommitAlignments(w io.Writer, commits []commit) error {
	for i := range commits {
		_, err := w.Write([]byte(fmt.Sprintf("Commit %v\n%+0.3g\n\t%v\n\n",
			commits[i].alignment.Format(), commits[i].alignment, commits[i].message)))
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteAuthorAlignments(w io.Writer, authors []author) error {
	for i := range authors {
		_, err := w.Write([]byte(fmt.Sprintf("%v", authors[i].Stats())))
		if err != nil {
			return err
		}
		if i != len(authors)-1 {
			_, err = w.Write([]byte("\n"))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func WriteNLPTags(w io.Writer, commits []commit) (err error) {
	return walkCommits(commits, writeNLP(w))
}

func writeNLP(w io.Writer) func(*commit, []prose.Token) {
	return func(c *commit, tokens []prose.Token) {
		for i := range tokens {
			w.Write([]byte(tokens[i].Text + " "))
		}
		w.Write([]byte("\n"))
		for cursor := range tokens {
			taglen := len(tokens[cursor].Tag) + 1
			txtlen := len(tokens[cursor].Text) + 1
			tag := tokens[cursor].Tag + spaces(max(0, txtlen-taglen))
			w.Write([]byte(tag + " "))
			// fmt.Print(tag, " ")
		}
		w.Write([]byte("\n"))
	}
}
