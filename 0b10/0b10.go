package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"log"
	"net/http"
)

// Compo == Component
// Component is customisable, independent and reusable UI element.
// It works by embedding app.Compo into a struct.
// This means that I can define UI elements for myself.
type Monika struct {
	app.Compo
}

// the (h *Monika) simply tells that the struct has the Render method. This is pretty nice.
func (h *Monika) Render() app.UI {
	return app.H1().Text("Just Monika")
}

// main is the entry point for the server.
func main() {
	// We simply give to the "/" route the struct Monika back. Struct Monika is an app Component that returns a header with text Just Monika.
	// What we end up seeing is the header Just Monika on the localhost:8000/ page.
	app.Route("/", func() app.Composer { return &Monika{} })

	// Run
	app.RunWhenOnBrowser()

	// by the net/http we create the handler for the route "/". We simply give the name and description of the page.
	http.Handle("/", &app.Handler{
		Name:        "Monika",
		Description: "There can be only one",
	})

	// If there is an error, log it and shut down the server
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
