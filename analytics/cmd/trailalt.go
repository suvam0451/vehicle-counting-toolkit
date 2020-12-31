/*Package cmd entry point for alternate trail detection algorithm*/
package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/suvam0451/trafficdetection/traffic"
)

// trailaltCmd represents the trailalt command
var trailaltCmd = &cobra.Command{
	Use:   "trailalt",
	Short: `Alternate algorithm for generating trails. See "help" for details`,
	Long: `
  	Information
  	----------------
    	Generates frame-by-frame information for vehicle objects visible (indicated by unique tag ID).

  	Output details
  	-------------
    	\out_traildetection_alt\veh_A.json      -->		All results with "no pruning"
    	\out_traildetection_alt\veh_A_c.json    -->		Objects with less than 5 data points are pruned


  	The following default configuration is applied. Use a config file to override.
  	------------------------------------
    	Positive Reinforcement 		: 2 points
    	Negative Reinforcement		: 1 points (negative)
    	X threshold(default)		: 0.00075
    	Y threshold(default)		: 0.00075
	`,
	Run: func(cmd *cobra.Command, args []string) {
		traffic.DetectTrailCustom("inputnew", traffic.ModelParameters{
			Rewards:            2,
			Penalty:            -1,
			XThreshold:         0.00075,
			YThreshold:         0.00075,
			EliminateThreshold: -2,
		})
	},
}

func init() {
	rootCmd.AddCommand(trailaltCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trailaltCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trailaltCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
