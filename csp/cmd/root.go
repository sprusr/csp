package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().IntVarP(&iterations, "iterations", "i", 40, "number of algorithm iterations")
	rootCmd.PersistentFlags().IntVarP(&instance, "instance", "I", 0, "instance on which to operate")
	rootCmd.PersistentFlags().IntVarP(&runs, "runs", "r", 40, "number of algorithm runs")
}

var iterations int
var instance int
var runs int

var rootCmd = &cobra.Command{
	Use:   "csp",
	Short: "Solves the CSP using Computational Intelligence",
	Long:  "Uses various Computational Intelligence algorithms to produce solutions for the Cutting Stock Problem.",
}

// Execute executes the root command
func Execute() {
	rootCmd.Execute()
}
