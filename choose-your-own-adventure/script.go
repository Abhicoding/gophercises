package main

import (
	"fmt"
	"strconv"
)

var current Plot

func Script(storyMap map[string]Plot) {
	fileToRead := readFileNamesFromArgs()
	if fileToRead == nil {
		fmt.Println("Please provide a story file")
		return
	}
	current := storyMap["intro"]
	printPlot(current)
	c := make(chan int)
	for {
		if len(current.Options) == 0 {
			fmt.Println("\n\n ** The End **")
			return
		}
		go readInput(c)
		response := <-c
		if len(current.Options) < response || response < 0 {
			fmt.Println("Incorrect option value")
			printOptions(current)
			continue
		}
		current = storyMap[current.Options[response-1].Arc]
		printPlot(current)
	}
}

func readInput(c chan int) {
	var i int
	var err error
	for {
		i, err = strconv.Atoi(readUserInput())
		if err != nil {
			fmt.Println("Enter the number against the options")
			continue
		}
		break
	}
	c <- i
}

func printPlot(plot Plot) {
	fmt.Println("\n - \n")
	for _, tale := range plot.Story {
		fmt.Println(tale)
	}
	printOptions(plot)
}

func printOptions(plot Plot) {
	for i, option := range plot.Options {
		fmt.Printf("%d > %s\n", i+1, option.Text)
	}
}
