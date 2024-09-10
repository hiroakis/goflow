package main

import (
	"fmt"
	"os"

	"github.com/hiroakis/goflow"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <file> <function>\n", os.Args[0])
		return
	}

	fileName := os.Args[1]
	funcName := os.Args[2]
	if err := goflow.AnalyzeFunction(fileName, funcName); err != nil {
		fmt.Printf("goflow.AnalyzeFunction failed: %v", err)
		return
	}
	fmt.Println(goflow.GetUML())
}
