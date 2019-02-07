# clita
CLI based Text Assistant for dealing with files easily and smartly.  

## Installation

Download pre-compiled binaries for Mac, Windows and Linux from the [releases](https://github.com/clita/clita/releases) page.  

### Mac OSX Installation of pre-compiled binaries
Please submit an issue if they don't work for you.  

1. Make sure you download and place the binary in a folder that's on your `$PATH`.  If you are unsure what this means, go to *step 2*. Otherwise, skip to *step 3*
2. Create a `bin` directory under your home directory.
```
$ mkdir ~/bin
$ cd ~/bin
```
3. Add the following line at the end of your `~/.bash_profile` file.  [Link with instructions](https://natelandau.com/my-mac-osx-bash_profile/) on how to find this file
```sh
export PATH=$PATH:$HOME/bin
```
4. Download the `clita` binary for OSX and rename it.  
```sh
$ wget https://github.com/clita/clita/releases/download/v0.1.0-alpha/clita-darwin-amd64  
$ mv clita-darwin-amd64 clita
```
5. Finally, make the binary an executable file and you are good to go!
```
$ chmod +x clita
```

## Usage: 
```
  Clita is Command line application which serves as an assistant
        while writing, reading or editing your files.
        It has many features which includes autocompletion, autocorrection, comparing files.

  Usage:
    clita [command]

  Available Commands:
    autocomplete A brief description of your command
    diff         Tool for comparing files or strings passed as arguments
    help         Help about any command
    spellcheck   Spelling checking module for clita.

  Flags:
    -h, --help   help for clita

  Use "clita [command] --help" for more information about a command.
```  

### Examples: 
**Diff**   
this example passes strings arguments using "-s" flag but you can also pass files as arguments.  
like: `clita diff /path/to/file1 /path/to/file2`
```diff
$ clita diff -s "Hello World" "Hello"
  
# Output:  
  First Argument:  
- Hello World

  Second Argument:  
+ Hello 
```  

**SpellCheck**  
```diff
$ clita spellcheck --color "Speling Errurs IN somethink. Whutever; unusuel misteakes? Hellothereworld"

# Output: 
- Speling Errurs IN somethink. Whutever; unusuel misteakes? Hellothereworld
+ Spelling Errors In something. Whatever; unusual mistakes? Hello there world
```  

**Words Segmentation**  
```diff
$ clita spellcheck --segments "thisisatestofsegmentationofaverylongsequenceofwords"

# Output: 
this is a test of segmentation of a very long sequence of words
```  

**Autocorrect**  
```
$ clita autocomplete "hell"

// Output: 
// Results (maximum 5) for target similarity: 0.30
// match: hell             frequency: 3177030      similarity: 1.00
// match: hello all        frequency: 525838       similarity: 0.44
// match: hello and        frequency: 339673       similarity: 0.44
// match: hello to         frequency: 256602       similarity: 0.50
// match: hello from       frequency: 277949       similarity: 0.40
```  

```
$ clita autocomplete --threshold 0.2 --maxresults 10 "hell"

// Output:
// Results (maximum 10) for target similarity: 0.20
// match: hell             frequency: 3177030      similarity: 1.00
// match: hello all        frequency: 525838       similarity: 0.44
// match: hello and        frequency: 339673       similarity: 0.44
// match: hello to         frequency: 256602       similarity: 0.50
// match: hello from       frequency: 277949       similarity: 0.40
// match: hello everyone   frequency: 384383       similarity: 0.29
// match: hell is          frequency: 178001       similarity: 0.57
// match: hello there      frequency: 233646       similarity: 0.36
// match: hell and         frequency: 137950       similarity: 0.50
// match: say hello        frequency: 154092       similarity: 0.44
```  

  
