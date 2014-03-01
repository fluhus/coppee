package main

import (
	"os"
	"fmt"
	"path/filepath"
	"arguments"
)

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
		relpath, relpatherr := filepath.Rel(basedir, path)
		if relpatherr != nil { panic(relpatherr.Error()) }

		for _,rule := range rules {
			// Check if filename matches
			// 'path' is source file name
			if globalMatch(rule.src, relpath) {
				// Destination file
				dst := rule.src.ReplaceAllString(path, rule.dst)

				// Overwrite?
				if !overwrite && fexists(dst) {
					continue
				}
				
				// Print
				if verbose {
					fmt.Println("copying: '" + path + "' to '" + dst + "'")
				}
				
				// Copy
				_, cerr := fcopy(dst, path)
				
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
	fmt.Println("Hi")
	r,e := readRules(".coppee")
	if e != nil { panic(e.Error()) }
	filepath.Walk(".", copier(".", r, true, true, true))

	arguments.NewParser()
	// var i int
	// flag.IntVar(&i, "integer", 13, "descriptionnn")
	// flag.Parse()
	// fmt.Println(i)
}

/*
TODOes:
1. make command line arguments
2. parse them so it can be run from command line
3. check if target dir exists

*/






