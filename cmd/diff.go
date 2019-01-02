// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"github.com/clita/diff"
	"github.com/spf13/cobra"
	"os"
	"io/ioutil"
	"path/filepath"
	ui "github.com/gizak/termui"
)

var colorWords bool
var splitBy string
var strInp bool

// checking for errors
func checkUiError(e error) (error) {
    if e != nil {
        return fmt.Errorf("Error in opening the UI")
    }
    return nil
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	absPath, err1 := filepath.Abs(name)
    if _, err := os.Stat(absPath); err != nil || err1!=nil {
    	if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

// Function to display the coloured changes, splitBy = "words" for considering each word as one unit
// or splitBy = "lines" for considering each line as separate unit.
func display(firstString string, secondString string, splitBy string) {

	err := ui.Init()
	checkUiError(err)
	defer ui.Close()

	leftString, rightString := diff.FindColouredChanges(firstString, secondString, splitBy)
	leftStringWindow := ui.NewPar(leftString)
	leftStringWindow.Height = ui.TermHeight()
	leftStringWindow.Y = 0
	leftStringWindow.WrapLength=(6*ui.TermWidth())/14
	leftStringWindow.BorderLabel = "First Argument"
	leftStringWindow.BorderFg = ui.ColorYellow

	rightStringWindow := ui.NewPar(rightString)
	rightStringWindow.Height = ui.TermHeight()
	rightStringWindow.Y = 0
	rightStringWindow.BorderLabel = "Second Argument"
	rightStringWindow.BorderFg = ui.ColorYellow

	// build
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, leftStringWindow),
			ui.NewCol(6, 0, rightStringWindow)))

	// calculate layout
	ui.Body.Align()

	ui.Render(ui.Body)

	ui.Handle("q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("<Resize>", func(e ui.Event) {
		payload := e.Payload.(ui.Resize)
		ui.Body.Width = payload.Width
		leftStringWindow.Height = ui.TermHeight()
		rightStringWindow.Height = ui.TermHeight()
		leftStringWindow.WrapLength=(5*ui.TermWidth())/10
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Loop()
}

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Args: cobra.ExactArgs(2),
	Short: "Tool for comparing files passed as arguments",
	Long: `A tool to generate git-diff like coloured output showing comparison between files
	passed as arguments.`,
	Run: func(cmd *cobra.Command, args []string) {

		initCheck()
		
		if strInp == true {

			display(args[0], args[1], splitBy)

		} else{

			absPathArg0, err0 := filepath.Abs(args[0])
			absPathArg1, err1 := filepath.Abs(args[1])

			if (Exists(absPathArg0)) && (Exists(absPathArg1)) && err0==nil && err1==nil {

				firstString, _ := ioutil.ReadFile(absPathArg0)
				secondString,_ := ioutil.ReadFile(absPathArg1)

    			display(string(firstString), string(secondString), splitBy)
			} else {

				fmt.Printf("Invalid file name specified: \"%s\" and \"%s\"\n", absPathArg0, absPathArg1)

			}
		}

	},
}

func init() {
	diffCmd.PersistentFlags().BoolVarP(&strInp, "strings", "s", false, "Passing strings rather than files as arguments for comparing. Default is false")
	diffCmd.PersistentFlags().BoolVar(&colorWords, "color-words", false, "To consider each word as a different unit. Default is false")

	rootCmd.AddCommand(diffCmd)
}

func initCheck() {
	if colorWords == true {

		splitBy = "words"

	} else {

		splitBy = "lines"

	}

}
