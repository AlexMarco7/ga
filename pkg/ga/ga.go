package ga

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Options struct {
	PopulationSize int
	MaxGeneration  int
	MutationRate   float32
}

type Rules interface {
	Create() interface{}
	Fitness(interface{}) float64
	Crossover(interface{}, interface{}) interface{}
	Mutate(interface{}) interface{}
	HasFinished(int, interface{}, float64) bool
}

type Organism struct {
	DNA     interface{}
	Fitness float64
}

type GA struct {
	Rules   Rules
	Options Options
	CPUNum  int
}

func Run(rules Rules, options Options) interface{} {
	g := GA{rules, options, runtime.NumCPU()}

	rand.Seed(time.Now().UTC().UnixNano())

	population := g.createPopulation()

	found := false
	generation := 0
	var result interface{}

	for !found {
		generation++
		bestOrganism := g.getBest(population)

		hasFinished := g.Rules.HasFinished(generation, bestOrganism.DNA, bestOrganism.Fitness)

		if hasFinished {
			found = true
			result = bestOrganism.DNA
		} else {
			maxFitness := bestOrganism.Fitness
			pool := g.createPool(population, maxFitness)
			population = g.naturalSelection(pool, population)
		}

	}
	return result

}

func (g *GA) createPopulation() (population []Organism) {
	population = make([]Organism, g.Options.PopulationSize)

	var wg sync.WaitGroup
	for i := 1; i <= g.Options.PopulationSize; i++ {
		wg.Add(1)
		go func(i int) {
			population[i] = Organism{DNA: g.Rules.Create()}
			population[i].Fitness = g.Rules.Fitness(population[i].DNA)
			wg.Done()
		}(i - 1)
		if i%g.CPUNum == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
	return
}

// Get the best organism
func (g *GA) getBest(population []Organism) Organism {
	best := 0.0
	index := 0
	for i := 0; i < len(population); i++ {
		if population[i].Fitness > best {
			index = i
			best = population[i].Fitness
		}
	}
	return population[index]
}

func (g *GA) createPool(population []Organism, maxFitness float64) []Organism {
	pool := make([]Organism, 0)

	// create a pool for next generation
	var wg sync.WaitGroup
	for i := 1; i <= len(population); i++ {
		wg.Add(1)
		go func(i int) {
			num := int((population[i].Fitness / maxFitness) * 100)
			for n := 0; n < num; n++ {
				pool = append(pool, population[i])
			}
			wg.Done()
		}(i - 1)
		if i%g.CPUNum == 0 {
			wg.Wait()
		}
	}
	wg.Wait()

	return pool
}

// perform natural selection to create the next generation
func (g *GA) naturalSelection(pool []Organism, population []Organism) []Organism {
	next := make([]Organism, len(population))

	var wg sync.WaitGroup
	for i := 1; i <= len(population); i++ {
		wg.Add(1)
		go func(i int) {
			r1, r2 := rand.Intn(len(pool)), rand.Intn(len(pool))
			a := pool[r1]
			b := pool[r2]

			childDNA := g.Rules.Crossover(a.DNA, b.DNA)
			child := Organism{DNA: childDNA}

			if rand.Float32() < g.Options.MutationRate {
				child = Organism{DNA: g.Rules.Mutate(childDNA)}
			}

			child.Fitness = g.Rules.Fitness(child.DNA)

			next[i] = child

			wg.Done()
		}(i - 1)
		if i%g.CPUNum == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
	return next
}
