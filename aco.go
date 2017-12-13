package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// ACO performs Ant Colony Optimisation on the given stock cutting instance
func ACO(instance SCInstance, searchTime float64, evaporation float64, deposition float64, pheremoneImportance float64, costImportance float64) SCSolution {
	rand.Seed(int64(time.Now().Nanosecond()))

	bestSolution := instance.GenerateRandomSolution()

	// initialise pheremones
	initialPheremone := 0.2
	pheremones := make(map[[3]int]float64)
	for i := 0; i < len(instance.StockLengths); i++ {
		for j := 0; j < len(instance.OrderLengths); j++ {
			for k := j; k < len(instance.OrderLengths); k++ {
				pheremones[[3]int{i, instance.OrderLengths[j], instance.OrderLengths[k]}] = initialPheremone
				pheremones[[3]int{i, instance.OrderLengths[k], instance.OrderLengths[j]}] = initialPheremone
			}
			pheremones[[3]int{i, 0, instance.OrderLengths[j]}] = initialPheremone
		}
	}

	// startTime := time.Now()
	// endTime := startTime.UnixNano() + int64(searchTime*1000000000)

	for run := 0; run < int(searchTime); run++ {
		//for time.Now().UnixNano() < endTime {
		// initialize solution lengths 2D slice
		solution := SCSolution{Instance: &instance}
		solution.Lengths = make([][]int, len(instance.StockLengths))
		for i := 0; i < len(solution.Lengths); i++ {
			solution.Lengths[i] = make([]int, 0)
		}

		// construct a slice of ordered lengths and count
		orders := []int{}
		orderCount := 0
		for i := 0; i < len(instance.OrderLengths); i++ {
			orders = append(orders, instance.OrderQuantities[i])
			orderCount += instance.OrderQuantities[i]
		}

		// while we still have orders to be cut, keep cutting!
		for i := 0; i < orderCount; i++ {
			potentials := []potential{}

			// work out the pheremones for each possibility
			for j := range instance.StockLengths {
				for k := range instance.OrderLengths {
					if orders[k] > 0 {
						for l := range instance.StockLengths {
							from := 0
							if len(solution.Lengths[j]) != 0 {
								from = solution.Lengths[j][len(solution.Lengths[j])-1]
							}

							potentials = append(potentials, potential{
								Pheremone:  pheremones[[3]int{l, from, instance.OrderLengths[k]}],
								StockIndex: j,
								OrderIndex: k})
						}
					}
				}
			}

			// sort.Sort(byPotential(potentials))
			//
			// // pick the best option
			// orders[potentials[0].OrderIndex]--
			// solution.Lengths[potentials[0].StockIndex] = append(solution.Lengths[potentials[0].StockIndex], instance.OrderLengths[potentials[0].OrderIndex])

			// shuffle potentials
			shuffled := make([]potential, len(potentials))
			perm := rand.Perm(len(potentials))
			for i, v := range perm {
				shuffled[v] = potentials[i]
			}

			// set probabilistic stuff
			for i := range shuffled {
				shuffled[i].Probability = math.Pow(shuffled[i].Pheremone, pheremoneImportance) * math.Pow(1/float64(instance.OrderLengths[shuffled[i].OrderIndex]), costImportance)
			}
			for i := 1; i < len(shuffled); i++ {
				shuffled[i].Probability = shuffled[i].Probability + shuffled[i-1].Probability
				//fmt.Println(shuffled[i].Probability)
			}

			// choose what to do
			r := rand.Float64() * shuffled[len(shuffled)-1].Probability
			for i := range shuffled {
				if shuffled[i].Probability >= r {
					//fmt.Println(i, r)
					orders[shuffled[i].OrderIndex]--
					solution.Lengths[shuffled[i].StockIndex] = append(solution.Lengths[shuffled[i].StockIndex], instance.OrderLengths[shuffled[i].OrderIndex])
					break
				}
			}
		}

		solution.Cost = solution.GetCost()

		if solution.Cost < bestSolution.Cost {
			bestSolution = solution
			// fmt.Println(solution)
		}

		// update pheremones
		for i := range solution.Lengths {
			if len(solution.Lengths[i]) > 0 {
				pheremones[[3]int{i, 0, solution.Lengths[i][0]}] = pheremones[[3]int{i, 0, solution.Lengths[i][0]}] + (deposition / float64(solution.Cost))
			}
			for j := 0; j < len(solution.Lengths[i])-1; j++ {
				pheremones[[3]int{solution.Lengths[i][j], solution.Lengths[i][j+1]}] = pheremones[[3]int{solution.Lengths[i][j], solution.Lengths[i][j+1]}] + (deposition / float64(solution.Cost))
				for k := j + 1; k < len(solution.Lengths[i]); k++ {
					pheremones[[3]int{i, solution.Lengths[i][j], solution.Lengths[i][k]}] = pheremones[[3]int{i, solution.Lengths[i][j], solution.Lengths[i][k]}] + (deposition / float64(solution.Cost))
					pheremones[[3]int{i, solution.Lengths[i][k], solution.Lengths[i][j]}] = pheremones[[3]int{i, solution.Lengths[i][k], solution.Lengths[i][j]}] + (deposition / float64(solution.Cost))
				}
			}
		}
		// for i := range solution.Lengths {
		// 	if len(solution.Lengths[i]) > 0 {
		// 		pheremones[[3]int{i, 0, solution.Lengths[i][0]}] = pheremones[[3]int{i, 0, solution.Lengths[i][0]}] * (deposition / float64(solution.Cost))
		// 	}
		// 	for j := 0; j < len(solution.Lengths[i])-1; j++ {
		// 		pheremones[[3]int{i, solution.Lengths[i][j], solution.Lengths[i][j+1]}] = pheremones[[3]int{i, solution.Lengths[i][j], solution.Lengths[i][j+1]}] * (deposition / float64(solution.Cost))
		// 	}
		// }

		// decay pheremones
		for i := 0; i < len(instance.StockLengths); i++ {
			for j := 0; j < len(instance.OrderLengths); j++ {
				for k := j; k < len(instance.OrderLengths); k++ {
					pheremones[[3]int{i, instance.OrderLengths[j], instance.OrderLengths[k]}] = pheremones[[3]int{i, instance.OrderLengths[j], instance.OrderLengths[k]}] * evaporation
					pheremones[[3]int{i, instance.OrderLengths[k], instance.OrderLengths[j]}] = pheremones[[3]int{i, instance.OrderLengths[k], instance.OrderLengths[j]}] * evaporation
				}
				pheremones[[3]int{i, 0, instance.OrderLengths[j]}] = pheremones[[3]int{i, 0, instance.OrderLengths[j]}] * evaporation
			}
		}
	}

	fmt.Printf("\n~~~\nBest route: %v\nCost: %v\n~~~\n", bestSolution.Lengths, bestSolution.Cost)

	return bestSolution
}

type potential struct {
	Pheremone   float64
	StockIndex  int
	OrderIndex  int
	Probability float64
}
type byPotential []potential

func (s byPotential) Len() int {
	return len(s)
}
func (s byPotential) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byPotential) Less(i, j int) bool {
	return s[i].Pheremone < s[j].Pheremone
}
