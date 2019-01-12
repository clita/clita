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
	ui "github.com/gizak/termui"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var INF = 1000000000
var colorWords bool
var splitBy string
var strInp bool
var controlsString string = "[q ](fg-black,bg-white):Exit   [n ](fg-black,bg-white):Next Page  [u ](fg-black,bg-white):Prev Page "
var metaStr string = "(fg-green,fg-bold)"
var pageIndex int

// checking for errors
func checkUiError(e error) error {
	if e != nil {
		return fmt.Errorf("Error in opening the UI")
	}
	return nil
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	absPath, err1 := filepath.Abs(name)
	if _, err := os.Stat(absPath); err != nil || err1 != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// // Function to process the strings in width and also for scrolling
func wrapString(str string, width int) string {
	delimRegex := regexp.MustCompile("[^\n \t]+")
	wordsdRegex := regexp.MustCompile("\n?[ \t]+")

	delimsArr := delimRegex.Split(str, INF)
	wordsArr := wordsdRegex.Split(str, INF)

	currWidth := 0
	modStr := ""

	for i := 0; i < len(wordsArr) && i < len(delimsArr); i++ {

		if len(delimsArr[i]) != 0 && delimsArr[i][0] == '\n' {
			modStr += delimsArr[i]
			currWidth = 0
		} else {
			currWidth += runewidth.StringWidth(delimsArr[i])
			modStr += delimsArr[i]
		}

		if strings.Contains(wordsArr[i], "fg-bold") {
			currWidth += runewidth.StringWidth(wordsArr[i]) - runewidth.StringWidth(metaStr)
			if currWidth > width {
				modStr += "\n" + wordsArr[i]
				currWidth = runewidth.StringWidth(wordsArr[i]) - runewidth.StringWidth(metaStr)
			} else {
				modStr += wordsArr[i]
			}
		} else {
			currWidth += runewidth.StringWidth(wordsArr[i])
			if currWidth > width {
				modStr += "\n" + wordsArr[i]
				currWidth = runewidth.StringWidth(wordsArr[i])
			} else {
				modStr += wordsArr[i]
			}
		}

	}

	return modStr

}

// Function to find min of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Function to process string for paging
func process(str string, index int, width int, height int) string {
	modString := wrapString(str, width)
	modStringArr := strings.Split(modString, "\n")

	y := min(height, len(modStringArr))
	x := index

	if x >= y {
		return strings.Join(modStringArr[y:], "\n")
	} else {
		return strings.Join(modStringArr[x:y], "\n")
	}
}

// Function to display the coloured changes, splitBy = "words" for considering each word as one unit
// or splitBy = "lines" for considering each line as separate unit.
func display(firstString string, secondString string, splitBy string) {

	err := ui.Init()
	checkUiError(err)
	defer ui.Close()

	leftString, rightString := diff.FindColouredChanges(firstString, secondString, splitBy)
	pageIndex = 0

	modLeftString := process(leftString, pageIndex, ui.TermWidth()/2-2, ui.TermHeight()-5)
	modRightString := process(rightString, pageIndex, ui.TermWidth()/2-2, ui.TermHeight()-5)

	leftStringWindow := ui.NewParagraph(modLeftString)
	leftStringWindow.Height = ui.TermHeight() - 3
	leftStringWindow.Y = 0
	leftStringWindow.BorderLabel = "First Argument"
	leftStringWindow.BorderFg = ui.ColorYellow

	rightStringWindow := ui.NewParagraph(modRightString)
	rightStringWindow.Height = ui.TermHeight() - 3
	rightStringWindow.Y = 0
	rightStringWindow.BorderLabel = "Second Argument"
	rightStringWindow.BorderFg = ui.ColorYellow

	controlsWindow := ui.NewParagraph(controlsString)
	controlsWindow.Height = 3
	controlsWindow.Y = 0
	// controlsWindow.BorderFg = ui.ColorGreen

	// build
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, leftStringWindow),
			ui.NewCol(6, 0, rightStringWindow)),
		ui.NewRow(
			ui.NewCol(8, 2, controlsWindow)))

	// calculate layout
	ui.Body.Align()

	ui.Render(ui.Body)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			ui.Body.Width = payload.Width
			pageIndex = 0
			modLeftString := process(leftString, pageIndex, ui.TermWidth()-2, ui.TermHeight()-5)
			modRightString := process(rightString, pageIndex, ui.TermWidth()-2, ui.TermHeight()-5)
			leftStringWindow.Text = modLeftString
			rightStringWindow.Text = modRightString
			leftStringWindow.Height = ui.TermHeight() - 3
			rightStringWindow.Height = ui.TermHeight() - 3
			ui.Body.Align()
			ui.Clear()
			ui.Render(ui.Body)
		case "n":
			pageIndex += 1
			modLeftString := process(leftString, pageIndex, leftStringWindow.Width, ui.TermHeight()-5)
			modRightString := process(rightString, pageIndex, rightStringWindow.Width, ui.TermHeight()-5)
			leftStringWindow.Text = modLeftString
			rightStringWindow.Text = modRightString
			ui.Body.Align()
			ui.Clear()
			ui.Render(ui.Body)
		case "u":
			if pageIndex > 0 {
				pageIndex -= 1
				modLeftString := process(leftString, pageIndex, leftStringWindow.Width, ui.TermHeight()-5)
				modRightString := process(rightString, pageIndex, rightStringWindow.Width, ui.TermHeight()-5)
				leftStringWindow.Text = modLeftString
				rightStringWindow.Text = modRightString
				ui.Body.Align()
				ui.Clear()
				ui.Render(ui.Body)
			}
		}
	}
}

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Args:  cobra.ExactArgs(2),
	Short: "Tool for comparing files passed as arguments",
	Long: `A tool to generate git-diff like coloured output showing comparison between files
	passed as arguments.`,
	Run: func(cmd *cobra.Command, args []string) {

		initCheck()

		if strInp == true {

			display(args[0], args[1], splitBy)

		} else {

			absPathArg0, err0 := filepath.Abs(args[0])
			absPathArg1, err1 := filepath.Abs(args[1])

			if (Exists(absPathArg0)) && (Exists(absPathArg1)) && err0 == nil && err1 == nil {

				firstString, _ := ioutil.ReadFile(absPathArg0)
				secondString, _ := ioutil.ReadFile(absPathArg1)

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
