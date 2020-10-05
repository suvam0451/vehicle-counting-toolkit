/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
	traffic "gitlab.com/suvam0451/trafficdetection/traffic"
)

// trailaltCmd represents the trailalt command
var trailaltCmd = &cobra.Command{
	Use:   "trailalt",
	Short: "Alternate algorithm for generating trails.",
	Long: `


The following default configuration is applied. Use a config file to override.
------------------------------------
Positive Reinforcement 		: 2 points
Negative Reinforcement		: 1 points (negative)
X threshold(default)		: 0.00075
Y threshold(default)		: 0.00075
	`,
	Run: func(cmd *cobra.Command, args []string) {
		traffic.DetectTrailCustom("inputnew", traffic.ModelParameters{
			Upvote:             2,
			Downvote:           -1,
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
