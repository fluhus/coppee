package main

import (
	"os"
	"fmt"
	"bufio"
	"errors"
	"regexp"
	"strings"
)

// Holds a template source file and its destination.
type copyRule struct {
	src *regexp.Regexp  // source file template regex
	dst string          // destination path
	negated bool        // if true, will copy non-matching files
}

// String representation of a rule. For debugging.
func (cr copyRule) String() string {
	negated := ""
	if cr.negated { negated = " negated" }
	return fmt.Sprintf("(\"%s\" \"%s\"%s)", cr.src.String(), cr.dst, negated)
}

// Checks that given string is a valid regular expression.
// Returns true if valid, false if invalid.
func isRegexp(s string) bool {
	_,err := regexp.Compile(s)
	return err == nil
}

// Describes the current state of the reader.
type readerState bool

// Allowed states
const (
	nextIsTemplate readerState = true
	nextIsTarget   readerState = false
)

// Whitespace characters and BOM to be trimmed
// before parsing an input line.
const charsToTrim = " \t\n\r\xef\xbb\xbf"

// Reads copy rules from a config file. Returns an error if file not found
// or badly formatted.
func readRules(path string) (rules []copyRule, err error) {
	// Open file
	f, ferr := os.Open(path)
	if ferr != nil {
		err = errors.New("file not found: " + path)
		return
	}

	b := bufio.NewReader(f)
	state := nextIsTemplate

	// Scan lines
	for r, rerr := b.ReadString('\n'); rerr == nil; r, rerr = b.ReadString('\n') {
		// Trim spaces and BOM
		r = strings.Trim(r, charsToTrim)

		// Skip empty lines
		if len(r) == 0 { continue }

		// Skip comments
		if len(r) >= 2 && r[0:2] == "//" { continue }

		// If expecting template, create regex
		if state == nextIsTemplate {
		
			// Check for negation - denotd by '!'
			negated := false
			if r[0] == '!' {
				negated = true
				r = strings.Trim(r[1:], charsToTrim)
				
				// Only '!' without a template is not allowed
				if len(r) == 0 {
					err = errors.New("expected regular expression after " +
							"negation operator")
					return
				}
				
			// Else, remove escaping backslashes "\!"
			} else {
				escapeNegationRE := regexp.MustCompile("^\\\\(\\\\*!.*)")
				r = escapeNegationRE.ReplaceAllString(r, "$1")
			}
			
			// Verify regex
			if !isRegexp(r) {
				err = errors.New("invalid regular expression: " + r)
				return
			}
			
			// We're all done, add regex
			rules = append(rules, copyRule{regexp.MustCompile(r), "", negated})
			
			state = !state
			
		// If expecting target, add to last rule
		} else {
			rules[len(rules) - 1].dst = r
			state = !state
		}
	}
	
	// Check that last rule has a target
	if rules[len(rules) - 1].dst == "" {
		err = errors.New("source with no target: " +
				rules[len(rules) - 1].src.String())
		return
	}

	return
}


