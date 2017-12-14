package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sprusr/csp/csp"
)

func init() {
	rootCmd.AddCommand(randomCmd)
}

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Perform the ACO algorithm",
	Long:  "For the given instance, execute the Ant Colony Optimisation algorithm.",
	Run: func(cmd *cobra.Command, args []string) {
		instances := csp.GetSCInstances()
		time := float64(iterations)
		// solution := csp.RandomSearch(instances[instance], time)
		// fmt.Printf("\n~~~\nBest route: %v\nCost: %v\n~~~\n", solution.Lengths, solution.Cost)
		avg := 0.0
		for i := 0; i < runs; i++ {
			avg = ((avg * float64(i)) + csp.RandomSearch(instances[instance], time).Cost) / (float64(i) + 1)
		}
		fmt.Printf("Random best: %v\n", avg)
	},
}
