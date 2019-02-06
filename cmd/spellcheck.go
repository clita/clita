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
	"fmt"
	"bufio"
	"os"
	"log"
	"regexp"
	"github.com/clita/spellcheck"
	"github.com/clita/diff"
	"github.com/spf13/cobra"
)

var segments bool
var color bool
var suggest bool
var correctedResult string

// spellcheckCmd represents the spellcheck command
var spellcheckCmd = &cobra.Command{
	Use:   "spellcheck",
	Short: "Spelling checking module for clita.",
	Run: func(cmd *cobra.Command, args []string) {
		spellcheck.Init()

		if suggest == false {
			if len(args) < 1 {
				log.Fatal("String not provided")
			}

			if segments == true {
				correctedResult = spellcheck.WordSegments(args[0])
			} else {
				correctedResult = spellcheck.Correctsentence(args[0])
			}
	
			if color == true {
				fmt.Println(diff.FindColouredChanges(args[0], correctedResult, "words"))
			} else {
				fmt.Println(correctedResult)
			}

		} else {
			// Improving error model
			fmt.Println("Enter each word suggestion in new line as mistake:correct pair.")
			fmt.Println("Example: mistaek:mistake")
			
			scanner := bufio.NewScanner(os.Stdin)
			re := regexp.MustCompile("^[a-zA-Z0-9']+:[a-zA-Z0-9']+$")
			correctRe := regexp.MustCompile("[^:]+$")
			mistakeRe := regexp.MustCompile("^[^:]+")

			for scanner.Scan() {
				inp := scanner.Text()
				
				if inp == "$$$" {
					break
				}

				match := re.FindString(inp)

				if inp != match {
					log.Fatal("Wrong input!")
				}

				correctWord := correctRe.FindString(inp)
				mistake	    := mistakeRe.FindString(inp)

				if _, present := spellcheck.ErrorModel[correctWord]; present == false {
					spellcheck.ErrorModel[correctWord] = make(map[string]int)
				}

				spellcheck.ErrorModel[correctWord][mistake]++
			}

			spellcheck.SaveMaps()
			fmt.Println("Saved succesfully!")

			if scanner.Err() != nil {
				// handle error.
				log.Fatal(scanner.Err())
			}
		}
	},
}

func init() {
	spellcheckCmd.PersistentFlags().BoolVar(&segments, "segments", false, "Segmentation of text into meaning full words.")
	spellcheckCmd.PersistentFlags().BoolVar(&color, "color", false, "To produce coloured diff of original string and it's correction. Default is false")
	spellcheckCmd.PersistentFlags().BoolVar(&suggest, "suggest", false, "Suggest mistake:correct pairs of words to improve it's working.")

	rootCmd.AddCommand(spellcheckCmd)
}
