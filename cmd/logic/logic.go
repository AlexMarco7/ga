package main

import (
	"fmt"
	"math/rand"
	"time"

	expr "github.com/alexMarco7/ga/pkg/expression"
	"github.com/alexMarco7/ga/pkg/ga"
)

type LogicRules struct {
	InputValues  [][]bool
	OutputValues []bool
}

func (lr LogicRules) Create() interface{} {
	return expr.CreateExpression(4, len(lr.InputValues[0]))
}

func (lr LogicRules) Fitness(dna interface{}) float64 {
	e := dna.(expr.Expression)

	count := 0

	for i, inputs := range lr.InputValues {
		if e.Execute(inputs) == lr.OutputValues[i] {
			count++
		}
	}

	return float64(count) + (1 / float64(e.Complexity()) * 0.1)
}
func (lr LogicRules) Crossover(dna1 interface{}, f1 float64, dna2 interface{}, f2 float64) interface{} {
	e1 := dna1.(expr.Expression)
	e2 := dna2.(expr.Expression)

	return expr.MergeExpression(e1, f1, e2, f2)
}
func (lr LogicRules) Mutate(dna interface{}) interface{} {
	e := dna.(expr.Expression)
	return expr.MutateExpression(e, len(lr.InputValues[0]))
}

func (lr LogicRules) HasFinished(generation int, dna interface{}, fitness float64) bool {
	e := dna.(expr.Expression)
	fmt.Printf("\n generation: %d | %s | fitness: %2f", generation, e.ToString(), fitness)
	return generation > 10000
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	start := time.Now()

	ga.Run(LogicRules{
		InputValues: [][]bool{
			{false, false},
			{false, true},
			{true, false},
			{true, true},
		},
		OutputValues: []bool{
			false,
			false,
			false,
			true,
		},
	}, ga.Options{
		PopulationSize: 3000,
		MutationRate:   0.05,
	})
	elapsed := time.Since(start)
	fmt.Printf("\nTime taken: %s\n", elapsed)

}
