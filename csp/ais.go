package csp

import (
	"math/rand"
	"sort"
	"time"
)

// CLONALG performs the CLONALG AIS algorithm with the given parameters
func CLONALG(instance SCInstance, searchTime float64, populationSize int, replacementThreshold int, cloneSizeFactor int, maxAge int) SCSolution {
	// seed randomness
	rand.Seed(int64(time.Now().Nanosecond()))

	bestSolution := instance.GenerateRandomSolution()

	// Initialisation
	population := generatePopulation(instance, populationSize)

	//startTime := time.Now()
	//endTime := startTime.UnixNano() + int64(searchTime*1000000000)

	for t := 0; t < int(searchTime); t++ {
		//for time.Now().UnixNano() < endTime {
		//for bestSolution.Cost > optimum {
		clonePool := []SCSolution{}

		// Cloning
		for i := 0; i < len(population); i++ {
			population[i].Age++

			if population[i].Age > maxAge {
				population = append(population[:i], population[i+1:]...)
			} else {
				for j := 0; j < cloneSizeFactor; j++ {
					clonePool = append(clonePool, population[i])
				}
			}
		}

		// Mutation
		for i := 0; i < len(clonePool); i++ {
			neighbours := clonePool[i].GetTwoOptNeighbourhood()
			clonePool[i] = neighbours[rand.Intn(len(neighbours))]
		}

		// Selection
		population = append(population, clonePool...)
		sort.Sort(byFitness(population))
		population = population[:populationSize]

		// Metadynamics
		population = append(population[:populationSize-replacementThreshold], generatePopulation(instance, replacementThreshold)...)
		sort.Sort(byFitness(population))

		if population[0].Cost < bestSolution.Cost {
			bestSolution = population[0]
		}

		// // Logging
		// for i := 0; i < len(population); i++ {
		// 	fmt.Println(population[i])
		// }
		// fmt.Println()
	}

	return bestSolution
}

func generatePopulation(instance SCInstance, size int) []SCSolution {
	routes := []SCSolution{}

	for i := 0; i < size; i++ {
		routes = append(routes, instance.GenerateRandomSolution())
	}

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

					newNeighbour.Cost = newNeighbour.GetCost()
					neighbours = append(neighbours, newNeighbour)
				}
			}
		}
	}

	return neighbours
}
