package main

import (
	"bytes"
	"os"
)

func main() {
	instances := ReadSCInstancesFromFile("cutting_instances_with_solutions.txt")

	instances[1].RandomSearch(100000)

	//CLONALG(instances[0], 20, 10, 5, 2000)

	ACO(instances[1], 100000, 0.5, 100, 1, 5)
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
	Instance          *SCInstance
	Lengths           [][]int
	Cost              int
	NormalisedFitness float64
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
		for j := range solution.Lengths[i] {
			buffer.WriteString(string(solution.Lengths[i][j]))
			buffer.WriteString(",")
		}
		buffer.WriteString(";")
	}

	return buffer.String()
}
