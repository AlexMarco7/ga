package expr

import (
	"fmt"
	"math"
	"math/rand"
)

type Operator int

const (
	EQ Operator = iota
	NOT
	AND
	OR
	NAND
	NOR
	XOR
	NXOR
)

type Expression struct {
	Expressions []Expression
	Input       int
	Operator    Operator
	Depth       int
}

func (e Expression) ToString() string {
	str := ""
	switch e.Operator {
	case NOT:
		{
			str += "NOT("
		}
	case AND:
		{
			str += "AND("
		}
	case OR:
		{
			str += "OR("
		}
	case NAND:
		{
			str += "NAND("
		}
	case NOR:
		{
			str += "NOR("
		}
	case XOR:
		{
			str += "XOR("
		}
	case NXOR:
		{
			str += "NXOR("
		}
	}

	if e.Operator != 0 {
		for i, ee := range e.Expressions {
			if i != 0 {
				str += ","
			}
			str += ee.ToString()
		}
		str += ")"
	} else {
		str += fmt.Sprintf("%s", string(65+e.Input))
	}

	return str
}

func (e Expression) Complexity() int {
	if e.Operator != EQ {
		c := 1
		for _, ee := range e.Expressions {
			c += ee.Complexity()
		}
		return c
	} else {
		return 1
	}
}

func (e Expression) Fitness(input [][]bool, output []bool, complexityRate float64) float64 {
	count := 0

	for i, input := range input {
		if e.Execute(input) == output[i] {
			count++
		}
	}

	if complexityRate == 0.0 {
		return float64(count)
	} else {
		return float64(count) + (1 / float64(e.Complexity()) * complexityRate)
	}
}

func (e Expression) Execute(input []bool) bool {
	switch e.Operator {
	case EQ:
		{
			return input[e.Input]

		}
	case NOT:
		{
			return !e.Expressions[0].Execute(input)
		}
	case AND:
		{
			b := true
			for _, ee := range e.Expressions {
				b = b && ee.Execute(input)
				if !b {
					break
				}
			}
			return b
		}
	case OR:
		{
			b := false
			for _, ee := range e.Expressions {
				b = b || ee.Execute(input)
				if b {
					break
				}
			}
			return b
		}
	case NAND:
		{
			b := true
			for _, ee := range e.Expressions {
				b = b && ee.Execute(input)
				if !b {
					break
				}
			}
			return !b
		}
	case NOR:
		{
			b := false
			for _, ee := range e.Expressions {
				b = b || ee.Execute(input)
				if b {
					break
				}
			}
			return !b
		}
	case XOR:
		{
			return !e.Expressions[0].Execute(input) != !e.Expressions[1].Execute(input)
		}
	case NXOR:
		{
			return !e.Expressions[0].Execute(input) != !e.Expressions[1].Execute(input)
		}
	}
	return false
}

func Create(depth int, maxDepth int, inputLength int) Expression {
	if depth < 1 {
		depth = 1
	}

	leaf := depth >= maxDepth-1 || random(1/float64(maxDepth))
	if leaf {
		return Expression{
			Input:    rand.Intn(inputLength),
			Operator: EQ,
			Depth:    depth,
		}
	} else {
		tp := Operator(1 + rand.Intn(7))
		if tp == NOT {
			return Expression{
				Expressions: []Expression{Create(depth+1, maxDepth, inputLength)},
				Operator:    tp,
			}
		} else {
			n := 2 + rand.Intn(inputLength)

			if tp == XOR || tp == NXOR {
				n = 2
			}

			ex := Expression{
				Expressions: []Expression{},
				Operator:    tp,
			}
			for i := 0; i < n; i++ {
				ex.Expressions = append(ex.Expressions, Create(depth+1, maxDepth, inputLength))
			}
			return (ex)
		}
	}
}

