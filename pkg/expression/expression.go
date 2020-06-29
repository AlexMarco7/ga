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
)

type Expression struct {
	Expressions []Expression
	Input       int
	Operator    Operator
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
		str += fmt.Sprintf("%d", e.Input)
	}

	return str
}

func (e Expression) Complexity() int {
	if e.Operator != 0 {
		c := 1
		for _, ee := range e.Expressions {
			c += ee.Complexity()
		}
		return c
	} else {
		return 1
	}
}

func (e Expression) Execute(inputs []bool) bool {
	switch e.Operator {
	case EQ:
		{
			return inputs[e.Input]

		}
	case NOT:
		{
			return !e.Expressions[0].Execute(inputs)
		}
	case AND:
		{
			b := true
			for _, ee := range e.Expressions {
				b = b && ee.Execute(inputs)
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
				b = b || ee.Execute(inputs)
				if b {
					break
				}
			}
			return b
		}
	}
	return false
}

func CreateExpression(depth int, inputLength int) Expression {
	leaf := depth == 0 || rand.Float64() < 0.1
	if leaf {
		return Expression{
			Input:    rand.Intn(inputLength),
			Operator: 0,
		}
	} else {
		tp := Operator(1 + rand.Intn(3))
		if tp == NOT {
			return Expression{
				Expressions: []Expression{CreateExpression(depth-1, inputLength)},
				Operator:    tp,
			}
		} else { //and - or
			n := 2 + rand.Intn(4)
			ex := Expression{
				Expressions: []Expression{},
				Operator:    tp,
			}
			for i := 0; i < n; i++ {
				ex.Expressions = append(ex.Expressions, CreateExpression(depth-1, inputLength))
			}
			return OptimizeExpression(ex)
		}
	}
}

func MergeExpression(a Expression, af float64, b Expression, bf float64) Expression {

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
	} else if rand.Float64() > r {
		tmp := a
		a = b
		b = tmp
	}

	newExpressions := []Expression{}

	for i := 0; i < int(math.Max(float64(len(a.Expressions)), float64(len(b.Expressions)))); i++ {
		if a.Operator == NOT && i > 0 {
			continue
		}
		if i < len(a.Expressions) && i < len(b.Expressions) {
			newExpressions = append(newExpressions, MergeExpression(a.Expressions[i], 0, b.Expressions[i], 0))
		} else {
			if i < len(a.Expressions) {
				if len(newExpressions) == 0 || rand.Float64() <= r {
					newExpressions = append(newExpressions, a.Expressions[i])
				}
			} else {
				if len(newExpressions) == 0 || rand.Float64() > r {
					newExpressions = append(newExpressions, b.Expressions[i])
				}
			}
		}
	}

	return OptimizeExpression(Expression{
		Operator:    a.Operator,
		Expressions: newExpressions,
	})
}

func OptimizeExpression(e Expression) Expression {

	if e.Operator == NOT {
		if e.Expressions[0].Operator == NOT {
			e = e.Expressions[0].Expressions[0]
		} else {
			e.Expressions[0] = OptimizeExpression(e.Expressions[0])
		}
	} else if e.Operator == AND || e.Operator == OR {
		onlyEq := true
		for i, ee := range e.Expressions {
			onlyEq = onlyEq && ee.Operator == EQ
			e.Expressions[i] = OptimizeExpression(ee)
		}
		if onlyEq || len(e.Expressions) == 1 {
			e = e.Expressions[0]
		}
	}

	return e

}

func MutateExpression(e Expression, inputLength int) Expression {
	if rand.Float64() < 0.01 {
		return CreateExpression(1, inputLength)
	}

	change := rand.Float64() < 0.2
	switch e.Operator {
	case EQ:
		{
			if change {
				e.Input = rand.Intn(inputLength)
			}
		}
	case NOT:
		{
			if change {
				e = e.Expressions[0]
			}
		}
	case AND:
		{
			if change {
				e.Operator = OR
			}
			if len(e.Expressions) > 1 && rand.Float64() < 0.2 {
				i := rand.Intn(len(e.Expressions))
				e.Expressions = append(e.Expressions[:i], e.Expressions[i+1:]...)
			}

			for i, ee := range e.Expressions {
				e.Expressions[i] = MutateExpression(ee, inputLength)
			}
		}
	case OR:
		{
			if change {
				e.Operator = AND
			}
			if len(e.Expressions) > 1 && rand.Float64() < 0.2 {
				i := rand.Intn(len(e.Expressions))
				e.Expressions = append(e.Expressions[:i], e.Expressions[i+1:]...)
			}

			for i, ee := range e.Expressions {
				e.Expressions[i] = MutateExpression(ee, inputLength)
			}
		}
	}
	return OptimizeExpression(e)
}
