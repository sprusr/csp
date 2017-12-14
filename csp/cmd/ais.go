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
		time := float64(iterations)
		for i := 0; i < runs; i++ {
			solution := csp.CLONALG(instances[instance], time, 20, 10, 5, 100)
			fmt.Println(solution.Cost)
		}
	},
}
