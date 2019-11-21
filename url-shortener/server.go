package main

import (
	"fmt"
	"log"
	"net/http"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var urlMap map[string]*string

func main() {
	// Reads YAML
	yamlFile, err := ioutil.ReadFile("redirect.yml")
	if err != nil {
		log.Printf("Redirect.yml err   #%v ", err)
	}
	// Parse YAML
	err = yaml.Unmarshal(yamlFile, &urlMap)
	if err != nil {
		log.Printf("YAML parse err #%v ", err)
	}
	// Custom Handler
	myMux := http.NewServeMux()
	for short, long := range urlMap {
		myMux.Handle(fmt.Sprintf("/%s", short), http.RedirectHandler(*long, 307))
	}
	log.Fatal(http.ListenAndServe(":8100", myMux))
}
