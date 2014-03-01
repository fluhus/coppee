package main

import (
	"os"
	"fmt"
	"path/filepath"
	"arguments"
)

// If true, will only pretend to copy the files. For debugging.
const pretend = false

// Returns a copier walk-function.
// basedir:          base walking directory
// rules:            copying rules
// collapseOnError:  if true, process will stop upon copying error
// overwrite:        if true, will overwrite existing files
// verbose:          if true, will print copies made and errors
func copier(
		basedir string,
		rules []copyRule,
		collapseOnError bool,
		overwrite bool,
		verbose bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// Relative path - matching will be to relative path
		relPath, relPathErr := filepath.Rel(basedir, path)
		if relPathErr != nil { panic(relPathErr.Error()) } // should not happen

		for _,rule := range rules {
			// Check if filename matches
			// 'path' is source file name
			if globalMatch(rule.src, relPath) {
				// Destination file
				dst := rule.src.ReplaceAllString(path, rule.dst)

				// Overwrite?
				if !overwrite && fexists(dst) {
					if verbose {
						fmt.Println("skipping: '" + path + "'")
					}
					continue
				}
				
				// Print
				if verbose {
					fmt.Println("copying: '" + path + "' to '" + dst + "'")
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
	// Accepted argument format: TODO
	p := arguments.NewParser()
	args, err := p.Parse(os.Args[1:])
	
	if err != nil {
		fmt.Println("argument error: " + err.Error())
		return
	}
	if len(args) > 1 {
		fmt.Println("argument error: too many arguments.")
		return
	}
	if len(args) == 0 {
		fmt.Println("Premature version.\nWill overwrite existing files." +
				"\n\nUsage:\ncoppee <path>")
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
	walker := copier(inputDir, rules, true, true, true)
	err = filepath.Walk(inputDir, walker)
	if err != nil {
		fmt.Println("copy error: " + err.Error())
		return
	}
}







