# spellcheck
Spell checking module for clita.

## Usage 
**Initialisatio**  
```go
  spellcheck.Init()
```

**Correcting a text**  
```go
  spellcheck.Correctsentence("Speling Errurs IN somethink. Whutever; unusuel misteakes?")
  // Ex. Output: Spelling Errors In something. Whatever; unusual mistakes?
```  

**Words Segmentation**
```go
  spellcheck.WordSegments("thisisatestofsegmentationofaverylongsequenceofwords")
  // Ex. Output: this is a test of segmentation of a very long sequence of words
```

