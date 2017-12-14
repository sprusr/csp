package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sprusr/csp/csp"
)

func init() {
	rootCmd.AddCommand(compareCmd)
}

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Perform a comparison between algorithms",
	Long:  "Averages the performance of the different algorithms after a defined number of iterations, over a defined number of runs.",
	Run: func(cmd *cobra.Command, args []string) {
		instances := csp.GetSCInstances()
		time := float64(iterations)

		avg := 0.0
		for i := 0; i < runs; i++ {
			avg = ((avg * float64(i)) + csp.ACO(instances[instance], time, 0.95, 100, 5, 1).Cost) / (float64(i) + 1)
		}
		fmt.Printf("ACO best: %v\n", avg)

		avg = 0
		for i := 0; i < runs; i++ {
			avg = ((avg * float64(i)) + csp.CLONALG(instances[instance], time, 20, 10, 5, 100).Cost) / (float64(i) + 1)
		}
		fmt.Printf("AIS best: %v\n", avg)

		avg = 0
		for i := 0; i < runs; i++ {
			avg = ((avg * float64(i)) + csp.RandomSearch(instances[instance], time).Cost) / (float64(i) + 1)
		}
		fmt.Printf("Random best: %v\n", avg)
	},
}
