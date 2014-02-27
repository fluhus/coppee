package main

// General helper functions

import (
	"regexp"
	"errors"
	"os"
	"bufio"
	"io"
)

// Searches for a global match of the string to the given regexp.
// Returns true is matches, false if not.
func globalMatch(re *regexp.Regexp, s string) bool {
	return len(re.FindString(s)) == len(s)
}

// Checks whether a file exists.
// Returns true if exists, false if not.
func fexists(file string) bool {
	if _,e := os.Stat(file); os.IsNotExist(e) {
		return false
	}
	return true
}

// Copies a file.
// Returns nil if successful, or the relevant error if not.
func fcopy(dst, src string) (written int64, err error) {
	// Open input file
	in, ine := os.Open(src)
	if ine != nil {
		switch {
		case os.IsNotExist(ine):
			ine = errors.New("File not found: " + src)
		case os.IsPermission(ine):
			ine = errors.New("No permission to read: " + src)
		}
		
		return 0, ine
	}
	defer in.Close()
	
	// Open output file
	out, oute := os.Create(dst)
	if oute != nil {
		if os.IsPermission(oute) {
			oute = errors.New("No permission to write: " + dst)
		}
	
		return 0, oute
	}
	defer out.Close()
	
	// Open buffers
	bin := bufio.NewReader(in)
	bout := bufio.NewWriter(out)
	defer bout.Flush()
	
	// Go!
	return io.Copy(bout, bin)
}

