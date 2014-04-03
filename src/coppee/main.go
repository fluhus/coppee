package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"flag"
)

// Returns a copier walk-function.
// basedir:          base walking directory, must match the directory walked
// rules:            copying rules
// collapseOnError:  if true, process will stop upon copying error
// overwrite:        if true, will overwrite existing files
// verbose:          if true, will print copies made and errors
// pretend:          if true, will only pretend to copy files
func copier(
		basedir string,
		rules []copyRule,
		collapseOnError bool,
		overwrite bool,
		verbose bool,
		pretend bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// Relative path - matching will be to relative path
		relPath, relPathErr := filepath.Rel(basedir, path)
		if relPathErr != nil { panic(relPathErr.Error()) } // should not happen
		
		// Don't try to copy directories
		if info.IsDir() {
			return nil
		}

		for _,rule := range rules {
			// Check if filename matches
			// 'path' is source file name
			if globalMatch(rule.src, relPath) {
				// Destination file
				dst := rule.src.ReplaceAllString(path, rule.dst)

				// Overwrite?
				if !overwrite && fexists(dst) {
					if verbose {
						fmt.Println("skipping: '" + path + "' to '" + dst + "'")
					}
					continue
				}
				
				// Print
				if verbose {
					fmt.Println("copying:  '" + path + "' to '" + dst + "'")
				}
				
				// Copy
				var cerr error
				if !pretend {
					_, cerr = fcopy(dst, path)
				}
				
				// If copy failed
				if cerr != nil {
					if verbose {
						fmt.Println(cerr.Error())
					}
					if collapseOnError {
						// Returning the error will cause the walking to stop
						return cerr
					}
				}
			}
		}
		
		return nil
	}
}

func main() {
	// Parse arguments
	// TODO move from 'arguments' to 'flag'
	p := flag.NewFlagSet("coppee", flag.ContinueOnError)
	p.SetOutput(ioutil.Discard)
	overwrite := p.Bool("o", false, "overwrite existing files")
	quiet     := p.Bool("q", false, "quiet mode - no verbose prints")
	pretend   := p.Bool("p", false, "pretend mode - does not copy files")
	err, args := p.Parse(os.Args[1:]), p.Args()
	
	if err != nil {
		fmt.Println("argument error: " + err.Error())
		return
	}
	if len(args) > 1 {
		fmt.Println("argument error: too many arguments.")
		return
	}
	if len(args) == 0 {
		// TODO add a parser-dependent default print
		fmt.Println("*** Premature version ***\n\n" +
				"Usage:\n" +
				"coppee [-o] [-q] [-p] <dir>\n\n" +
				"Arguments:\n" +
				"-o\tOverwrite existing target files. (default: false)\n" +
				"-q\tQuiet mode, disable verbose prints. (default: false)\n" +
				"-p\tPretend to copy, only print what will be copied. (default: false)\n" +
				"dir\tTarget directory. Must contain an instruction file named '.coppee'.")
		return
	}
	
	// Check directory
	inputDir := args[0]
	
	if stat, err := os.Stat(inputDir); os.IsNotExist(err) || !stat.IsDir() {
		fmt.Println("path error: '" + inputDir + "' is not a valid directory.")
		return
	}
	
	// Parse copy rules
	inputFile := filepath.Join(inputDir, ".coppee")
	rules, err := readRules(inputFile)
	if err != nil {
		fmt.Println("rule parse error: " + err.Error())
		return
	}
	
	// Copy files!
	walker := copier(inputDir, rules, true, *overwrite, !*quiet, *pretend)
	err = filepath.Walk(inputDir, walker)
	if err != nil {
		fmt.Println("copy error: " + err.Error())
		return
	}
	
	// Notify debug mode
	if *pretend {
		fmt.Println("*** PRETEND MODE. NO FILES WERE COPIED. ***")
	}
}






