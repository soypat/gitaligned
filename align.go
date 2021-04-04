package main

import (
	"math"

	"github.com/jdkato/prose/v2"
)

const alignmentThreshold = 0.3

type alignment struct {
	// These are axes on our alignment chart.
	// that go from -1 to 1
	ChaoticLaw, Morality float64
}

func (a alignment) Format() (format string) {
	var good, lawful, evil, chaotic bool
	good = a.Morality > alignmentThreshold
	lawful = a.ChaoticLaw > alignmentThreshold
	evil = a.Morality < -alignmentThreshold
	chaotic = a.ChaoticLaw < -alignmentThreshold
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

func SetCommitAlignments(commits []commit, authors []author) error {
	return walkCommits(commits, func(c *commit, t []prose.Token) {
		c.alignment = getAlignment(c, t)
	})
}

func SetAuthorAlignments(commits []commit, authors []author) error {
	walkCommits(commits, func(c *commit, t []prose.Token) {
		c.user.commits++
		a := getAlignment(c, t)
		c.user.accumulator.Morality += a.Morality
		c.user.accumulator.ChaoticLaw += a.ChaoticLaw
	})
	for i := range authors {
		if authors[i].commits > 0 {
			authors[i].alignment.ChaoticLaw = math.Erf(authors[i].accumulator.ChaoticLaw)
			authors[i].alignment.Morality = math.Erf(authors[i].accumulator.Morality)
		}
	}
	return nil
}

func getAlignment(c *commit, t []prose.Token) (a alignment) {
	tlen := len(t)
	if tlen <= 2 {
		a.Morality = -1
		return
	}
	// first word is verb. nice to read these commits
	if t[0].Tag == "VB" || t[0].Tag == "VBZ" {
		a.Morality = 1
	}
	adjectives := 0
	for i := 1; i < tlen; i++ {
		switch t[i].Tag {
		case "NN":
			continue // NN (noun, singular or mass) could be just about anything
		case "JJ": // adjective
			adjectives++
		case "DT":
			// determiners are just noise in small messages
			a.Morality -= 0.1 * (10 - float64(min(tlen, 10)))
		}
	}
	a.Morality += math.Min(float64(adjectives)*0.4, 1)
	a.ChaoticLaw -= (float64(adjectives)/float64(tlen) - 0.1) * 3
	a.Morality = capNorm(1, a.Morality)
	a.ChaoticLaw = capNorm(1, a.ChaoticLaw)
	return
}

func capNorm(c, f float64) float64 {
	if math.Signbit(f) {
		return math.Max(-c, f)
	}
	return math.Min(c, f)
}
