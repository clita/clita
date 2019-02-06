# diff
Diff module for clita application.  

### Install:
```sh
  go get github.com/clita/diff
```  

### Usage: 
```go

  // Return strings formatted according to colour codes
  diff.FindColouredChanges(firstString, secondString, splitBy, stdoutPrint)
  
  // firstString, secondString - Strings to compare
  // splitBy = "words" compared strings word by word, splitBy = "lines" to compare strings line by line
  // stdoutPrint boolean to tell which colour codes to use, 
  //            = true returns strings formatted according to terminal escape colour codes
  //            = false returns strings formatted according to gizak/termui colour codes
```
