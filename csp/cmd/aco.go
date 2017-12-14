package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sprusr/csp/csp"
)

func init() {
	rootCmd.AddCommand(acoCmd)
}

var acoCmd = &cobra.Command{
	Use:   "aco",
	Short: "Perform the ACO algorithm",
	Long:  "For the given instance, execute the Ant Colony Optimisation algorithm.",
	Run: func(cmd *cobra.Command, args []string) {
		instances := csp.GetSCInstances()
		time := float64(iterations)
		for i := 0; i < runs; i++ {
			solution := csp.ACO(instances[instance], time, 0.5, 10, 1, 5)
			fmt.Println(solution.Cost)
		}
		//fmt.Printf("\n~~~\nBest route: %v\nCost: %v\n~~~\n", solution.Lengths, solution.Cost)
	},
}
