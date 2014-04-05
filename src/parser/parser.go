package parser

import (
	"os"
	"fmt"
	"bufio"
	"errors"
	"regexp"
	"strings"
)

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
func ReadRules(path string) (rules []CopyRule, err error) {
	// Open file
	f, ferr := os.Open(path)
	if ferr != nil {
		err = errors.New("file not found: " + path)
		return
	}

	b := bufio.NewReader(f)
	state := nextIsTemplate

	// Scan lines
	lineNumber := 0       // number of current line
	lastTemp := ""        // last encountered template
	lastTempLine := -1    // line of the last encountered template,
	                      // for error reporting
	for r, rerr := b.ReadString('\n'); rerr == nil; r, rerr = b.ReadString('\n') {
		lineNumber++
		
		// Trim spaces and BOM
		r = strings.Trim(r, charsToTrim)

		// Skip empty lines
		if len(r) == 0 { continue }

		// Skip comments
		if len(r) >= 2 && r[0:2] == "//" { continue }

		//TODO move rule making to the rules' territory
		// If expecting template, create regex
		if state == nextIsTemplate {
			lastTemp = r
			lastTempLine = lineNumber
		
			// Check for negation - denotd by '!'
			negated := false
			if r[0] == '!' {
				negated = true
				r = strings.Trim(r[1:], charsToTrim)
				
				// Only '!' without a template is not allowed
				if len(r) == 0 {
					err = errors.New(
							fmt.Sprintf("line %d: expected regular expression " +
							"after negation operator", lineNumber))
					return
				}
				
			// Else, remove escaping backslashes "\!"
			} else {
				escapeNegationRE := regexp.MustCompile("^\\\\(\\\\*!.*)")
				r = escapeNegationRE.ReplaceAllString(r, "$1")
			}
			
			// Verify regex
			if !isRegexp(r) {
				err = errors.New(fmt.Sprintf("line %d: invalid regular " +
						"expression: %s", lineNumber, r))
				return
			}
			
			// We're all done, add regex
			if negated {
				rules = append(rules, newNegatedRule())
			} else {
				rules = append(rules, newSimpleRule())
			}
			
			rules[len(rules) - 1].setTemplate(r)
			
		// If expecting target, add to last rule
		} else {
			rules[len(rules) - 1].setTarget(r)
		}
		
		state = !state
	}
	
	// Check that last rule has a target
	if state == nextIsTarget {
		err = errors.New(fmt.Sprintf("line %d: source with no target: %s",
				lastTempLine, lastTemp))
		return
	}

	return
}


