// Parses command line arguments.
package arguments

// BUG(amit): Improve help print for floats.
// BUG(amit): Add free argument parsers.

import (
	"fmt"
	"errors"
)

// Holds argument settings and does the parsing.
type Parser struct {
	argOrder []string                   // holds the order of received symbols
	val      map[string]valueParser     // from symbol to corresponding value-argument parser
	noval    map[string]noValueParser   // from symbol to corresponding no-value-argument parser
}

// Creates a new empty parser.
func NewParser() *Parser {
	return &Parser{ nil,
		make(map[string]valueParser),
		make(map[string]noValueParser) }
}

// Adds an integer argument.
// symbol:      string that matches this argument.
// description: argument description for help prints.
// defaultVal:  default value for when the argument is not given.
// Returned pointer will have the argument's value after calling Parse().
func (p *Parser) AddInt(
		symbol string,
		description string,
		defaultval int) *int {
		
	// Empty symbol is not accepted
	if symbol == "" { panic("Cannot add an empty symbol.") }
	
	// Existing symbols are not accepted
	_, ok1 := p.val[symbol]
	_, ok2 := p.noval[symbol]
	if ok1 || ok2 {
		panic("Symbol '" + symbol + "' already exists.")
	}
	
	// Add to parser
	p.argOrder = append(p.argOrder, symbol)
	parser := newIntParser(symbol, description, defaultval)
	p.val[symbol] = parser
	
	return parser.p
}

// Adds a float argument.
// symbol:      string that matches this argument.
// description: argument description for help prints.
// defaultVal:  default value for when the argument is not given.
// Returned pointer will have the argument's value after calling Parse().
func (p *Parser) AddFloat(
		symbol string,
		description string,
		defaultval float64) *float64 {
		
	// Empty symbol is not accepted
	if symbol == "" { panic("Cannot add an empty symbol.") }
	
	// Existing symbols are not accepted
	_, ok1 := p.val[symbol]
	_, ok2 := p.noval[symbol]
	if ok1 || ok2 {
		panic("Symbol '" + symbol + "' already exists.")
	}
	
	// Add to parser
	p.argOrder = append(p.argOrder, symbol)
	parser := newFloatParser(symbol, description, defaultval)
	p.val[symbol] = parser
	
	return parser.p
}

// Adds a string argument.
// symbol:      string that matches this argument.
// description: argument description for help prints.
// defaultVal:  default value for when the argument is not given.
// Returned pointer will have the argument's value after calling Parse().
func (p *Parser) AddString(
		symbol string,
		description string,
		defaultval string) *string {
		
	// Empty symbol is not accepted
	if symbol == "" { panic("Cannot add an empty symbol.") }
	
	// Existing symbols are not accepted
	_, ok1 := p.val[symbol]
	_, ok2 := p.noval[symbol]
	if ok1 || ok2 {
		panic("Symbol '" + symbol + "' already exists.")
	}
	
	// Add to parser
	p.argOrder = append(p.argOrder, symbol)
	parser := newStringParser(symbol, description, defaultval)
	p.val[symbol] = parser
	
	return parser.p
}

// Adds a bool argument. This argument does not expect a value after its
// symbol.
// symbol:      string that matches this argument.
// description: argument description for help prints.
// defaultVal:  default value for when the argument is not given.
// Returned pointer will have the argument's value after calling Parse().
func (p *Parser) AddBool(
		symbol string,
		description string,
		defaultval bool) *bool {
		
	// Empty symbol is not accepted
	if symbol == "" { panic("Cannot add an empty symbol.") }
	
	// Existing symbols are not accepted
	if p.has(symbol) {
		panic("Symbol '" + symbol + "' already exists.")
	}
	
	// Add to parser
	p.argOrder = append(p.argOrder, symbol)
	parser := newBoolParser(symbol, description, defaultval)
	p.noval[symbol] = parser
	
	return parser.p
}

// Returns true iff a value argument with the given symbol is defined.
func (p *Parser) hasVal(symbol string) bool {
	_,ok := p.val[symbol]
	return ok
}

// Returns true iff a no-value (bool) argument with the given symbol is defined.
func (p *Parser) hasNoVal(symbol string) bool {
	_,ok := p.noval[symbol]
	return ok
}

// Returns true iff a value or no-value (bool) argument with the given symbol
// is defined.
func (p *Parser) has(symbol string) bool {
	return p.hasVal(symbol) || p.hasNoVal(symbol)
}

// Parses the given arguments (usually called on os.Args[1:]).
// freeArgs will contain the free arguments, that don't belong to any defined
// flag.
// err will be nil if parsing went ok. If not err will be defined with an
// informative user-printable error.
func (p *Parser) Parse(args []string) (freeArgs []string, err error) {
	for i := 0; i < len(args); i++ {
		// Check if value argument
		if p.hasVal(args[i]) {
			// Check if a following value exists (not last and next is not
			// a symbol)
			if i == len(args) - 1 || p.has(args[i+1]) {
				// No value exists -> error
				err = errors.New("value expected after '" + args[i] + "'.")
				return
			}
			
			// Try to parse
			parseErr := p.val[args[i]].parse(args[i+1])
			if parseErr != nil {
				err = parseErr
				return
			}
			
			// Skip next argument
			i++
			
		// Check if no-value argument
		} else if p.hasNoVal(args[i]) {
			// Try to parse
			parseErr := p.noval[args[i]].parse()
			if parseErr != nil {
				err = parseErr
				return
			}
			
		// Free argument
		} else {
			freeArgs = append(freeArgs, args[i])
		}
	}
	
	return
}

// Creates a sequence of n spaces.
// Inefficient, so don't use regularly.
func spaces(n int) string {
	s := ""
	
	for i := 0; i < n; i++ {
		s += " "
	}
	
	return s
}

// Returns a user-printable help message for the defined
// arguments.
func (p *Parser) HelpPrint() string {
	// Find length of longest symbol, for justification
	longest := 0
	for _,v := range p.argOrder {
		if len(v) > longest {
			longest = len(v)
		}
	}
	
	// Print arguments by their order
	result := ""
	for i,v := range p.argOrder {
		// If more that 1 argument, add new line
		if i > 0 {
			result += "\n"
		}
		
		// Find argument's description
		var description string
		if p.hasVal(v) {
			description = p.val[v].helpPrint()
		} else if p.hasNoVal(v) {
			description = p.noval[v].helpPrint()
		} else {
			panic("argument '" + v + "' found but not defined")
		}
		
		// Add i'th argument
		result += fmt.Sprintf("%s%s: %s",
				spaces(2 + longest - len(v)),
				v,
				description)
	}
	
	return result
}

