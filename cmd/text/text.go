package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/alexMarco7/ga/pkg/ga"
)

type TextRules struct {
	Target []byte
}

func (tr TextRules) Create() interface{} {
	ba := make([]byte, len(tr.Target))
	for i := 0; i < len(tr.Target); i++ {
		ba[i] = byte(rand.Intn(95) + 32)
	}
	return ba
}
func (tr TextRules) Fitness(dna interface{}) float64 {
	score := 0
	value := (dna).([]byte)
	for i := 0; i < len(value); i++ {
		if value[i] == tr.Target[i] {
			score++
		}
	}
	return float64(score)
}
func (tr TextRules) Crossover(dna1 interface{}, dna2 interface{}) interface{} {

	a := dna1.([]byte)
	b := dna2.([]byte)
	child := make([]byte, len(a))

	mid := rand.Intn(len(a))
	for i := 0; i < len(a); i++ {
		if i > mid {
			child[i] = a[i]
		} else {
			child[i] = b[i]
		}
	}
	return child
}
func (tr TextRules) Mutate(dna interface{}) interface{} {
	value := (dna).([]byte)
	for {
		idx := rand.Intn(len(value))
		if value[idx] != tr.Target[idx] {
			value[idx] = byte(rand.Intn(95) + 32)
			break
		}
	}
	return value
}

func (tr TextRules) HasFinished(generation int, dna interface{}, fitness float64) bool {
	fmt.Printf("\r generation: %d | %s | fitness: %2f", generation, dna, fitness)
	return bytes.Compare(dna.([]byte), tr.Target) == 0
}

func main() {
	start := time.Now()

	fmt.Printf("%d", runtime.NumCPU())

	ga.Run(TextRules{
		Target: []byte(`Ouviram do Ipiranga as margens placidas.De um povo heroico o brado retumbante.E o sol da liberdade, em raios fulgidos.Brilhou no ceu da patria nesse instante.`),
	}, ga.Options{
		PopulationSize: 1000,
		MutationRate:   0.1,
	})
	elapsed := time.Since(start)
	fmt.Printf("\nTime taken: %s\n", elapsed)
}
