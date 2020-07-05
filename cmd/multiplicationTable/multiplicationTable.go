package main

import (
	"fmt"
	"math/rand"
	"time"

	expr "github.com/alexmarco7/ga/pkg/expression"
	"github.com/alexmarco7/ga/pkg/ga"
)

type MultiplicationTableRules struct {
	InputValues  [][]bool
	OutputValues []bool
	MutationRate float64
	MaxDepth     int
}

func (r MultiplicationTableRules) Create() interface{} {
	return expr.Create(1, r.MaxDepth, len(r.InputValues[0]))
}

func (r MultiplicationTableRules) Fitness(dna interface{}) float64 {
	e := dna.(expr.Expression)

	return e.CalcFitness(r.InputValues, r.OutputValues, 0.1)
}
func (r MultiplicationTableRules) Crossover(dna1 interface{}, dna2 interface{}) interface{} {
	e1 := dna1.(expr.Expression)
	e2 := dna2.(expr.Expression)

	return expr.Merge(e1, e2)
}
func (r MultiplicationTableRules) Mutate(dna interface{}) interface{} {
	e := dna.(expr.Expression)
	return expr.Mutate(e, r.MutationRate, r.MaxDepth, len(r.InputValues[0]))
}

func (r MultiplicationTableRules) HasFinished(generation int, dna interface{}, fitness float64) bool {
	e := dna.(expr.Expression)
	fmt.Printf("\n generation: %d | %s | fitness: %2f", generation, e.ToString(), fitness)

	finished := fitness > float64(len(r.InputValues))
	//finished := generation > 10000

	return finished
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	start := time.Now()

	inputValues := [][]bool{}
	outputValues := [5][]bool{[]bool{}, []bool{}, []bool{}, []bool{}, []bool{}}

	for i := 0; i <= 9; i++ {
		for j := 0; j <= 9; j++ {
			x := i * j
			input := []bool{
				i&16 == 16,
				i&8 == 8,
				i&4 == 4,
				i&2 == 2,
				i&1 == 1,
				j&16 == 16,
				j&8 == 8,
				j&4 == 4,
				j&2 == 2,
				j&1 == 1,
			}
			inputValues = append(inputValues, input)
			outputValues[0] = append(outputValues[0], x&16 == 16)
			outputValues[1] = append(outputValues[1], x&8 == 8)
			outputValues[2] = append(outputValues[2], x&4 == 4)
			outputValues[3] = append(outputValues[3], x&2 == 2)
			outputValues[4] = append(outputValues[4], x&1 == 1)
		}
	}

	for i, output := range outputValues {
		result := ga.Run(MultiplicationTableRules{
			InputValues:  inputValues,
			OutputValues: output,
			MutationRate: 0.5,
			MaxDepth:     10,
		}, ga.Options{
			PopulationSize: 5000,
			SurvivalRate:   0.001,
			MutationRate:   0.5,
		}).(expr.Expression)
		fmt.Printf("\n output: %d  expression: %s", i+1, expr.Optimize(result).ToString())
	}

	elapsed := time.Since(start)
	fmt.Printf("\nTime taken: %s\n", elapsed)

}
