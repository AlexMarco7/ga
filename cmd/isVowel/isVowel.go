package main

import (
	"fmt"
	"math/rand"
	"time"

	expr "github.com/alexmarco7/ga/pkg/expression"
	"github.com/alexmarco7/ga/pkg/ga"
)

type IsVowelRules struct {
	InputValues  [][]bool
	OutputValues []bool
	MutationRate float64
	MaxDepth     int
}

func (r IsVowelRules) Create() interface{} {
	return expr.Create(1, r.MaxDepth, len(r.InputValues[0]))
}

func (r IsVowelRules) Fitness(dna interface{}) float64 {
	e := dna.(expr.Expression)

	return e.Fitness(r.InputValues, r.OutputValues, 0.1)
}
func (r IsVowelRules) Crossover(dna1 interface{}, f1 float64, dna2 interface{}, f2 float64) interface{} {
	e1 := dna1.(expr.Expression)
	e2 := dna2.(expr.Expression)

	return expr.Merge(e1, f1, e2, f2)
}
func (r IsVowelRules) Mutate(dna interface{}) interface{} {
	e := dna.(expr.Expression)
	return expr.Mutate(e, r.MutationRate, r.MaxDepth, len(r.InputValues[0]))
}

func (r IsVowelRules) HasFinished(generation int, dna interface{}, fitness float64) bool {
	e := dna.(expr.Expression)
	fmt.Printf("\n generation: %d | %s | fitness: %2f", generation, "" /*e.ToString()*/, fitness)

	finished := fitness > float64(len(r.InputValues))
	//finished := generation > 10000

	if finished {
		fmt.Printf("\n %s", expr.Optimize(e).ToString())
	}

	return finished
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	start := time.Now()

	inputStr := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}
	inputValues := [][]bool{}
	outputValues := []bool{}

	for _, c := range inputStr {
		bt := []byte(c)[0]
		input := []bool{
			bt&1 == 1,
			bt&2 == 2,
			bt&4 == 4,
			bt&8 == 8,
			bt&16 == 16,
			bt&32 == 32,
			bt&64 == 64,
			bt&128 == 128,
		}
		inputValues = append(inputValues, input)
		outputValues = append(outputValues, (c == "A" || c == "E" || c == "I" || c == "O" || c == "U" || c == "a" || c == "e" || c == "i" || c == "o" || c == "u"))
	}

	ga.Run(IsVowelRules{
		InputValues:  inputValues,
		OutputValues: outputValues,
		MutationRate: 0.5,
		MaxDepth:     6,
	}, ga.Options{
		PopulationSize: 5000,
		SurvivalRate:   0.001,
		MutationRate:   0.5,
	})
	elapsed := time.Since(start)
	fmt.Printf("\nTime taken: %s\n", elapsed)

}
