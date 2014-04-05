package parser

// Copy rules implementation

import (
	"os"
	"fmt"
	"regexp"
	"helpers"
)

// Represents a template and its target.
//
// The apply function applies the rule on the given file name, and checks
// whether or not it should be copied and to where. 'shouldCopy' will be true
// iff the file should be copied, and 'target' will be the name of the target
// file.
//
// The setter functions set the template and target of the rule. Template is
// asserted to be a valid regex.
type CopyRule interface {
	Apply(fileName string) (target string, shouldCopy bool)
	setTemplate(template string)
	setTarget(target string)
}

// A simple copy rule. Performs a simple regex match of the input and template.
type simpleRule struct {
	template *regexp.Regexp
	target   string
}

func (rule *simpleRule) Apply(fileName string) (target string, shouldCopy bool) {
	// Check for match
	if helpers.GlobalMatch(rule.template, fileName) {
		return rule.template.ReplaceAllString(fileName, rule.target), true
	} else {
		return "", false
	}
}

func (rule *simpleRule) setTemplate(template string) {
	rule.template = regexp.MustCompile(template)
}

func (rule *simpleRule) setTarget(target string) {
	rule.target = target
}

// Returns a simple rule instance.
func newSimpleRule() CopyRule {
	return &simpleRule{}
}

// Returns true on mismatches to the template.
type negatedRule simpleRule

// Path separator for regular expressions
var pathSeparatorRegex = regexp.QuoteMeta(string(os.PathSeparator))

// Captures a files path and splits it as follows:
// ${0} will match the entire file path and name.
// ${1} will match the directory path preceding the file name.
// ${2} will match the file's prefix (until the last period).
// ${3} will match the file's suffix, including the period.
var negatedCaptureRegex = regexp.MustCompile(fmt.Sprintf(
		"^((?:[^%s]*%s)*)([^\\.]*(?:\\.[^\\.]*)*?)(\\.[^\\.]*)?$",
		pathSeparatorRegex, pathSeparatorRegex))

// Capturing groups available for the target:
// ${0} will match the entire file path and name.
// ${1} will match the directory path preceding the file name.
// ${2} will match the file's prefix (until the last period).
// ${3} will match the file's suffix, including the period.
func (rule *negatedRule) Apply(fileName string) (target string, shouldCopy bool) {
	// Check for mismatch
	if !helpers.GlobalMatch(rule.template, fileName) {
		return negatedCaptureRegex.ReplaceAllString(fileName, rule.target), true
	} else {
		return "", false
	}
}

func (rule *negatedRule) setTemplate(template string) {
	(*simpleRule)(rule).setTemplate(template)
}

func (rule *negatedRule) setTarget(target string) {
	(*simpleRule)(rule).setTarget(target)
}

func newNegatedRule() CopyRule {
	return &negatedRule{}
}

