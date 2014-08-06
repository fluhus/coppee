// An automatic file copier. Execute with no parameters for explanation on the
// accepted parameters. See the readme for detailed usage instructions.
package main

import (
	"os"
	"fmt"
	"flag"
	"parser"
	"walker"
	"io/ioutil"
	"path/filepath"
)

func main() {
	// Parse arguments
	p := flag.NewFlagSet("coppee", flag.ContinueOnError)
	p.SetOutput(ioutil.Discard)
	overwrite := p.Bool("o", false, "overwrite existing files")
	quiet     := p.Bool("q", false, "quiet mode - no verbose prints")
	pretend   := p.Bool("p", false, "pretend mode - does not copy files")
	inputFile := p.String("f", "", "instruction file (relative to shell path)")
	err, args := p.Parse(os.Args[1:]), p.Args()
	
	if err != nil {
		fmt.Println("argument error: " + err.Error())
		return
	}
	if len(args) > 1 {
		fmt.Println("argument error: too many arguments.")
		return
	}
	if len(args) == 0 {
		// TODO add a parser-dependent default print
		fmt.Println("*** Premature version " + version + " ***\n\n" +
				"Usage:\n" +
				"coppee [-o] [-q] [-p] [-f] <dir>\n\n" +
				"Arguments:\n" +
				"-o\toverwrite existing target files. (default: false)\n" +
				"-q\tquiet mode, disable verbose prints. (default: false)\n" +
				"-p\tpretend to copy, only print what will be copied. (default: false)\n" +
				"-f\tinput instruction file. (default: <dir>\\.coppee)\n" +
				"dir\ttarget directory.")
		return
	}
	
	// Check directory
	inputDir := args[0]
	
	if stat, err := os.Stat(inputDir); os.IsNotExist(err) || !stat.IsDir() {
		fmt.Println("path error: '" + inputDir + "' is not a valid directory.")
		return
	}
	
	// Parse copy rules
	if *inputFile == "" {
		*inputFile = filepath.Join(inputDir, ".coppee")
	}
	
	rules, err := parser.ReadRules(*inputFile)
	if err != nil {
		fmt.Println("rule parse error: " + err.Error())
		return
	}
	
	// Copy files!
	walker := walker.Copier(inputDir, rules, true, *overwrite, !*quiet, *pretend)
	err = filepath.Walk(inputDir, walker)
	if err != nil {
		fmt.Println("copy error: " + err.Error())
		return
	}
	
	// Notify debug mode
	if *pretend {
		fmt.Println("*** PRETEND MODE. NO FILES WERE COPIED. ***")
	}
}






