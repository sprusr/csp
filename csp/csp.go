package csp

import (
	"math/rand"
	"time"
)

// SCInstance represents a stock cutting problem instance
type SCInstance struct {
	StockLengths    []int
	StockCosts      []float64
	OrderLengths    []int
	OrderQuantities []int
}

// GetSCInstances returns a SCInstance from the given file
func GetSCInstances() []SCInstance {
	return []SCInstance{
		SCInstance{
			StockLengths:    []int{10, 13, 15},
			StockCosts:      []float64{100, 130, 150},
			OrderLengths:    []int{3, 4, 5, 6, 7, 8, 9, 10},
			OrderQuantities: []int{5, 2, 1, 2, 4, 2, 1, 3}},
		SCInstance{
			StockLengths:    []int{4300, 4250, 4150, 3950, 3800, 3700, 3550, 3500},
			StockCosts:      []float64{86, 85, 83, 79, 68, 66, 64, 63},
			OrderLengths:    []int{2350, 2250, 2200, 2100, 2050, 2000, 1950, 1900, 1850, 1700, 1650, 1350, 1300, 1250, 1200, 1150, 1100, 1050},
			OrderQuantities: []int{2, 4, 4, 15, 6, 11, 6, 15, 13, 5, 2, 9, 3, 6, 10, 4, 8, 3}},
		SCInstance{
			StockLengths:    []int{120, 115, 110, 105, 100},
			StockCosts:      []float64{12, 11.5, 11, 10.5, 10},
			OrderLengths:    []int{21, 22, 24, 25, 27, 29, 30, 31, 32, 33, 34, 35, 38, 39, 42, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 59, 60, 61, 63, 65, 66, 67},
			OrderQuantities: []int{13, 15, 7, 5, 9, 9, 3, 15, 18, 17, 4, 17, 20, 9, 4, 19, 4, 12, 15, 3, 20, 14, 15, 6, 4, 7, 5, 19, 19, 6, 3, 7, 20, 5, 10, 17}}}
}

// SCSolution represents a stock cutting problem solution
type SCSolution struct {
	Instance *SCInstance
	Lengths  [][]int
	Cost     float64
	Age      int
	Str      string
}

// GetCost returns the cost of a solution
func (solution *SCSolution) GetCost() float64 {
	cost := 0.0

	for i := range solution.Lengths {
		count := 0

		for j := range solution.Lengths[i] {
			count += solution.Lengths[i][j]

			if count > solution.Instance.StockLengths[i] {
				cost += solution.Instance.StockCosts[i]
				count = solution.Lengths[i][j]
			}
		}

		if count > 0 {
			cost += solution.Instance.StockCosts[i]
		}
	}

	return cost
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
func RandomSearch(instance SCInstance, searchTime float64) SCSolution {
	rand.Seed(int64(time.Now().Nanosecond()))

	var newSolution SCSolution
	bestSolution := instance.GenerateRandomSolution()

	// startTime := time.Now()
	// endTime := startTime.UnixNano() + int64(searchTime*1000000000)

	for i := 0; i < int(searchTime); i++ {
		//for time.Now().UnixNano() < endTime {
		//for bestSolution.Cost > 5000 {
		newSolution = instance.GenerateRandomSolution()

		if newSolution.Cost < bestSolution.Cost {
			bestSolution = newSolution
			// fmt.Println(newSolution)
		}
	}

	//fmt.Printf("\n~~~\nBest route: %v\nCost: %v\n~~~\n", bestSolution.Lengths, bestSolution.Cost)

	return bestSolution
}
