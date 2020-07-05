package ga

import (
	"math"
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"time"
)

type Options struct {
	PopulationSize int
	SurvivalRate   float64
	MaxGeneration  int
	MutationRate   float64
}

type Rules interface {
	Create() interface{}
	Fitness(interface{}) float64
	Crossover(interface{}, interface{}) interface{}
	Mutate(interface{}) interface{}
	HasFinished(int, interface{}, float64) bool
}

type Organism struct {
	DNA        interface{}
	Fitness    float64
	FitnessAcc float64
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

	survivorsQty := int(math.Floor(options.SurvivalRate * float64(options.PopulationSize)))

	for !found {
		generation++
		fitnessAcc := 0.0
		fitnessAcc, population = g.sort(population)
		bestOrganism := population[len(population)-1]

		hasFinished := g.Rules.HasFinished(generation, bestOrganism.DNA, bestOrganism.Fitness)

		if hasFinished {
			found = true
			result = bestOrganism.DNA
		} else {
			survivors := population[len(population)-survivorsQty:]
			population = g.naturalSelection(population[survivorsQty:], fitnessAcc)
			population = append(population, survivors...)
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

func (g *GA) sort(population []Organism) (float64, []Organism) {
	acc := 0.0
	sort.Slice(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})
	for i, o := range population {
		acc += o.Fitness
		population[i].FitnessAcc = acc

	}
	return acc, population
}

func (g *GA) naturalSelection(population []Organism, fitnessAcc float64) []Organism {

	couples := [][2]Organism{}

	for i := 0; i < len(population); i++ {
		a := population[len(population)-1]
		b := population[len(population)-1]

		ra := fitnessAcc * rand.Float64()
		rb := fitnessAcc * rand.Float64()

		aFound := false
		bFound := false

		for j := len(population) - 1; j >= 0; j-- {
			o := population[j]

			if o.FitnessAcc-o.Fitness <= ra && o.FitnessAcc > ra {
				aFound = true
				a = o
			}
			if o.FitnessAcc-o.Fitness <= rb && o.FitnessAcc > rb {
				bFound = true
				b = o
			}

			if aFound && bFound {
				break
			}
		}

		couples = append(couples, [2]Organism{a, b})
	}

	var wg sync.WaitGroup

	next := make([]Organism, len(population))

	for i, c := range couples {
		wg.Add(1)

		go func(a Organism, b Organism, i int) {
			childDNA := g.Rules.Crossover(a.DNA, b.DNA)

			if rand.Float64() < g.Options.MutationRate {
				childDNA = g.Rules.Mutate(childDNA)
			}

			child := Organism{DNA: childDNA}
			child.Fitness = g.Rules.Fitness(child.DNA)

			next[i] = child

			wg.Done()
		}(c[0], c[1], i)

		if (i+1)%g.CPUNum == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
	return next
}
