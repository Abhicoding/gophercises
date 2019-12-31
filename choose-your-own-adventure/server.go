package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

func Server(storyMap map[string]Plot) {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", templateHandler)
	log.Fatal(http.ListenAndServe(":1111", nil))
}

func templateHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	t := template.New("index.html")
	t, _ = t.ParseFiles("./index.html")

	arc := strings.TrimPrefix(r.URL.Path, "/")
	if arc != "" {
		section := storyMap[arc]
		if err := t.Execute(w, section); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
		return
	}
	section := storyMap["intro"]
	if err := t.Execute(w, section); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
