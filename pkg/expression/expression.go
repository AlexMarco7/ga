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
	Fitness     float64
}

func (e Expression) ToString() string {
	str := ""
	switch e.Operator {
	case EQ:
		str += fmt.Sprintf("%s", string(65+e.Input))
	case NOT:
		str += "NOT("
	case AND:
		str += "AND("
	case OR:
		str += "OR("
	case NAND:
		str += "NAND("
	case NOR:
		str += "NOR("
	case XOR:
		str += "XOR("
	case NXOR:
		str += "NXOR("
	default:
		str += "*"
	}

	if e.Operator != EQ {
		for i, ee := range e.Expressions {
			if i != 0 {
				str += ","
			}
			str += ee.ToString()
		}
		str += ")"
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

func (e Expression) CalcFitness(input [][]bool, output []bool, complexityRate float64) float64 {
	count := 0

	for i, input := range input {
		if e.Execute(input) == output[i] {
			count++
		}
	}
	f := float64(count)

	if complexityRate != 0.0 {
		f = float64(count) + (1 / float64(e.Complexity()) * complexityRate)
	}

	e.Fitness = f

	return f
}

func (e Expression) Execute(input []bool) bool {

	defer func() {
		if r := recover(); r != nil {
			println("panic: ", e.ToString())
			panic(r)
		}
	}()

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
		op := []Operator{NOT, AND, OR, NOR, XOR, NOR}[rand.Intn(6)]

		n := int(math.Min(math.Max(float64(rand.Intn(inputLength)), float64(minItems(op))), float64(maxItems(op))))

		ex := Expression{
			Expressions: []Expression{},
			Operator:    op,
		}
		for i := 0; i < n; i++ {
			ex.Expressions = append(ex.Expressions, Create(depth+1, maxDepth, inputLength))
		}
		return (ex)
	}
}

func Merge(a Expression, b Expression) Expression {
	af := a.Fitness
	bf := b.Fitness

	if bf > af {
		tmp := a
		a = b
		b = tmp

	}
	r := 0.5

	if af+bf > 0 {
		r = af / (af + bf)
	}

	if a.Operator == EQ && b.Operator == EQ {
		return Expression{
			Input: a.Input,
		}
	} else if rand.Float64() <= r/2 {
		return Expression{
			Operator:    a.Operator,
			Input:       a.Input,
			Expressions: a.Expressions,
			Depth:       a.Depth,
		}
	} else if rand.Float64() <= (1-r)/2 {
		return Expression{
			Operator:    b.Operator,
			Input:       b.Input,
			Expressions: b.Expressions,
			Depth:       b.Depth,
		}
	} else if a.Operator == EQ && b.Operator != EQ {
		a = b
	} /*else if rand.Float64() <= 0.2 {
		return Expression{
			Operator:    OR,
			Expressions: []Expression{a, b},
		}
	}*/

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
				if (len(newExpressions) < minItems(a.Operator) || rand.Float64() <= r) && len(newExpressions) < maxItems(a.Operator) {
					newExpressions = append(newExpressions, a.Expressions[i])
				}
			} else {
				if (len(newExpressions) < minItems(a.Operator) || rand.Float64() <= 1-r) && len(newExpressions) < maxItems(a.Operator) {
					newExpressions = append(newExpressions, b.Expressions[i])
				}
			}
		}
	}

	if len(newExpressions) == 0 {
		panic(newExpressions)
	}

	return Optimize(Expression{
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

func maxItems(o Operator) int {
	switch o {
	case NOT:
		return 1
	case XOR:
		return 2
	case NXOR:
		return 2
	default:
		return 999999
	}
}

func Optimize(e Expression) Expression {
	ne := e

	if ne.Operator == NOT {
		if ne.Expressions[0].Operator == NOT {
			ne = Expression{
				Operator:    ne.Expressions[0].Expressions[0].Operator,
				Input:       ne.Expressions[0].Expressions[0].Input,
				Expressions: ne.Expressions[0].Expressions[0].Expressions,
				Depth:       ne.Depth,
			}
		} else {
			ne.Expressions = []Expression{Optimize(ne.Expressions[0])}
		}
	} else if ne.Operator == AND || ne.Operator == OR || ne.Operator == NAND || ne.Operator == NOR {
		newExpressions := []Expression{}
		eqMap := map[int]Expression{}
		for _, ee := range ne.Expressions {
			ee = Optimize(ee)
			if ne.Operator == EQ {
				eqMap[ne.Input] = ee
			} else {
				newExpressions = append(newExpressions, ee)
			}
		}
		for _, v := range eqMap {
			newExpressions = append(newExpressions, v)
		}
		ne.Expressions = newExpressions

		if len(ne.Expressions) == 1 {
			ne = Expression{
				Operator:    ne.Expressions[0].Operator,
				Input:       ne.Expressions[0].Input,
				Expressions: ne.Expressions[0].Expressions,
				Depth:       ne.Depth,
			}
		}
	}

	return ne

}

func Mutate(e Expression, mutationRate float64, maxDepth int, inputLength int) Expression {
	ne := e
	if random(mutationRate) {
		switch ne.Operator {
		case EQ:
			{
				if random(mutationRate) {
					ne.Input = rand.Intn(inputLength)
				}
			}
		case NOT:
			{
				if random(mutationRate) {
					ne = Expression{
						Operator:    ne.Expressions[0].Operator,
						Input:       ne.Expressions[0].Input,
						Expressions: ne.Expressions[0].Expressions,
						Depth:       ne.Depth,
					}
				} else {
					ne.Expressions = []Expression{Mutate(ne.Expressions[0], mutationRate, maxDepth, inputLength)}
				}
			}
		case XOR:
			fallthrough
		case NXOR:
			{
				ne.Expressions = []Expression{
					Mutate(ne.Expressions[0], mutationRate, maxDepth, inputLength),
					Mutate(ne.Expressions[1], mutationRate, maxDepth, inputLength),
				}
				if random(mutationRate) {
					ne.Operator = []Operator{XOR, NXOR}[rand.Intn(2)]
				}
			}
			break
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

					for _, ee := range ne.Expressions {
						if !random((1 / float64(len(ne.Expressions)))) {
							newExpressions = append(newExpressions, Mutate(ee, mutationRate, maxDepth, inputLength))
						}
					}
					if len(newExpressions) == 0 || random((1 / float64(len(newExpressions)))) {
						newExpressions = append(newExpressions, Create(ne.Depth+1, maxDepth, inputLength))
					}
					ne.Expressions = newExpressions

				} else if random(mutationRate) {
					idx := int(math.Floor(1 / float64(len(ne.Expressions))))
					if idx < len(ne.Expressions) && idx >= 0 {
						ne = Expression{
							Operator:    ne.Expressions[idx].Operator,
							Input:       ne.Expressions[idx].Input,
							Expressions: ne.Expressions[idx].Expressions,
							Depth:       ne.Depth,
						}
					}
				} else if random(mutationRate) {
					ne.Operator = []Operator{AND, OR, NAND, NOR}[rand.Intn(4)]
					newExpressions := []Expression{}
					for _, ee := range ne.Expressions {
						newExpressions = append(newExpressions, Mutate(ee, mutationRate, maxDepth, inputLength))
					}
					ne.Expressions = newExpressions
				}
			}
		}
	} else if random(mutationRate) {
		ne = Optimize(e)
	} else if random(mutationRate) {
		ne = Create(ne.Depth, maxDepth, inputLength)
	} else {
		ne = Expression{
			Input:    rand.Intn(inputLength),
			Operator: EQ,
			Depth:    ne.Depth,
		}
	}

	return ne
}

func random(rate float64) bool {
	return rand.Float64() < rate
}
