// Copyright © 2019 VAIBHAV THAKKAR <vaibhav.thakkar.22.12.99@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"path/filepath"
	"log"
	"github.com/clita/autocomplete"
	"github.com/spf13/cobra"
)

var threshold float64
var maxresults int
var training_file string

// autocompleteCmd represents the autocomplete command
var autocompleteCmd = &cobra.Command{
	Use:   "autocomplete",
	Args:  cobra.ExactArgs(1),
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if threshold < 0 || threshold > 1 {
			log.Fatal("Invalid threshold value!")
		} else if maxresults < 0 {
			log.Fatal("Invalid max results value!")
		} else {
			autocomplete.Init(threshold, maxresults, training_file)
			_ = autocomplete.Autocomplete(args[0], true)
		}
	},
}

func init() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	training_file = dir + "/autocomplete.txt"

	autocompleteCmd.PersistentFlags().Float64VarP(&threshold, "threshold", "t", 0.3, "Minimum similarity between input and suggested strings between 0 and 1. Default - 0.3")
	autocompleteCmd.PersistentFlags().IntVarP(&maxresults, "maxresults", "m", 5, "Maximum number of suggested results. Default is 5")
	autocompleteCmd.PersistentFlags().StringVarP(&training_file, "file", "f", training_file, "Training file to use for suggestion.")

	rootCmd.AddCommand(autocompleteCmd)
}
