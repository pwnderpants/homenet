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

// IncrementHandler handles the counter increment request
func IncrementHandler(w http.ResponseWriter, r *http.Request) {
	// This is a simple example - in a real app you'd use sessions or database
	// For demo purposes, we'll just return a new counter value
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
		<div class="text-center">
			<h2 class="text-2xl font-bold text-white mb-4">Counter: <span id="counter">1</span></h2>
			<button 
				hx-post="/increment" 
				hx-target="#counter-section"
				class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition-colors duration-200">
				Increment
			</button>
		</div>
	`))
}
