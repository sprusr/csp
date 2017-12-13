package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sprusr/csp/csp"
)

func init() {
	rootCmd.AddCommand(aisCmd)
}

var aisCmd = &cobra.Command{
	Use:   "ais",
	Short: "Perform the AIS algorithm",
	Long:  "For the given instance, execute the CLONALG Artificial Immune System algorithm.",
	Run: func(cmd *cobra.Command, args []string) {
		instances := csp.GetSCInstances()
		time := float64(100)
		solution := csp.CLONALG(instances[2], time, 20, 10, 5, 100)
		fmt.Printf("\n~~~\nBest route: %v\nCost: %v\n~~~\n", solution.Lengths, solution.Cost)
	},
}
