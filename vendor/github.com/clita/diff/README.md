# diff
Diff module for clita application. Implemented using Memoized Longest Common Subsequence algorithm.

### Install:
```sh
  go get github.com/clita/diff
```  

### Usage: 
```go

  // Return strings formatted according to colour codes
  diff.FindColouredChanges(firstString, secondString, splitBy)
  
  // firstString, secondString - Strings to compare
  // splitBy = "words" compared strings word by word, splitBy = "lines" to compare strings line by line
```
