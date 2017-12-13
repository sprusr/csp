package csp

import (
	"math"
	"math/rand"
	"time"
)

// ACO performs Ant Colony Optimisation on the given stock cutting instance
func ACO(instance SCInstance, searchTime float64, evaporation float64, deposition float64, pheremoneImportance float64, costImportance float64) SCSolution {
	rand.Seed(int64(time.Now().Nanosecond()))

	bestSolution := instance.GenerateRandomSolution()

	// initialise pheremones
	initialPheremone := 1.0
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

			// set probabilistic stuff
			for i := range potentials {
				potentials[i].Probability = math.Pow(potentials[i].Pheremone, pheremoneImportance) * math.Pow(1/float64(instance.OrderLengths[potentials[i].OrderIndex]), costImportance)
			}
			for i := 1; i < len(potentials); i++ {
				potentials[i].Probability = potentials[i].Probability + potentials[i-1].Probability
				//fmt.Println(potentials[i].Probability)
			}

			// choose what to do
			r := rand.Float64() * potentials[len(potentials)-1].Probability
			for i := range potentials {
				if potentials[i].Probability >= r {
					//fmt.Println(i, r)
					orders[potentials[i].OrderIndex]--
					solution.Lengths[potentials[i].StockIndex] = append(solution.Lengths[potentials[i].StockIndex], instance.OrderLengths[potentials[i].OrderIndex])
					break
				}
			}
		}

		solution.Cost = solution.GetCost()

		if solution.Cost < bestSolution.Cost {
			bestSolution = solution
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
	return s[i].Pheremone > s[j].Pheremone
}
