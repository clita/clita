package autocomplete

import (
	"fmt"
	"github.com/alediaferia/prefixmap"
	"math"
	"os"
	"strings"
	"sort"
	"regexp"
	"log"
	"strconv"
	"bufio"
)

var WordMap map[string]float64

var similarityThreshold float64
var datasource *prefixmap.PrefixMap
var maxResults int

type Match struct {
	Value      string
	Similarity float64
	Frequency  float64
}

type MatchList []Match

func minimum(value0 int, values ...int) int {
	min := value0
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

func (match Match) Print() {
	fmt.Printf("match: %s\tfrequency: %.0f\tsimilarity: %.2f\t\n", match.Value, match.Frequency, match.Similarity)
}

// Print matches (maximum = maxResults)
func PrintMatches(matches []Match) {
	length := int(math.Min(float64(len(matches)), float64(maxResults)))
	for i := 0; i<length; i++ {
		matches[i].Print()
	}
}

// Computes similarity betwoon two strings returning result between 0 to 1
func ComputeSimilarity(str1, str2 string) float64 {
	ld := LevenshteinDistance(str1, str2)
	maxLen := math.Max(float64(len(str1)), float64(len(str2)))

	return 1.0 - (float64(ld)/float64(maxLen))
}

// Function to compute Levenshtein Distance
func LevenshteinDistance(source, destination string) int {
	vec1 := make([]int, len(destination)+1)
	vec2 := make([]int, len(destination)+1)

	w1 := []rune(source)
	w2 := []rune(destination)

	// initializing vec1
	for i := 0; i < len(vec1); i++ {
		vec1[i] = i
	}

	// initializing the matrix
	for i := 0; i < len(w1); i++ {
		vec2[0] = i + 1

		for j := 0; j < len(w2); j++ {
			cost := 1
			if w1[i] == w2[j] {
				cost = 0
			}
			min := minimum(vec2[j]+1,
				vec1[j+1]+1,
				vec1[j]+cost)
			vec2[j+1] = min
		}

		for j := 0; j < len(vec1); j++ {
			vec1[j] = vec2[j]
		}
	}

	return vec2[len(w2)]
}

// Implementing sort.Interface - Len, Swap, Less by Match
func (m MatchList) Len() int {
    return len(m)
}
func (m MatchList) Swap(i, j int) {
    m[i], m[j] = m[j], m[i]
}
func (m MatchList) Less(i, j int) bool {
    return (m[i].Frequency)*(m[i].Similarity) > (m[j].Frequency)*(m[j].Similarity)
}

func Init(threshold float64, maxres int, training_file string) {
	similarityThreshold = threshold
	maxResults = maxres
	datasource = prefixmap.New()

	// Entering values into WordMap
	WordMap = make(map[string]float64)

	file, err := os.Open(training_file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wordPattern := regexp.MustCompile("^[0-9a-zA-Z']+(\\s)?[a-zA-Z']+")
	numberPattern := regexp.MustCompile("\\d+$")
	for scanner.Scan() {
		w := wordPattern.FindString(scanner.Text())
		n := numberPattern.FindString(scanner.Text())
		WordMap[strings.ToLower(w)], err = strconv.ParseFloat(n, 64)

		if err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// here we populate the datasource
	for word, _ := range WordMap {
		parts := strings.Split(strings.ToLower(word), " ")
		for _, part := range parts {
			datasource.Insert(part, word)
		}
	}
}

// Provide autocomplete suggestions on input string
func Autocomplete(input string, printbool bool) []Match {
	values := datasource.GetByPrefix(strings.ToLower(input))
	results := make([]Match, 0, len(values))
	for _, v := range values {
		value := v.(string)
		s := ComputeSimilarity(value, input)
		if s >= similarityThreshold {
			m := Match{value, s, WordMap[value]}
			results = append(results, m)
		}
	}

	sort.Sort(MatchList(results))

	if printbool {
		fmt.Printf("Results (maximum %d) for target similarity: %.2f\n", maxResults, similarityThreshold)
		PrintMatches(results)
	}

	return results
}

// func main() {
// 	input := "holl"
// 	if input == "" {
// 		fmt.Println("Please, specify an input string")
// 		os.Exit(1)
// 	}
	
// 	Init(0.1, 10, "autocomplete.txt")
// 	Autocomplete(input, true)
// }
