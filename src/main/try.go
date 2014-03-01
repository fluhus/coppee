package main

import (
	"fmt"
	"arguments"
	"os"
)

func main() {
	p := arguments.NewParser()
	s := p.AddString("-s", "", "aaa")
	
	args, err := p.Parse(os.Args[1:])
	
	fmt.Println(*s)
	fmt.Println(args)
	fmt.Println(err)
}

