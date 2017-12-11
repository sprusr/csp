package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// CLONALG performs the CLONALG AIS algorithm with the given parameters
func CLONALG(instance SCInstance, populationSize int, replacementThreshold int, cloneSizeFactor int, runs int) {
	// seed randomness
	rand.Seed(int64(time.Now().Nanosecond()))

	// Initialisation
	population := generatePopulation(instance, populationSize)

	for t := 0; t < runs; t++ {
		clonePool := []SCSolution{}

		// Cloning
		for i := 0; i < populationSize; i++ {
			for j := 0; j < cloneSizeFactor; j++ {
				clonePool = append(clonePool, population[i])
			}
		}

		// Mutation
		for i := 0; i < len(clonePool); i++ {
			neighbours := clonePool[i].GetTwoOptNeighbourhood()
			rand.Seed(int64(time.Now().Nanosecond()))
			clonePool[i] = neighbours[rand.Intn(len(neighbours))]
		}

		// TODO Ageing

		// Selection
		population = append(population, clonePool...)
		sort.Sort(byFitness(population))
		population = population[:populationSize]

		// Metadynamics
		population = append(population[:populationSize-replacementThreshold], generatePopulation(instance, replacementThreshold)...)
		sort.Sort(byFitness(population))
		for i := 0; i < len(population); i++ {
			population[i].NormalisedFitness = float64(population[0].GetCost()) / float64(population[i].GetCost())
		}

		// Logging
		for i := 0; i < len(population); i++ {
			fmt.Println(population[i])
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("~~~")
	fmt.Print("Best route: ")
	fmt.Println(population[0])
	fmt.Println("~~~")
	fmt.Println()
}

func generatePopulation(instance SCInstance, size int) []SCSolution {
	routes := []SCSolution{}

	for i := 0; i < size; i++ {
		routes = append(routes, instance.GenerateRandomSolution())
	}

	sort.Sort(byFitness(routes))

	return routes
}

type byFitness []SCSolution

func (p byFitness) Len() int {
	return len(p)
}
func (p byFitness) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p byFitness) Less(i, j int) bool {
	return p[i].GetCost() < p[j].GetCost()
}

// GetTwoOptNeighbourhood returns an array of neighbours for the given stock cutting problem solution
func (solution *SCSolution) GetTwoOptNeighbourhood() []SCSolution {
	neighbours := []SCSolution{}
	neighboursMap := make(map[string]bool)

	for i := range solution.Lengths {
		for j := range solution.Lengths[i] {
			for k := 0; k < len(solution.Lengths); k++ {
				if i != k {
					newNeighbour := SCSolution{Instance: solution.Instance}
					newNeighbour.Lengths = make([][]int, len(solution.Lengths))
					for m := range solution.Lengths {
						newNeighbour.Lengths[m] = make([]int, len(solution.Lengths[m]))
						copy(newNeighbour.Lengths[m], solution.Lengths[m])
					}

					newNeighbour.Lengths[k] = append(newNeighbour.Lengths[k], newNeighbour.Lengths[i][j])
					newNeighbour.Lengths[i] = append(newNeighbour.Lengths[i][:j], newNeighbour.Lengths[i][j+1:]...)

					solutionString := newNeighbour.String()

					if !neighboursMap[solutionString] {
						neighboursMap[solutionString] = true
						newNeighbour.Cost = newNeighbour.GetCost()
						neighbours = append(neighbours, newNeighbour)
					}
				}
			}
		}
	}

	return neighbours
}

// GenerateRandomSolution returns a random stock cutting problem solution for the given instance
func (instance SCInstance) GenerateRandomSolution() SCSolution {
	// initialize solution lengths 2D slice
	solution := SCSolution{Instance: &instance}
	solution.Lengths = make([][]int, len(instance.StockLengths))
	for i := 0; i < len(solution.Lengths); i++ {
		solution.Lengths[i] = make([]int, 0)
	}

	// construct a slice of ordered lengths and shuffle it
	orders := []int{}
	for i := 0; i < len(instance.OrderLengths); i++ {
		for j := 0; j < instance.OrderQuantities[i]; j++ {
			orders = append(orders, instance.OrderLengths[i])
		}
	}
	shuffledOrders := make([]int, len(orders))
	perm := rand.Perm(len(orders))
	for i, v := range perm {
		shuffledOrders[v] = orders[i]
	}

	// put the ordered lengths into random stocks that fit them
	for orderIndex := range shuffledOrders {
		stockIndex := rand.Intn(len(instance.StockLengths))
		for instance.StockLengths[stockIndex] < shuffledOrders[orderIndex] {
			stockIndex = rand.Intn(len(instance.StockLengths))
		}
		solution.Lengths[stockIndex] = append(solution.Lengths[stockIndex], shuffledOrders[orderIndex])
	}

	// set the solution Cost
	solution.Cost = solution.GetCost()

	return solution
}

// RandomSearch performs a random search on the given stock cutting instance
func (instance SCInstance) RandomSearch(steps int) SCSolution {
	rand.Seed(int64(time.Now().Nanosecond()))

	var newSolution SCSolution
	bestSolution := instance.GenerateRandomSolution()

	for i := 0; i < steps; i++ {
		//for bestSolution.Cost > 5000 {
		newSolution = instance.GenerateRandomSolution()

		if newSolution.Cost < bestSolution.Cost {
			bestSolution = newSolution
			fmt.Println(newSolution)
		}
	}

	fmt.Println()
	fmt.Println("~~~")
	fmt.Print("Best solution: ")
	fmt.Println(bestSolution)
	fmt.Println("~~~")
	fmt.Println()

	return bestSolution
}
