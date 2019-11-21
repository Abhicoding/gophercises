package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var urlMap map[string]*string

func main() {
	yamlFile, err := ioutil.ReadFile("redirect.yml")
	if err != nil {
		log.Printf("Redirect.yml err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &urlMap)
	if err != nil {
		log.Printf("YAML parse err #%v ", err)
	}

	myMux := http.NewServeMux()
	myMux.HandleFunc("/hello", hello)
	for short, long := range urlMap {
		myMux.Handle(fmt.Sprintf("/%s", short), http.RedirectHandler(*long, 307))
	}
	log.Fatal(http.ListenAndServe(":8100", myMux))
}

func hello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hellow!\n")
}
