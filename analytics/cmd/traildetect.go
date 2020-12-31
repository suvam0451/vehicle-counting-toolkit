/*Package cmd is entry point for the original algorithm */
package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	traffic "gitlab.com/suvam0451/trafficdetection/traffic"
	"gitlab.com/suvam0451/trafficdetection/utility"
)

// traildetectCmd represents the traildetect command
var traildetectCmd = &cobra.Command{
	Use:   "traildetect",
	Short: "Previous algorithm for generating trails. Use help for detailed information.",
	Long: `
	Information
	----------------
	  Generates "frame-by-frame" information for vehicle objects visible (indicated by unique tag ID).
      This variant is multi-threaded

	Output details
	-------------
	  \out_traildetection_alt\veh_A.json      -->		"frame-by-frame" info + no pruning
	  \out_traildetection_alt\veh_A_c.json    -->		"object-by-object" informations.
  
  
	The following default configuration is applied. Use a config file to override.
	-----------------------------------------------------------
	  Positive Reinforcement         : 2 points
	  Negative Reinforcement         : 1 points (negative)
	  X threshold(default)           : 0.00025
	  Y threshold(default)           : 0.00025
	  Elimination threshold          : 0
	  `,
	Run: func(cmd *cobra.Command, args []string) {
		if configBytes, err := utility.ReadJSON("./config.json"); err == nil {
			tmp := traffic.ConfigFileSchema{}
			json.Unmarshal(configBytes, &tmp)
			traffic.DetectTrail(tmp.InputFiles.TrailDetectAlt, tmp.TrailDetectAlt)
		} else {
			panic("config.json file could not be found. Make sure you have appropriate permissions set")
		}
	},
}

func init() {
	rootCmd.AddCommand(traildetectCmd)
}
