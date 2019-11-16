package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type round struct {
	problem  string
	solution int
}

var correct int
var timer = 5

func main() {

	lines := readCSV("./problems.csv")
	if len(lines) == 0 {
		fmt.Println("Error:", "No problems found in the file provided")
	}

	fmt.Println(
		fmt.Sprintf(`Every question has a %d sec timer. If it expires, then the quiz is over 
Press Enter/Return key to start the quiz.`,
			timer,
		))
	readUserInput()
	var data round
	message := make(chan string)
	correct = 0
problemloop:
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
		go func() {
			fmt.Println("Awaiting input ...")
			response := readUserInput()
			message <- response
		}()
		select {
		case value := <-message:
			input, err := strconv.Atoi(value)
			if err != nil {
				continue problemloop
			}
			if input == data.solution {
				correct++
				continue problemloop
			}
		case <-time.After(time.Duration(timer) * time.Second):
			fmt.Println("Timer expired")
			break problemloop
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
