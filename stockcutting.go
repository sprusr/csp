package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// SCInstance represents a stock cutting problem instance
type SCInstance struct {
	StockLengths    []int
	StockCosts      []int
	OrderLengths    []int
	OrderQuantities []int
}

// SCSolution represents a stock cutting problem solution
type SCSolution struct {
	Instance *SCInstance
	Lengths  [][][]int
	Cost     int
}

func main() {
	fmt.Printf("todo\n")
	instances := ReadSCInstancesFromFile("cutting_instances_with_solutions.txt")

	instances[0].RandomSearch(30000)
}

// RandomSearch performs a random search on the given stock cutting instance
func (instance SCInstance) RandomSearch(steps int) {
	var newSolution SCSolution
	bestSolution := instance.GenerateRandomSolution()

	for i := 0; i < steps; i++ {
		//for bestSolution.Cost > 5000 {
		newSolution = instance.GenerateRandomSolution()
		fmt.Print(newSolution)

		if newSolution.Cost < bestSolution.Cost {
			bestSolution = newSolution
			fmt.Print(" << NEW BEST")
		}

		fmt.Println()
	}

	fmt.Println()
	fmt.Println("~~~")
	fmt.Print("Best solution: ")
	fmt.Println(bestSolution)
	fmt.Println("~~~")
	fmt.Println()
}

// ReadSCInstancesFromFile returns a SCInstance from the given file
func ReadSCInstancesFromFile(file string) []SCInstance {
	tspFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer tspFile.Close()

	return []SCInstance{
		SCInstance{
			StockLengths:    []int{10, 13, 15},
			StockCosts:      []int{100, 130, 150},
			OrderLengths:    []int{3, 4, 5, 6, 7, 8, 9, 10},
			OrderQuantities: []int{5, 2, 1, 2, 4, 2, 1, 3}},
		SCInstance{
			StockLengths:    []int{4300, 4250, 4150, 3950, 3800, 3700, 3550, 3500},
			StockCosts:      []int{86, 85, 83, 79, 68, 66, 64, 63},
			OrderLengths:    []int{2350, 2250, 2200, 2100, 2050, 2000, 1950, 1900, 1850, 1700, 1650, 1350, 1300, 1250, 1200, 1150, 1100, 1050},
			OrderQuantities: []int{2, 4, 4, 15, 6, 11, 6, 15, 13, 5, 2, 9, 3, 6, 10, 4, 8, 3}}}
}

// GenerateRandomSolution returns a random stock cutting problem solution for the given instance
func (instance SCInstance) GenerateRandomSolution() SCSolution {
	// seed randomness
	rand.Seed(int64(time.Now().Nanosecond()))

	// initialize solution Lengths 3D array
	solution := SCSolution{Instance: &instance}
	solution.Lengths = make([][][]int, len(instance.StockLengths))
	for i := 0; i < len(solution.Lengths); i++ {
		solution.Lengths[i] = make([][]int, 1)
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

	// put the ordered lengths into random stocks that they will fit in
	for orderIndex := range shuffledOrders {
		// randomly select a stock length
		stockIndex := rand.Intn(len(instance.StockLengths))
		for instance.StockLengths[stockIndex] < shuffledOrders[orderIndex] {
			stockIndex = rand.Intn(len(instance.StockLengths))
		}

		endBinIndex := len(solution.Lengths[stockIndex]) - 1
		endBinSum := 0

		for i := range solution.Lengths[stockIndex][endBinIndex] {
			endBinSum += solution.Lengths[stockIndex][endBinIndex][i]
		}

		if endBinSum+shuffledOrders[orderIndex] > instance.StockLengths[stockIndex] {
			endBinIndex++
			solution.Lengths[stockIndex] = append(solution.Lengths[stockIndex], make([]int, 0))
		}

		solution.Lengths[stockIndex][endBinIndex] = append(solution.Lengths[stockIndex][endBinIndex], shuffledOrders[orderIndex])
	}

	// set the solution Cost
	solution.Cost = solution.GetCost()

	return solution
}

// GetCost returns the cost of a solution
func (solution *SCSolution) GetCost() int {
	cost := 0

	for i := 0; i < len(solution.Lengths); i++ {
		for j := 0; j < len(solution.Lengths[i]); j++ {
			// for k := 0; k < len(solution.Lengths[i][j]); k++ {
			// 	//
			// }
			cost += solution.Instance.StockCosts[i]
		}
	}

	return cost
}
