// Exports the walk function that traverses the directory tree
// and copies the appropriate files.
package walker

import (
	"path/filepath"
	"parser"
	"os"
	"fmt"
	"helpers"
)

// Returns a copier walk function.
//
// basedir: base walking directory, must match the directory walked
//
// rules: array of copying rules by whom to copy
//
// collapseOnError: if true, process will stop upon copying error
//
// overwrite: if true, will overwrite existing files
//
// verbose: if true, will print copies made and errors
//
// pretend: if true, will only pretend to copy files
func Copier(
		basedir string,
		rules []parser.CopyRule,
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
			if helpers.GlobalMatch(rule.Src, relPath) {
				// Destination file
				dst := rule.Src.ReplaceAllString(path, rule.Dst)

				// Overwrite?
				if !overwrite && helpers.FExists(dst) {
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
					_, cerr = helpers.FCopy(dst, path)
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
