package diff

import (
	"errors"
	"regexp"
	"strings"
)

const INF = 10000000000
var memo[][] string

// Function that returns coloured strings by finding deletions im firstString and additions in secondString
// If splitBy is "words", one word is considered as one unit
// Else if splitBy is "lines", one line is considered as one unit
func FindColouredChanges(firstString string, secondString string, splitBy string) (string, string) {

	var colouredStringLeft, colouredStringRight string

	if splitBy == "words" {

		splitRegex := regexp.MustCompile("[\n \t]+")

		input1 := splitRegex.Split(firstString, INF)
		input2 := splitRegex.Split(secondString, INF)

		memo = make([][]string, len(input1))

		for i := 0; i < len(input1); i++ {
			memo[i] = make([]string, len(input2))
		}

		lcsWordString := lcsByWords(input1[0:], input2[0:], 0, 0)
		lcsWordString = trim(lcsWordString, ' ')
		// fmt.Println(lcsWordString)

		colouredStringLeft, _ = findLCSDeletions(firstString, lcsWordString, secondString, "words")
		colouredStringRight, _ = findLCSAdditions(firstString, lcsWordString, secondString, "words")

	} else if splitBy == "lines" {

		input1 := strings.Split(firstString, "\n")
		input2 := strings.Split(secondString, "\n")

		memo = make([][]string, len(input1))

		for i := 0; i < len(input1); i++ {
			memo[i] = make([]string, len(input2))
		}

		lcsLineString := lcsByLines(input1[0:], input2[0:], 0, 0)
		lcsLineString = trim(lcsLineString, '\n')
		// fmt.Println(lcsLineString)

		colouredStringLeft, _ = findLCSDeletions(firstString, lcsLineString, secondString, "lines")
		colouredStringRight, _ = findLCSAdditions(firstString, lcsLineString, secondString, "lines")

		colouredStringLeft = trim(colouredStringLeft, '\n')
		colouredStringRight = trim(colouredStringRight, '\n')
	}

	return colouredStringLeft, colouredStringRight
}

// Function that returns lcs string formed by concatenating words in lcs word array with spaces
// when one word is considered as one unit
func lcsByWords(s1 []string, s2 []string, indexS1 int, indexS2 int) string {
	lenghtS1, lenghtS2 := len(s1), len(s2)

	if indexS1 == lenghtS1 || indexS2 == lenghtS2 {
		return ""
	}

	if len(memo[indexS1][indexS2]) != 0 {
		return memo[indexS1][indexS2]
	}

	if s1[indexS1] == s2[indexS2] {
		if indexS1 == (lenghtS1-1) || indexS2 == (lenghtS2-1) {
			memo[indexS1][indexS2] = string(s1[indexS1])
			return memo[indexS1][indexS2]
		} else {
			memo[indexS1][indexS2] = string(s1[indexS1]) + " " + lcsByWords(s1, s2, indexS1+1, indexS2+1)
			return memo[indexS1][indexS2]
		}
	}

	result1 := lcsByWords(s1[0:], s2[0:], indexS1+1, indexS2)
	result2 := lcsByWords(s1[0:], s2[0:], indexS1, indexS2+1)

	if len(result1) >= len(result2) {
		memo[indexS1][indexS2] = result1
		return memo[indexS1][indexS2]
	} else {
		memo[indexS1][indexS2] = result2
		return memo[indexS1][indexS2]
	}
}

// Function to calculate lcs string when each line of string is considered as separate unit
func lcsByLines(s1 []string, s2 []string, indexS1 int, indexS2 int) string {

	lenghtS1, lenghtS2 := len(s1), len(s2)

	if indexS1 == lenghtS1 || indexS2 == lenghtS2 {
		return ""
	}

	if len(memo[indexS1][indexS2]) != 0 {
		return memo[indexS1][indexS2]
	}

	if s1[indexS1] == s2[indexS2] {

		if indexS1 == (lenghtS1-1) || indexS2 == (lenghtS2-1) {
			memo[indexS1][indexS2] = string(s1[indexS1])
		} else {
			memo[indexS1][indexS2] = string(s1[indexS1]) + "\n" + lcsByLines(s1, s2, indexS1+1, indexS2+1)
		}

		return memo[indexS1][indexS2]
	}

	result1 := lcsByLines(s1[0:], s2[0:], indexS1+1, indexS2)
	result2 := lcsByLines(s1[0:], s2[0:], indexS1, indexS2+1)

	if len(result1) >= len(result2) {
		memo[indexS1][indexS2] = result1
	} else {
		memo[indexS1][indexS2] = result2
	}

	return memo[indexS1][indexS2]
}

