package main

import (
    "fmt"
	"log"
    "net/http"
	"text/template"
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// stuff to set headers so that clients will be like: OH wow, this is server sent events waoow
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Access-Control-Allow-Origin", "*")

	// Presumably this below must be in a loop and needs to be flushed.
	fmt.Fprintf(w, "<div id=\"chat\">HOI</div>")

	// Because here we simply close the server sent event, but this needs to be active at all times.
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error while parsing index.html!\nMessage: %s", err)
		return
	}
	data := struct {
		Name string
	} {
		Name: "Monika",
	}
	template.Execute(w, data)
}

func main() {
	http.HandleFunc("/", indexHandler)
    http.HandleFunc("/events", sseHandler)
	fmt.Println("http://127.0.0.1:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println(err.Error())
    }
}
