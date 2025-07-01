package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/pwnderpants/homenet/internal/database"
)

// Use the Movie struct from database package
type Movie = database.Movie

// PageData represents the data passed to templates
type PageData struct {
	Title string
	Count int
}

// MovieBoardData represents the data for the movie board page
type MovieBoardData struct {
	Title      string
	Movies     []Movie
	MovieCount int
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

// MovieBoardHandler handles the movie board page request
func MovieBoardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/movie-board.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// Get movies from database
	movies, err := database.GetAllMovies()
	if err != nil {
		http.Error(w, "Failed to load movies: "+err.Error(), http.StatusInternalServerError)

		return
	}

	data := MovieBoardData{
		Title:      "Movie Board",
		Movies:     movies,
		MovieCount: len(movies),
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// AddMovieHandler handles adding a new movie
func AddMovieHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	title := r.FormValue("title")
	yearStr := r.FormValue("year")
	genre := r.FormValue("genre")
	streaming := r.FormValue("streaming")
	notes := r.FormValue("notes")
	imdbLink := r.FormValue("imdb_link")

	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)

		return
	}

	// Parse year
	year := 0
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			http.Error(w, "Invalid year", http.StatusBadRequest)

			return
		}
	}

	// Create new movie
	newMovie := Movie{
		Title:     title,
		Year:      year,
		Genre:     genre,
		Streaming: streaming,
		Notes:     notes,
		IMDBLink:  imdbLink,
	}

	// Add to database
	movieID, err := database.AddMovie(newMovie)
	if err != nil {
		http.Error(w, "Failed to add movie: "+err.Error(), http.StatusInternalServerError)

		return
	}

	// Set the ID for HTML generation
	newMovie.ID = movieID

	// Get all movies from database to return the complete updated list
	allMovies, err := database.GetAllMovies()
	if err != nil {
		http.Error(w, "Failed to get movies: "+err.Error(), http.StatusInternalServerError)

		return
	}

	// Return the complete movie list HTML for HTMX to replace
	w.Header().Set("Content-Type", "text/html")

	// Generate HTML for all movies
	var moviesHTML string
	for _, movie := range allMovies {
		movieHTML := `
		<div class="bg-gray-700 rounded-lg p-4 border border-gray-600">
			<div class="flex justify-between items-start">
				<div class="flex-1">
					<h4 class="text-lg font-semibold text-white">` + movie.Title + `</h4>
					<div class="flex items-center space-x-4 mt-2 text-sm text-gray-300">`

		if movie.Year > 0 {
			movieHTML += `<span class="bg-gray-600 px-2 py-1 rounded">` + strconv.Itoa(movie.Year) + `</span>`
		}

		if movie.Genre != "" {
			movieHTML += `<span class="bg-blue-600 px-2 py-1 rounded">` + movie.Genre + `</span>`
		}

		if movie.Streaming != "" {
			movieHTML += `<span class="bg-green-600 px-2 py-1 rounded">` + movie.Streaming + `</span>`
		}

		movieHTML += `
					</div>`

		if movie.IMDBLink != "" {
			movieHTML += `<div class="mt-2"><a href="` + movie.IMDBLink + `" target="_blank" class="text-blue-400 hover:text-blue-300 text-sm">View on IMDB</a></div>`
		}

		if movie.Notes != "" {
			movieHTML += `<p class="text-gray-400 mt-2 text-sm">` + movie.Notes + `</p>`
		}

		movieHTML += `
				</div>
				<div class="flex space-x-2">
					<button 
						data-movie-id="` + strconv.Itoa(movie.ID) + `"
						data-movie-title="` + movie.Title + `"
						data-movie-year="` + strconv.Itoa(movie.Year) + `"
						data-movie-genre="` + movie.Genre + `"
						data-movie-streaming="` + movie.Streaming + `"
						data-movie-notes="` + movie.Notes + `"
						data-movie-imdb="` + movie.IMDBLink + `"
						onclick="openEditModal(this)"
						class="text-blue-400 hover:text-blue-300 transition-colors duration-200">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
						</svg>
					</button>
					<button 
						hx-delete="/movie-board/delete/` + strconv.Itoa(movie.ID) + `"
						hx-target="closest div"
						hx-swap="outerHTML"
						hx-confirm="Are you sure you want to delete this movie?"
						class="text-red-400 hover:text-red-300 transition-colors duration-200">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
						</svg>
					</button>
				</div>
			</div>
		</div>`

		moviesHTML += movieHTML
	}

	// If no movies, show the "no movies" message
	if len(allMovies) == 0 {
		moviesHTML = `
		<div class="text-center py-8">
			<svg class="w-16 h-16 text-gray-600 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2m-9 0h10m-10 0a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V6a2 2 0 00-2-2"></path>
			</svg>
			<p class="text-gray-400">No movies added yet. Add your first movie above!</p>
		</div>`
	}

	// Create the response with out-of-band update for movie count
	response := moviesHTML + `
	<div id="movie-count" hx-swap-oob="true">
		` + strconv.Itoa(len(allMovies)) + ` movies in your list
	</div>`

	w.Write([]byte(response))
}