// This function finds the words(if splitBy = "words") / lines(if splitBy = "lines") that are not present in
// originalString but are there in modifiedString using the lcsString
//
// if splitBy == "words", lcsString is the string formed by concatenating words in lcs by words array with spaces
// when one word is considered as one unit
//
// if splitBy == "lines", lcsString is string formed by concatenating lines by '\n' in lcs by lines array
// when each line of string is considered as separate unit
func findLCSAdditions(originalString string, lcsString string, modifiedString string, splitBy string) (string, error) {
	// Intial Check
	if (len(lcsString) > len(modifiedString)) || (len(lcsString) > len(originalString)) {
		return "", errors.New("Wrong inputs!")
	}

	var finalColouredString string

	splitRegex := regexp.MustCompile("[\n \t]+")
	delimRegex := regexp.MustCompile("[^\n \t]+")

	if splitBy == "words" {
		modifiedStringArray := splitRegex.Split(modifiedString, INF)
		delimArray := delimRegex.Split(modifiedString, INF)

		var lcsArray []string
		if len(lcsString) != 0 {
			lcsArray = strings.Split(lcsString, " ")
		}

		finalColouredString = delimArray[0]

		j := 0
		for i := 0; i < len(modifiedStringArray); i++ {
			if (j < len(lcsArray)) && (modifiedStringArray[i] == lcsArray[j]) {
				finalColouredString = finalColouredString + modifiedStringArray[i]

				if i+1 < len(delimArray) {
					finalColouredString = finalColouredString + delimArray[i+1]
				}

				j++
			} else {
				finalColouredString = finalColouredString + "\033[32m" + modifiedStringArray[i] + "\033[0m"
				
				if i+1 < len(delimArray) {
					finalColouredString = finalColouredString + delimArray[i+1]
				}
			}
		}

	} else if splitBy == "lines" {
		var lcsArray []string
		if len(lcsString) != 0 {
			lcsArray = strings.Split(lcsString, "\n")
		}

		modifiedStringArray := strings.Split(modifiedString, "\n")

		j := 0
		for i := 0; i < len(modifiedStringArray); i++ {
			if (j < len(lcsArray)) && (modifiedStringArray[i] == lcsArray[j]) {
				finalColouredString = finalColouredString + modifiedStringArray[i] + "\n"
				j++
			} else {
				finalColouredString = finalColouredString + "\033[32m+" + modifiedStringArray[i] + "\033[0m\n"
			}
		}

	}

	return finalColouredString, nil
}

// This function finds the words(if splitBy = "words") / lines(if splitBy = "lines") that are present in
// originalString but not in modifiedString using lcsString
//
// if splitBy == "words", lcsString is the string formed by concatenating words in lcs by words array with spaces
// when one word is considered as one unit
//
// if splitBy == "lines", lcsString is string formed by concatenating lines by '\n' in lcs by lines array
// when each line of string is considered as separate unit
func findLCSDeletions(originalString string, lcsString string, modifiedString string, splitBy string) (string, error) {
	// Intial Check
	if (len(lcsString) > len(modifiedString)) || (len(lcsString) > len(originalString)) {
		return "", errors.New("Wrong inputs!")
	}

	var finalColouredString string = ""

	if splitBy == "words" {

		splitRegex := regexp.MustCompile("[\n \t]+")
		delimRegex := regexp.MustCompile("[^\n \t]+")

		originalStringArray := splitRegex.Split(originalString, INF)
		delimArray := delimRegex.Split(originalString, INF)

		var lcsArray []string
		if len(lcsString) != 0 {
			lcsArray = strings.Split(lcsString, " ")
		}

		finalColouredString = delimArray[0]

		j := 0
		for i := 0; i < len(originalStringArray); i++ {
			if (j < len(lcsArray)) && (originalStringArray[i] == lcsArray[j]) {
				finalColouredString = finalColouredString + originalStringArray[i]
				
				if i+1 < len(delimArray) {
					finalColouredString = finalColouredString + delimArray[i+1]
				}

				j++
			} else {
				finalColouredString = finalColouredString + "\033[31m" + originalStringArray[i] + "\033[0m"

				if i+1 < len(delimArray) {
					finalColouredString = finalColouredString + delimArray[i+1]
				}
			}
		}

	} else if splitBy == "lines" {

		var lcsArray []string
		if len(lcsString) != 0 {
			lcsArray = strings.Split(lcsString, "\n")
		}

		originalStringArray := strings.Split(originalString, "\n")

		j := 0
		for i := 0; i < len(originalStringArray); i++ {
			if (j < len(lcsArray)) && (originalStringArray[i] == lcsArray[j]) {
				finalColouredString = finalColouredString + originalStringArray[i] + "\n"
				j++
			} else {
				finalColouredString = finalColouredString + "\033[31m-" + originalStringArray[i] + "\033[0m\n"
			}
		}
	}

	return finalColouredString, nil
}

// Remove last newline of string
func trim(str string, ch byte) string {
	sz := len(str)
	if sz > 0 && str[sz-1] == ch {
		str = str[:sz-1]
	}
	return str
}
