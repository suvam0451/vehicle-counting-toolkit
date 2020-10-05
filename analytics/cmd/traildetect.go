/*Package cmd is entry point for the original algorithm */
package cmd

import (
	"github.com/spf13/cobra"
	traffic "gitlab.com/suvam0451/trafficdetection/traffic"
)

// traildetectCmd represents the traildetect command
var traildetectCmd = &cobra.Command{
	Use:   "traildetect",
	Short: "Previous algorithm for generating trails. Use help for detailed information.",
	Long: `
	Information
	----------------
	  Generates frame-by-frame information for vehicle objects visible (indicated by unique tag ID).
      This variant is multi-threaded

	Output details
	-------------
	  \out_traildetection_alt\veh_A.json      -->		All results with "no pruning"
	  \out_traildetection_alt\veh_A_c.json    -->		Objects with less than 5 data points are pruned
  
  
	The following default configuration is applied. Use a config file to override.
	-----------------------------------------------------------
	  Positive Reinforcement         : 2 points
	  Negative Reinforcement         : 1 points (negative)
	  X threshold(default)           : 0.00025
	  Y threshold(default)           : 0.00025
	  Elimination threshold          : 0
	  `,
	Run: func(cmd *cobra.Command, args []string) {
		traffic.DetectTrail("./input_traildetect", traffic.ModelParameters{
			Upvote:             2,
			Downvote:           -1,
			XThreshold:         0.00025,
			YThreshold:         0.00025,
			EliminateThreshold: 0,
		})
	},
}

func init() {
	rootCmd.AddCommand(traildetectCmd)
}
