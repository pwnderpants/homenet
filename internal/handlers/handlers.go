package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/pwnderpants/homenet/internal/database"
)

// Use the Movie struct from database package
type Movie = database.Movie
type TVShow = database.TVShow

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

// TVShowBoardData represents the data for the TV show board page
type TVShowBoardData struct {
	Title       string
	TVShows     []TVShow
	TVShowCount int
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

// TVShowBoardHandler handles the TV show board page request
func TVShowBoardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/tv-shows-board.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// Get TV shows from database
	tvShows, err := database.GetAllTVShows()
	if err != nil {
		http.Error(w, "Failed to load TV shows: "+err.Error(), http.StatusInternalServerError)

		return
	}

	data := TVShowBoardData{
		Title:       "TV Shows Board",
		TVShows:     tvShows,
		TVShowCount: len(tvShows),
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

// AddTVShowHandler handles adding a new TV show
func AddTVShowHandler(w http.ResponseWriter, r *http.Request) {
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
	activeSeason := r.FormValue("active_season") == "on"

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

	// Create new TV show
	newTVShow := TVShow{
		Title:        title,
		Year:         year,
		Genre:        genre,
		Streaming:    streaming,
		Notes:        notes,
		IMDBLink:     imdbLink,
		ActiveSeason: activeSeason,
	}

	// Add to database
	tvShowID, err := database.AddTVShow(newTVShow)

	if err != nil {
		http.Error(w, "Failed to add TV show: "+err.Error(), http.StatusInternalServerError)

		return
	}

	// Set the ID for HTML generation
	newTVShow.ID = tvShowID

	// Get all TV shows from database to return the complete updated list
	allTVShows, err := database.GetAllTVShows()

	if err != nil {
		http.Error(w, "Failed to get TV shows: "+err.Error(), http.StatusInternalServerError)

		return
	}

	// Return the complete TV show list HTML for HTMX to replace
	w.Header().Set("Content-Type", "text/html")

	// Generate HTML for all TV shows
	var tvShowsHTML string

	for _, tvShow := range allTVShows {
		tvShowHTML := `
		<div class="bg-gray-700 rounded-lg p-4 border border-gray-600">
			<div class="flex justify-between items-start">
				<div class="flex-1">
					<h4 class="text-lg font-semibold text-white">` + tvShow.Title + `</h4>
					<div class="flex items-center space-x-4 mt-2 text-sm text-gray-300">`

		if tvShow.Year > 0 {
			tvShowHTML += `<span class="bg-gray-600 px-2 py-1 rounded">` + strconv.Itoa(tvShow.Year) + `</span>`
		}

		if tvShow.Genre != "" {
			tvShowHTML += `<span class="bg-blue-600 px-2 py-1 rounded">` + tvShow.Genre + `</span>`
		}

		if tvShow.Streaming != "" {
			tvShowHTML += `<span class="bg-green-600 px-2 py-1 rounded">` + tvShow.Streaming + `</span>`
		}

		tvShowHTML += `
					</div>`

		if tvShow.IMDBLink != "" {
			tvShowHTML += `<div class="mt-2"><a href="` + tvShow.IMDBLink + `" target="_blank" class="text-blue-400 hover:text-blue-300 text-sm">View on IMDB</a></div>`
		}

		if tvShow.Notes != "" {
			tvShowHTML += `<p class="text-gray-400 mt-2 text-sm">` + tvShow.Notes + `</p>`
		}

		tvShowHTML += `
				</div>
				<div class="flex space-x-2">
					<button 
						data-tvshow-id="` + strconv.Itoa(tvShow.ID) + `"
						data-tvshow-title="` + tvShow.Title + `"
						data-tvshow-year="` + strconv.Itoa(tvShow.Year) + `"
						data-tvshow-genre="` + tvShow.Genre + `"
						data-tvshow-streaming="` + tvShow.Streaming + `"
						data-tvshow-notes="` + tvShow.Notes + `"
						data-tvshow-imdb="` + tvShow.IMDBLink + `"
						onclick="openEditModal(this)"
						class="text-blue-400 hover:text-blue-300 transition-colors duration-200">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
						</svg>
					</button>
					<button 
						hx-delete="/tv-shows-board/delete/` + strconv.Itoa(tvShow.ID) + `"
						hx-target="closest div"
						hx-swap="outerHTML"
						hx-confirm="Are you sure you want to delete this TV show?"
						class="text-red-400 hover:text-red-300 transition-colors duration-200">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
						</svg>
					</button>
				</div>
			</div>
		</div>`

		tvShowsHTML += tvShowHTML
	}

	// If no TV shows, show the "no TV shows" message
	if len(allTVShows) == 0 {
		tvShowsHTML = `
		<div class="text-center py-8">
			<svg class="w-16 h-16 text-gray-600 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2m-9 0h10m-10 0a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V6a2 2 0 00-2-2"></path>
			</svg>
			<p class="text-gray-400">No TV shows added yet. Add your first TV show above!</p>
		</div>`
	}

	// Create the response with out-of-band update for TV show count
	response := tvShowsHTML + `
	<div id="tvshow-count" hx-swap-oob="true">
		` + strconv.Itoa(len(allTVShows)) + ` TV shows in your list
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

// DeleteTVShowHandler handles deleting a TV show
func DeleteTVShowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	// Extract TV show ID from URL
	path := r.URL.Path
	idStr := path[len("/tv-shows-board/delete/"):]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid TV show ID", http.StatusBadRequest)
		return
	}

	// Delete from database
	err = database.DeleteTVShow(id)

	if err != nil {
		http.Error(w, "Failed to delete TV show: "+err.Error(), http.StatusInternalServerError)
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

// EditTVShowHandler handles editing an existing TV show
func EditTVShowHandler(w http.ResponseWriter, r *http.Request) {
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
	activeSeason := r.FormValue("active_season") == "on"

	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)

		return
	}

	// Parse TV show ID
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid TV show ID", http.StatusBadRequest)

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

	// Update TV show in database
	updatedTVShow := TVShow{
		ID:           id,
		Title:        title,
		Year:         year,
		Genre:        genre,
		Streaming:    streaming,
		Notes:        notes,
		IMDBLink:     imdbLink,
		ActiveSeason: activeSeason,
	}

	err = database.UpdateTVShow(updatedTVShow)

	if err != nil {
		http.Error(w, "Failed to update TV show: "+err.Error(), http.StatusInternalServerError)

		return
	}

	// Get all TV shows from database to return the complete updated list
	allTVShows, err := database.GetAllTVShows()

	if err != nil {
		http.Error(w, "Failed to get TV shows: "+err.Error(), http.StatusInternalServerError)

		return
	}

	// Return the complete TV show list HTML for HTMX to replace
	w.Header().Set("Content-Type", "text/html")

	// Generate HTML for all TV shows
	var tvShowsHTML string

	for _, tvShow := range allTVShows {
		tvShowHTML := `
		<div class="bg-gray-700 rounded-lg p-4 border border-gray-600">
			<div class="flex justify-between items-start">
				<div class="flex-1">
					<h4 class="text-lg font-semibold text-white">` + tvShow.Title + `</h4>
					<div class="flex items-center space-x-4 mt-2 text-sm text-gray-300">`

		if tvShow.Year > 0 {
			tvShowHTML += `<span class="bg-gray-600 px-2 py-1 rounded">` + strconv.Itoa(tvShow.Year) + `</span>`
		}

		if tvShow.Genre != "" {
			tvShowHTML += `<span class="bg-blue-600 px-2 py-1 rounded">` + tvShow.Genre + `</span>`
		}

		if tvShow.Streaming != "" {
			tvShowHTML += `<span class="bg-green-600 px-2 py-1 rounded">` + tvShow.Streaming + `</span>`
		}

		tvShowHTML += `
					</div>`

		if tvShow.IMDBLink != "" {
			tvShowHTML += `<div class="mt-2"><a href="` + tvShow.IMDBLink + `" target="_blank" class="text-blue-400 hover:text-blue-300 text-sm">View on IMDB</a></div>`
		}

		if tvShow.Notes != "" {
			tvShowHTML += `<p class="text-gray-400 mt-2 text-sm">` + tvShow.Notes + `</p>`
		}

		tvShowHTML += `
				</div>
				<div class="flex space-x-2">
					<button 
						data-tvshow-id="` + strconv.Itoa(tvShow.ID) + `"
						data-tvshow-title="` + tvShow.Title + `"
						data-tvshow-year="` + strconv.Itoa(tvShow.Year) + `"
						data-tvshow-genre="` + tvShow.Genre + `"
						data-tvshow-streaming="` + tvShow.Streaming + `"
						data-tvshow-notes="` + tvShow.Notes + `"
						data-tvshow-imdb="` + tvShow.IMDBLink + `"
						onclick="openEditModal(this)"
						class="text-blue-400 hover:text-blue-300 transition-colors duration-200">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
						</svg>
					</button>
					<button 
						hx-delete="/tv-shows-board/delete/` + strconv.Itoa(tvShow.ID) + `"
						hx-target="closest div"
						hx-swap="outerHTML"
						hx-confirm="Are you sure you want to delete this TV show?"
						class="text-red-400 hover:text-red-300 transition-colors duration-200">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
						</svg>
					</button>
				</div>
			</div>
		</div>`

		tvShowsHTML += tvShowHTML
	}

	// If no TV shows, show the "no TV shows" message
	if len(allTVShows) == 0 {
		tvShowsHTML = `
		<div class="text-center py-8">
			<svg class="w-16 h-16 text-gray-600 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2m-9 0h10m-10 0a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V6a2 2 0 00-2-2"></path>
			</svg>
			<p class="text-gray-400">No TV shows added yet. Add your first TV show above!</p>
		</div>`
	}

	// Create the response with out-of-band update for TV show count
	response := tvShowsHTML + `
	<div id="tvshow-count" hx-swap-oob="true">
		` + strconv.Itoa(len(allTVShows)) + ` TV shows in your list
	</div>`

	w.Write([]byte(response))
}

// RandomMovieHandler handles getting a random movie
func RandomMovieHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	// Get random movie from database
	randomMovie, err := database.GetRandomMovie()

	if err != nil {
		http.Error(w, "Failed to get random movie: "+err.Error(), http.StatusInternalServerError)

		return
	}

	if randomMovie == nil {
		// No movies in database
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<div class="text-center py-8">
				<svg class="w-16 h-16 text-gray-600 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2m-9 0h10m-10 0a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V6a2 2 0 00-2-2"></path>
				</svg>
				<p class="text-gray-400">No movies in your list yet. Add some movies first!</p>
			</div>
		`))
		return
	}

	// Return the random movie HTML for the modal
	w.Header().Set("Content-Type", "text/html")

	movieHTML := `
		<div class="bg-gray-700 rounded-lg p-6 border border-gray-600">
			<div class="text-center mb-4">
				<h3 class="text-2xl font-bold text-white mb-2">🎬 Your Random Movie Pick!</h3>
				<p class="text-gray-300">Here's what you should watch tonight:</p>
			</div>
			<div class="space-y-4">
				<div class="text-center">
					<h4 class="text-xl font-semibold text-white mb-2">` + randomMovie.Title + `</h4>
					<div class="flex items-center justify-center space-x-4 text-sm text-gray-300">`

	if randomMovie.Year > 0 {
		movieHTML += `<span class="bg-gray-600 px-2 py-1 rounded">` + strconv.Itoa(randomMovie.Year) + `</span>`
	}

	if randomMovie.Genre != "" {
		movieHTML += `<span class="bg-blue-600 px-2 py-1 rounded">` + randomMovie.Genre + `</span>`
	}

	if randomMovie.Streaming != "" {
		movieHTML += `<span class="bg-green-600 px-2 py-1 rounded">` + randomMovie.Streaming + `</span>`
	}

	movieHTML += `
					</div>
				</div>`

	if randomMovie.IMDBLink != "" {
		movieHTML += `
				<div class="text-center">
					<a href="` + randomMovie.IMDBLink + `" target="_blank" class="text-blue-400 hover:text-blue-300 text-sm">View on IMDB</a>
				</div>`
	}

	if randomMovie.Notes != "" {
		movieHTML += `
				<div class="text-center">
					<p class="text-gray-400 text-sm italic">"` + randomMovie.Notes + `"</p>
				</div>`
	}

	movieHTML += `
			</div>
		</div>`

	w.Write([]byte(movieHTML))
}
