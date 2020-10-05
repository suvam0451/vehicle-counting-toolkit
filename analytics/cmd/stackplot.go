package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	traffic "gitlab.com/suvam0451/trafficdetection/traffic"
)

// stackplotCmd represents the stackplot command
var stackplotCmd = &cobra.Command{
	Use:   "stackplot",
	Short: "Generate stack plots for output data from trailalt",
	Long:  `For every JSON file in input_stackplot, generates stackplot data to ./out_stackplot.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generating stackplots to ./stackplot...")
		traffic.GenerateStackplot()
	},
}

func init() {
	rootCmd.AddCommand(stackplotCmd)
}
