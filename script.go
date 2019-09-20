package main

import (
	"fmt"
	"encoding/csv"
	"os"
	"io"
)
func main() {
	rFile, err := os.Open("./problems.csv")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer rFile.Close()
	reader := csv.NewReader(rFile)
	lines, err := reader.Read()
	if err == io.EOF {
		fmt.Println("Error: ", err)
	}
	fmt.Println(lines)
}
