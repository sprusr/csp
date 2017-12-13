package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

func main() {
	instances := ReadSCInstancesFromFile("cutting_instances_with_solutions.txt")
	instance := instances[0]
	time := float64(800)

	// RandomSearch(instance, time)
	// CLONALG(instance, time, 20, 10, 5, 100)
	// ACO(instance, time, 0.5, 10, 1, 5)

	avg := 0
	for i := 0; i < 40; i++ {
		avg = ((avg * i) + CLONALG(instance, time, 20, 10, 5, 100).Cost) / (i + 1)
	}
	fmt.Println(avg)
}

// SCInstance represents a stock cutting problem instance
type SCInstance struct {
	StockLengths    []int
	StockCosts      []int
	OrderLengths    []int
	OrderQuantities []int
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

// SCSolution represents a stock cutting problem solution
type SCSolution struct {
	Instance *SCInstance
	Lengths  [][]int
	Cost     int
	Age      int
	Str      string
}

// GetCost returns the cost of a solution
func (solution *SCSolution) GetCost() int {
	cost := 0

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

func (solution *SCSolution) String() string {
	var buffer bytes.Buffer

	for i := range solution.Lengths {
		sort.Ints(solution.Lengths[i])
		for j := range solution.Lengths[i] {
			buffer.WriteString(strconv.Itoa(solution.Lengths[i][j]))
			buffer.WriteString(",")
		}
		buffer.WriteString(";")
	}

	solution.Str = buffer.String()

	return buffer.String()
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

	fmt.Printf("\n~~~\nBest route: %v\nCost: %v\n~~~\n", bestSolution.Lengths, bestSolution.Cost)

	return bestSolution
}