func Merge(a Expression, af float64, b Expression, bf float64) Expression {

	r := 0.5

	if af+bf > 0 {
		r = af / (af + bf)
	}

	if a.Operator == EQ && b.Operator == EQ {
		return Expression{
			Input: a.Input,
		}
	} else if a.Operator == EQ && b.Operator != EQ {
		a = b
	} else if rand.Float64() <= r/2 {
		return (a)
	} else if rand.Float64() <= (1-r)/2 {
		return (b)
	} else if rand.Float64() <= 0.2 {
		return Expression{
			Operator:    OR,
			Expressions: []Expression{a, b},
		}
	}

	newExpressions := []Expression{}

	for i := 0; i < int(math.Max(float64(len(a.Expressions)), float64(len(b.Expressions)))); i++ {
		if i < len(a.Expressions) && i < len(b.Expressions) {
			if rand.Float64() <= r {
				newExpressions = append(newExpressions, a.Expressions[i])
			} else {
				newExpressions = append(newExpressions, b.Expressions[i])
			}
		} else {
			if i < len(a.Expressions) {
				if len(newExpressions) < minItems(a.Operator) || rand.Float64() <= r {
					newExpressions = append(newExpressions, a.Expressions[i])
				}
			} else {
				if len(newExpressions) < minItems(a.Operator) || rand.Float64() <= 1-r {
					newExpressions = append(newExpressions, b.Expressions[i])
				}
			}
		}
	}

	if len(newExpressions) == 0 {
		panic(newExpressions)
	}

	return (Expression{
		Operator:    a.Operator,
		Expressions: newExpressions,
	})

}

func minItems(o Operator) int {
	switch o {
	case NOT:
		return 1
	case XOR:
		return 2
	case NXOR:
		return 2
	default:
		return 1
	}
}

func Optimize(e Expression) Expression {

	if e.Operator == NOT {
		if e.Expressions[0].Operator == NOT {
			e = e.Expressions[0].Expressions[0]
		} else {
			e.Expressions = []Expression{Optimize(e.Expressions[0])}
		}
	} else if e.Operator != EQ {
		newExpressions := []Expression{}
		eqMap := map[int]Expression{}
		for _, ee := range e.Expressions {
			ee = Optimize(ee)
			if ee.Operator == EQ {
				eqMap[ee.Input] = ee
			} else {
				newExpressions = append(newExpressions, ee)
			}
		}
		for _, v := range eqMap {
			newExpressions = append(newExpressions, v)
		}
		e.Expressions = newExpressions

		if len(e.Expressions) == 1 {
			e = e.Expressions[0]
			e.Depth--
		}
	}

	return e

}

func Mutate(e Expression, mutationRate float64, maxDepth int, inputLength int) Expression {
	if random(mutationRate) {
		switch e.Operator {
		case EQ:
			{
				if random(mutationRate) {
					e.Input = rand.Intn(inputLength)
				}
			}
		case NOT:
			{
				e.Expressions[0] = Mutate(e.Expressions[0], mutationRate, maxDepth, inputLength)
				if random(mutationRate) {
					e = e.Expressions[0]
					e.Depth--
				}
			}
		case XOR:
			fallthrough
		case NXOR:
			{
				e.Expressions[0] = Mutate(e.Expressions[0], mutationRate, maxDepth, inputLength)
				e.Expressions[1] = Mutate(e.Expressions[1], mutationRate, maxDepth, inputLength)
				if random(mutationRate) {
					e.Operator = []Operator{XOR, NXOR}[rand.Intn(2)]
				}
			}
		case AND:
			fallthrough
		case OR:
			fallthrough
		case NAND:
			fallthrough
		case NOR:
			{
				if random(mutationRate) {
					newExpressions := []Expression{}

					for _, ee := range e.Expressions {
						if !random((1 / float64(len(e.Expressions)))) {
							newExpressions = append(newExpressions, Mutate(ee, mutationRate, maxDepth, inputLength))
						}
					}
					if random((1 / float64(len(e.Expressions)))) {
						newExpressions = append(newExpressions, Create(e.Depth+1, maxDepth, inputLength))
					}
					e.Expressions = newExpressions
				} else if random(mutationRate) {
					idx := int(math.Floor(1 / float64(len(e.Expressions))))
					if idx < len(e.Expressions) && idx >= 0 {
						e = e.Expressions[idx]
						e.Depth--
					}
				} else if random(mutationRate) {
					e.Operator = []Operator{AND, OR, NAND, NOR}[rand.Intn(4)]
					newExpressions := []Expression{}
					for _, ee := range e.Expressions {
						newExpressions = append(newExpressions, Mutate(ee, mutationRate, maxDepth, inputLength))
					}
					e.Expressions = newExpressions
				}
			}
		}
		return (e)
	} else if random(mutationRate) {
		return Optimize(e)
	} else if random(mutationRate) {
		return Create(e.Depth, maxDepth, inputLength)
	} else {
		return Expression{
			Input:    rand.Intn(inputLength),
			Operator: EQ,
			Depth:    e.Depth,
		}
	}
}

func random(rate float64) bool {
	return rand.Float64() < rate
}
