package main

import (
	"fmt"
	"math/rand"
	"time"

	expr "github.com/alexmarco7/ga/pkg/expression"
	"github.com/alexmarco7/ga/pkg/ga"
)

type LogicRules struct {
	InputValues          [][]bool
	OutputValues         []bool
	BestFitness          float64
	SameBestFitnessCount int
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

	f := float64(count)

	if lr.SameBestFitnessCount >= 10 {
		lr.SameBestFitnessCount = 0
		f += (1 / float64(e.Complexity()) * 0.1)
	}

	return f
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
	//e := dna.(expr.Expression)
	fmt.Printf("\r generation: %d | %s | fitness: %2f", generation, "" /*e.ToString()*/, fitness)

	/*for i, inputs := range lr.InputValues {
		fmt.Printf("\n %v | %v | %v", inputs, e.Execute(inputs), lr.OutputValues[i])
	}*/

	if lr.BestFitness == fitness {
		lr.SameBestFitnessCount++
	}

	lr.BestFitness = fitness

	return fitness > float64(len(lr.InputValues)) //generation > 10000
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

	ga.Run(LogicRules{
		InputValues:  inputValues,
		OutputValues: outputValues,
	}, ga.Options{
		PopulationSize: 5000,
		MutationRate:   0.05,
	})
	elapsed := time.Since(start)
	fmt.Printf("\nTime taken: %s\n", elapsed)

}
