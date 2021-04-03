package main

import (
	"fmt"

	"github.com/jdkato/prose/v2"
)

type alignment struct {
	// These are axes on our alignment chart.
	// that go from -1 to 1
	Morality, ChaoticLaw float64
}

func SetAlignments(commits []commit, authors []author) error {
	for i := range commits {
		// alg := alignment{}
		doc, err := prose.NewDocument(commits[i].messsage)
		if err != nil {
			return err
		}
		fmt.Println("TOK:")
		for _, t := range doc.Tokens() {
			fmt.Printf("%+v\n", t)
		}
		fmt.Println("ENT:")
		for _, t := range doc.Entities() {
			fmt.Printf("%+v\n", t)
		}
	}
	return nil
}
