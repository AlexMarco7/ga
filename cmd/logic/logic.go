package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Expression struct {
	Expressions []Expression
	Input       int
	Type        int
}

func (e Expression) Print() string {
	str := ""
	switch e.Type {
	case 1:
		{
			str += "NOT("
		}
	case 2:
		{
			str += "AND("
		}
	case 3:
		{
			str += "OR("
		}
	}

	if e.Type != 0 {
		for i := 0; i < len(e.Expressions); i++ {
			if i != 0 {
				str += ","
			}
			str += e.Expressions[i].Print()
		}
		str += ")"
	} else {
		str += fmt.Sprintf("%d", e.Input)
	}

	return str
}

type LogicRules struct {
	Total int
}

func (lr LogicRules) Create() interface{} {
	return createExpression(0, lr.Total)
}

func (lr LogicRules) Fitness(dna interface{}) float64 {
	return 0.0
}
func (lr LogicRules) Crossover(dna1 interface{}, dna2 interface{}) interface{} {
	return nil
}
func (lr LogicRules) Mutate(dna interface{}) interface{} {
	return nil
}

func (lr LogicRules) HasFinished(generation int, dna interface{}, fitness float64) bool {
	fmt.Printf("\r generation: %d | %s | fitness: %2f", generation, dna, fitness)
	return false
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println(createExpression(0, 2).Print())
	/*
		start := time.Now()

		fmt.Printf("%d", runtime.NumCPU())

		ga.Run(LogicRules{}, ga.Options{
			PopulationSize: 1000,
			MutationRate:   0.1,
		})
		elapsed := time.Since(start)
		fmt.Printf("\nTime taken: %s\n", elapsed)
	*/
}

func createExpression(depth int, totalInput int) Expression {
	leaf := depth >= 4 || rand.Float64() < 0.1
	if leaf {
		return Expression{
			Input: rand.Intn(totalInput),
			Type:  0,
		}
	} else {
		tp := 1 + rand.Intn(3)
		if tp == 1 { //not
			return Expression{
				Expressions: []Expression{createExpression(depth+1, totalInput)},
				Type:        tp,
			}
		} else { //and - or
			n := 2 + rand.Intn(4)
			ex := Expression{
				Expressions: []Expression{},
				Type:        tp,
			}
			for i := 0; i < n; i++ {
				ex.Expressions = append(ex.Expressions, createExpression(depth+1, totalInput))
			}
			return ex
		}
	}
}
