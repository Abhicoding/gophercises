package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
)

func readFileNamesFromArgs() []string {
	file := os.Args
	if len(os.Args) > 1 {
		return file
	}
	return nil
}

func readJSON(file *os.File) (*map[string]Plot, error) {
	var plot map[string]Plot
	data, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(data), &plot)
	return &plot, nil
}

func readUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _, _ := reader.ReadLine()
	return string(text)
}
