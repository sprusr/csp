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
	Short: "Perform a random search",
	Long:  "For the given instance, execute a random search.",
	Run: func(cmd *cobra.Command, args []string) {
		instances := csp.GetSCInstances()
		time := float64(iterations)
		for i := 0; i < runs; i++ {
			solution := csp.RandomSearch(instances[instance], time)
			fmt.Println(solution.Cost)
		}
	},
}
