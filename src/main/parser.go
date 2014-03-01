package main

import (
	"os"
	//"fmt"
	"bufio"
	"errors"
	"regexp"
	"strings"
)

// Holds a template source file and its destination.
type copyRule struct {
	src *regexp.Regexp  // source file template regex
	dst string  // destination path
}

func (cr copyRule) String() string {
	return "{" + cr.src.String() + " " + cr.dst + "}"
}

// Checks that given string is a valid regular expression.
// Returns true if valid, false if invalid.
func isRegexp(s string) bool {
	_,err := regexp.Compile(s)
	return err == nil
}

// Reads copy rules from a config file.
func readRules(path string) (rules []copyRule, err error) {
	// Open file
	f, ferr := os.Open(path)
	if ferr != nil {
		err = errors.New("File not found: " + path)
		return
	}

	b := bufio.NewReader(f)

	// Scan lines
	var s []string
	for r, rerr := b.ReadString('\n'); rerr == nil; r, rerr = b.ReadString('\n') {
		// Trim spaces
		r = strings.Trim(r, " \t\n")

		// Skip empty line
		if len(r) == 0 { continue }

		// Skip comments
		if len(r) >= 2 && r[0:2] == "//" { continue }

		// Check regex
		if !isRegexp(r) {
			err = errors.New("Invalid regular expression: " + r)
			return
		}

		s = append(s, r)
	}
	
	// Check length - each source must be followed by a target, so an odd length
	// is not allowed.
	if len(s) % 2 == 1 {
		err = errors.New("Source with no target: " + s[len(s) - 1])
		return
	}

	// Create rules
	for i := 0; i < len(s); i += 2 {
		rules = append(rules, copyRule{regexp.MustCompile(s[i]), s[i+1]})
	}

	return
}