// DeleteMovieHandler handles deleting a movie
func DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	// Extract movie ID from URL
	path := r.URL.Path
	idStr := path[len("/movie-board/delete/"):]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)

		return
	}

	// Delete from database
	err = database.DeleteMovie(id)

	if err != nil {
		http.Error(w, "Failed to delete movie: "+err.Error(), http.StatusInternalServerError)

		return
	}

	// Return success status
	w.WriteHeader(http.StatusOK)
}

// EditMovieHandler handles editing an existing movie
func EditMovieHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	// Parse form data
	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	idStr := r.FormValue("id")
	title := r.FormValue("title")
	yearStr := r.FormValue("year")
	genre := r.FormValue("genre")
	streaming := r.FormValue("streaming")
	notes := r.FormValue("notes")
	imdbLink := r.FormValue("imdb_link")

	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)

		return
	}

	// Parse movie ID
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)

		return
	}

	// Parse year
	year := 0
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			http.Error(w, "Invalid year", http.StatusBadRequest)
			return
		}
	}

	// Update movie in database
	updatedMovie := Movie{
		ID:        id,
		Title:     title,
		Year:      year,
		Genre:     genre,
		Streaming: streaming,
		Notes:     notes,
		IMDBLink:  imdbLink,
	}

	err = database.UpdateMovie(updatedMovie)
	if err != nil {
		http.Error(w, "Failed to update movie: "+err.Error(), http.StatusInternalServerError)

		return
	}

	// Get all movies from database to return the complete updated list
	allMovies, err := database.GetAllMovies()
	if err != nil {
		http.Error(w, "Failed to get movies: "+err.Error(), http.StatusInternalServerError)

		return
	}

	// Return the complete movie list HTML for HTMX to replace
	w.Header().Set("Content-Type", "text/html")

	// Generate HTML for all movies
	var moviesHTML string

	for _, movie := range allMovies {
		movieHTML := `
		<div class="bg-gray-700 rounded-lg p-4 border border-gray-600">
			<div class="flex justify-between items-start">
				<div class="flex-1">
					<h4 class="text-lg font-semibold text-white">` + movie.Title + `</h4>
					<div class="flex items-center space-x-4 mt-2 text-sm text-gray-300">`

		if movie.Year > 0 {
			movieHTML += `<span class="bg-gray-600 px-2 py-1 rounded">` + strconv.Itoa(movie.Year) + `</span>`
		}

		if movie.Genre != "" {
			movieHTML += `<span class="bg-blue-600 px-2 py-1 rounded">` + movie.Genre + `</span>`
		}

		if movie.Streaming != "" {
			movieHTML += `<span class="bg-green-600 px-2 py-1 rounded">` + movie.Streaming + `</span>`
		}

		movieHTML += `
					</div>`

		if movie.IMDBLink != "" {
			movieHTML += `<div class="mt-2"><a href="` + movie.IMDBLink + `" target="_blank" class="text-blue-400 hover:text-blue-300 text-sm">View on IMDB</a></div>`
		}

		if movie.Notes != "" {
			movieHTML += `<p class="text-gray-400 mt-2 text-sm">` + movie.Notes + `</p>`
		}

		movieHTML += `
				</div>
				<div class="flex space-x-2">
					<button 
						data-movie-id="` + strconv.Itoa(movie.ID) + `"
						data-movie-title="` + movie.Title + `"
						data-movie-year="` + strconv.Itoa(movie.Year) + `"
						data-movie-genre="` + movie.Genre + `"
						data-movie-streaming="` + movie.Streaming + `"
						data-movie-notes="` + movie.Notes + `"
						data-movie-imdb="` + movie.IMDBLink + `"
						onclick="openEditModal(this)"
						class="text-blue-400 hover:text-blue-300 transition-colors duration-200">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
						</svg>
					</button>
					<button 
						hx-delete="/movie-board/delete/` + strconv.Itoa(movie.ID) + `"
						hx-target="closest div"
						hx-swap="outerHTML"
						hx-confirm="Are you sure you want to delete this movie?"
						class="text-red-400 hover:text-red-300 transition-colors duration-200">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
						</svg>
					</button>
				</div>
			</div>
		</div>`

		moviesHTML += movieHTML
	}

	// If no movies, show the "no movies" message
	if len(allMovies) == 0 {
		moviesHTML = `
		<div class="text-center py-8">
			<svg class="w-16 h-16 text-gray-600 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2m-9 0h10m-10 0a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V6a2 2 0 00-2-2"></path>
			</svg>
			<p class="text-gray-400">No movies added yet. Add your first movie above!</p>
		</div>`
	}

	// Create the response with out-of-band update for movie count
	response := moviesHTML + `
	<div id="movie-count" hx-swap-oob="true">
		` + strconv.Itoa(len(allMovies)) + ` movies in your list
	</div>`

	w.Write([]byte(response))
}
