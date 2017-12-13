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
		time := float64(100)
		solution := csp.ACO(instances[2], time, 0.5, 10, 1, 5)
		fmt.Printf("\n~~~\nBest route: %v\nCost: %v\n~~~\n", solution.Lengths, solution.Cost)
	},
}
