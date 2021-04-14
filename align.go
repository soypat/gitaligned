package main

import (
	"math"

	"github.com/jdkato/prose/v2"
)

const alignmentThreshold = 0.3

type alignment struct {
	// These are axes on our alignment chart.
	// that go from -1 to 1
	Licitness float64 `json:"licitness"`
	Morality  float64 `json:"morality"`
}

// Format returns human readable alignment.
// i.e. "Neutral Evil".
//
// The threshold is set by a global variable called `alignmentThreshold`
func (a alignment) Format() (format string) {
	var good, lawful, evil, chaotic bool
	good = a.Morality > alignmentThreshold
	lawful = a.Licitness > alignmentThreshold
	evil = a.Morality < -alignmentThreshold
	chaotic = a.Licitness < -alignmentThreshold
	if !evil && !good && !lawful && !chaotic {
		return "True Neutral"
	}
	switch {
	case lawful:
		format = "Lawful "
	case chaotic:
		format = "Chaotic "
	default:
		format = "Neutral "
	}
	if good {
		format += "Good"
	} else if evil {
		format += "Evil"
	} else {
		format += "Neutral"
	}
	return format
}

// SetCommitAlignments processes commits and assigns them an alignment
func SetCommitAlignments(commits []commit, authors []author) error {
	return walkCommits(commits, func(c *commit, t []prose.Token) {
		c.alignment = getAlignment(c, t)
	})
}

// SetAuthorAlignments processes authors and sets their alignment
func SetAuthorAlignments(commits []commit, authors []author) error {
	walkCommits(commits, func(c *commit, t []prose.Token) {
		c.User.Commits++
		a := getAlignment(c, t)
		c.User.accumulator.Morality += a.Morality
		c.User.accumulator.Licitness += a.Licitness
	})
	for i := range authors {
		if authors[i].Commits > 0 {
			authors[i].alignment.Licitness = math.Erf(authors[i].accumulator.Licitness)
			authors[i].alignment.Morality = math.Erf(authors[i].accumulator.Morality)
		}
	}
	return nil
}

func getAlignment(c *commit, t []prose.Token) (a alignment) {
	tlen := len(t)
	if edgeCases(t, &a) {
		return a
	}
	var adjectives, determiners, interjections int
	for i := 1; i < tlen; i++ {
		switch t[i].Tag {
		case "NN", "NNP", "NNS":
			continue // NN (noun, singular or mass) could be just about anything
		case "JJ":
			adjectives++
		case "DT", "WDT":
			determiners++
		case "UH":
			interjections++
		}
	}
	// interjections: uh, oops, ah
	a.Licitness -= float64(interjections)
	//determiners are just noise in small messages: an, a, one, my, the
	a.Morality -= float64(determiners) * 0.1 * (10 - float64(min(tlen, 10)))
	// adjectives
	a.Morality += math.Min(float64(adjectives)*0.4, 1)
	// if adjectives > 1 {
	a.Licitness -= (float64(adjectives)/float64(tlen) - 0.2) * 3
	// }

	// normalize values so that it is within alignment chart values: [-1,1]
	a.Morality = capNorm(1, a.Morality)
	a.Licitness = capNorm(1, a.Licitness)
	return
}

func capNorm(c, f float64) float64 {
	if math.Signbit(f) {
		return math.Max(-c, f)
	}
	return math.Min(c, f)
}

// edgeCases handles the edge cases of a git commit message
// without worrying much about NLP aspects of it. If it finds
// an edge case `a` should be modified accordingly.
//
// Returned bool indicates if the resulting alignment is final
// (no more processing needed).
func edgeCases(t []prose.Token, a *alignment) (finalAlignment bool) {
	tlen := len(t)
	if tlen <= 2 {
		a.Morality = -1
		return true
	}
	switch t[0].Tag {
	// first word is verb. nice to read these commits
	case "VB", "VBZ":
		a.Morality = 1
	}
	switch t[0].Text {
	// branch merging demonstrates organized development
	case "merge":
		a.Licitness = 1
	}
	return false
}
