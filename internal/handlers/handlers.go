package handlers

import (
	"html/template"
	"net/http"
)

// PageData represents the data passed to templates
type PageData struct {
	Title string
	Count int
}

// HomeHandler handles the main page request
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	data := PageData{
		Title: "HTMX + Go + Tailwind",
		Count: 0,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
