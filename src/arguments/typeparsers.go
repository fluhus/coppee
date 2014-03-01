package arguments

// Parsers for specific types.

import (
	"errors"
	"strconv"
)

// Parses a value argument.
type valueParser interface {
	parse(arg string) error
}

// Parses a boolean argument (that is not followed by a value).
type noValueParser interface {
	parse() error
}

// Common structure of a parser.
type generalParser struct {
	symbol string        // argument that matches this parser
	description string   // description for help print
	parsed bool          // indicates that this argument was already encountered
}

// Parses integers.
type intParser struct {
	generalParser    // embedded common parser fields
	defaultVal int   // default value
	p *int           // output pointer
}

// Returns a new integer parser.
func newIntParser(symbol string, description string, defaultVal int) *intParser {
	i := defaultVal
	return &intParser{ generalParser{symbol, description, false}, defaultVal, &i }
}

// Parses the given argument. Returns an error if fails.
func (p *intParser) parse(arg string) error {
	if p.parsed {
		return errors.New("too many occurrances of argument '" + p.symbol + "'.")
	}
		
	p.parsed = true
		
	// Try to parse
	i, err := strconv.Atoi(arg)
	if err != nil {
		return errors.New("'" + arg + "' is not a valid integer.")
	}
	
	// Assign output
	*(p.p) = i
	return nil
}

// Parses floats.
type floatParser struct {
	generalParser      // embedded common parser fields
	defaultVal float64 // default value
	p *float64         // output pointer
}

// Returns a new float parser.
func newFloatParser(symbol string, description string, defaultVal float64) *floatParser {
	f := defaultVal
	return &floatParser{ generalParser{symbol, description, false}, defaultVal, &f }
}

// Parses the given argument. Returns an error if fails.
func (p *floatParser) parse(arg string) error {
	if p.parsed {
		return errors.New("too many occurrances of argument '" + p.symbol + "'.")
	}
		
	p.parsed = true
		
	// Try to parse
	f, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return errors.New("'" + arg + "' is not a valid float.")
	}
	
	// Assign output
	*(p.p) = f
	return nil
}

// Parses strings.
type stringParser struct {
	generalParser      // embedded common parser fields
	defaultVal string  // default value
	p *string          // output pointer
}

// Returns a new string parser.
func newStringParser(symbol string, description string, defaultVal string) *stringParser {
	s := defaultVal
	return &stringParser{ generalParser{symbol, description, false}, defaultVal, &s }
}

// Parses the given argument. Returns an error if fails.
func (p *stringParser) parse(arg string) error {
	if p.parsed {
		return errors.New("too many occurrances of argument '" + p.symbol + "'")
	}
		
	p.parsed = true
	
	// Assign output
	*(p.p) = arg
	return nil
}

// Parses boolean arguments (with no following value).
type boolParser struct {
	generalParser    // embedded common parser fields
	defaultVal bool  // default value
	p *bool          // output pointer
}

// Returns a new bool parser.
func newBoolParser(symbol string, description string, defaultVal bool) *boolParser {
	b := defaultVal
	return &boolParser{ generalParser{symbol, description, false}, defaultVal, &b }
}

// Reacts to the argument. Returns an error if fails.
func (p *boolParser) parse() error {
	if p.parsed {
		return errors.New("too many occurrances of argument '" + p.symbol + "'.")
	}
		
	p.parsed = true
	
	// Assign output
	*(p.p) = !*(p.p)
	return nil
}
