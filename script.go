package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type round struct {
	problem  string
	solution int
}

var correct int

func main() {

	lines := readCSV("./problems.csv")
	if len(lines) == 0 {
		fmt.Println("Error:", "No problems found in the file provided")
	}
	var data round
	correct = 0
	for i, line := range lines {
		answer, err := strconv.Atoi(line[1])
		if err != nil {
			fmt.Println("Invalid solution to problem no. ", strconv.Itoa(i))
			continue
		}
		data = round{
			problem:  line[0],
			solution: answer,
		}
		fmt.Println(data.problem)
		response := readUserInput()
		input, err := strconv.Atoi(response)
		if input == data.solution {
			correct++
			continue
		}
	}
	fmt.Println("You got", correct, "/", len(lines), "answers correct")
}

func readUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _, _ := reader.ReadLine()
	return string(text)
}

func readCSV(
	filename string,
) [][]string {
	rFile, err := os.Open(filename)
	defer rFile.Close()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	reader := csv.NewReader(rFile)
	lines, err := reader.ReadAll()
	if err == io.EOF {
		fmt.Println("Error: ", err)
		return [][]string{}
	}
	return lines
}
