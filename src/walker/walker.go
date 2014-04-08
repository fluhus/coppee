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
// baseDir: base walking directory, must match the directory walked
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
		baseDir string,
		rules []parser.CopyRule,
		collapseOnError bool,
		overwrite bool,
		verbose bool,
		pretend bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// Path - relative to shell directory.
		// Relative path - relative to walked directory.
		// Matching will be to relative path.
		relPath, relPathErr := filepath.Rel(baseDir, path)
		if relPathErr != nil { panic(relPathErr.Error()) } // should not happen
		
		// Don't try to copy directories
		if info.IsDir() {
			return nil
		}

		for _,rule := range rules {
			// Apply rule
			relTarget, shouldCopy := rule.Apply(relPath)  // relative to walked dir
			target := filepath.Join(baseDir, relTarget)   // relative to shell dir
			
			if shouldCopy {
				// Overwrite?
				if !overwrite && helpers.FileExists(target) {
					if verbose {
						fmt.Println("skipping: '" + path + "' to '" + target + "'")
					}
					continue
				}
				
				// Print
				if verbose {
					fmt.Println("copying:  '" + path + "' to '" + target + "'")
				}
				
				// Copy
				var cerr error
				if !pretend {
					_, cerr = helpers.FileCopy(target, path)
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
