package main

import (
	"fmt"
	"os"
)

var storyMap map[string]Plot

func main() {
	filesToRead := readFileNamesFromArgs()
	if len(filesToRead) < 3 {
		fmt.Println("Please provide program file and plot file")
		return
	}

	plotFileName := filesToRead[2]
	file, err := os.Open(plotFileName)
	if err != nil {
		fmt.Println(err, "\nDefaulting to plot.json ...")
		file, _ = os.Open("plot.json")
	}
	storyMap, err := readJSON(file)
	if err != nil {
		fmt.Println("Bad plot file: ", err)
		return
	}
	switch filesToRead[1] {
	case "server":
		Server(*storyMap)
		return
	default:
		fmt.Println("Defaulting to terminal adventure game ...\n\n")
		Script(*storyMap)
		return
	}
}
