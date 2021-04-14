package main

import (
	"strings"

	"github.com/jdkato/prose/v2"
)

var replacecommits = strings.NewReplacer(".", " ", "(", " ", ")", " ", ":", " ", ",", " ",
	"[", " ", "]", " ", "\\", " ", `"`, " ", "'", " ", "!", " ", ";", " ", "?", " ",
	"/", " ", "<", " ", ">", " ")

func tokenizeCommits(commits []commit) ([]prose.Token, error) {
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
		msg := replacecommits.Replace(commits[i].Message)
		if i != len(commits)-1 {
			msg += " . "
		}
		allCommits.WriteString(msg)
	}
	// allstr :=allCommits.String()  // debugging purposes

	doc, err = prose.NewDocument(allCommits.String(),
		prose.WithExtraction(false), prose.WithSegmentation(false), prose.WithTokenization(false))
	return doc.Tokens(), err
}

// walkCommits is SLOW. This is because it processes all commit messages into one
//
func walkCommits(commits []commit, f func(*commit, []prose.Token)) error {
	tokens, err := tokenizeCommits(commits)
	if err != nil {
		return err
	}
	atCommit := 0
	last := -1
	for i := range tokens {
		if tokens[i].Tag == "." {
			f(&commits[atCommit], tokens[last+1:i])
			last = i
			atCommit++
		}
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
