package handlers

// Genres list for templating (fallback defaults)
var Genres = []string{
	"Action",
	"Adventure",
	"Animation",
	"Children's",
	"Comedy",
	"Crime",
	"Documentary",
	"Drama",
	"Fantasy",
	"Holiday",
	"Horror",
	"Mystery",
	"Romance",
	"Sci-Fi",
	"Thriller",
	"Western",
}

// StreamingServices list for templating (fallback defaults)
var StreamingServices = []string{
	"Amazon Prime",
	"Apple TV+",
	"Crunchyroll",
	"Disney+",
	"HBO Max",
	"Hulu",
	"Netflix",
	"Other",
	"Paramount+",
	"Peacock",
}

// Navigation list for templating
var Navigation = []NavItem{
	{URL: "/", Label: "Home", Icon: "M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"},
	{URL: "/movie-board", Label: "Movie Board", Icon: "M7 4V2a1 1 0 011-1h4a1 1 0 011 1v2h4a1 1 0 011 1v14a1 1 0 01-1 1H3a1 1 0 01-1-1V5a1 1 0 011-1h4zM9 4V3h6v1H9z"},
	{URL: "/tv-shows-board", Label: "TV Shows Board", Icon: "M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"},
	{URL: "/ai", Label: "AI", Icon: "M8.625 12a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H8.25m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H12m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0h-.375M21 12c0 4.556-4.03 8.25-9 8.25a9.764 9.764 0 0 1-2.555-.337A5.972 5.972 0 0 1 5.41 20.97a5.969 5.969 0 0 1-.474-.065 4.48 4.48 0 0 0 .978-2.025c.09-.457-.133-.901-.467-1.226C3.93 16.178 3 14.189 3 12c0-4.556 4.03-8.25 9-8.25s9 3.694 9 8.25Z"},
}

// FeatureCards list for the homepage
var FeatureCards = []FeatureCard{
	{
		Title:       "Movie Board",
		Description: "Manage your movie watchlist, add new films, and get random movie suggestions",
		URL:         "/movie-board",
		Icon:        "M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z",
		IconColor:   "blue",
		ButtonColor: "blue",
		ButtonText:  "Open Movie Board",
	},
	{
		Title:       "TV Shows Board",
		Description: "Track your favorite TV shows, manage seasons, and organize your binge-watching",
		URL:         "/tv-shows-board",
		Icon:        "M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z",
		IconColor:   "purple",
		ButtonColor: "purple",
		ButtonText:  "Open TV Shows Board",
	},
	{
		Title:       "Home Assistant",
		Description: "Control your smart devices, monitor your temps, automate your life",
		URL:         "http://has.gotpwnd.org:8123/",
		Icon:        "M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6",
		IconColor:   "sky",
		ButtonColor: "sky",
		ButtonText:  "Open Home Assistant",
	},
	{
		Title:       "AI",
		Description: "Homenet custom AI, running locally. No data is uploaded to anywhere",
		URL:         "/ai",
		Icon:        "M8.625 12a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H8.25m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H12m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0h-.375M21 12c0 4.556-4.03 8.25-9 8.25a9.764 9.764 0 0 1-2.555-.337A5.972 5.972 0 0 1 5.41 20.97a5.969 5.969 0 0 1-.474-.065 4.48 4.48 0 0 0 .978-2.025c.09-.457-.133-.901-.467-1.226C3.93 16.178 3 14.189 3 12c0-4.556 4.03-8.25 9-8.25s9 3.694 9 8.25Z",
		IconColor:   "pink",
		ButtonColor: "pink",
		ButtonText:  "Open Homenet AI",
	},
}

// AppColors defines the application's color scheme (fallback defaults)
var AppColors = ColorScheme{
	Primary:   "blue",
	Secondary: "purple",
	Success:   "green",
	Warning:   "yellow",
	Error:     "red",
	Info:      "sky",
	Neutral:   "gray",
}

// BadgeColors maps data types to their display colors (fallback defaults)
var BadgeColors = map[string]string{
	"year":      "gray",
	"genre":     "blue",
	"streaming": "green",
	"active":    "yellow",
}

// MovieFormText defines text for movie forms
var MovieFormText = FormText{
	AddNew:       "Add New Movie",
	Edit:         "Edit Movie",
	Cancel:       "Cancel",
	Save:         "Add Movie",
	Delete:       "Delete",
	Title:        "Movie Title",
	Year:         "Year",
	Genre:        "Genre",
	Streaming:    "Streaming Service",
	IMDBLink:     "IMDB Link (Optional)",
	Notes:        "Notes (Optional)",
	AvailableNow: "Available Now",
}

// TVShowFormText defines text for TV show forms
var TVShowFormText = FormText{
	AddNew:       "Add New TV Show",
	Edit:         "Edit TV Show",
	Cancel:       "Cancel",
	Save:         "Add TV Show",
	Delete:       "Delete",
	Title:        "TV Show Title",
	Year:         "Year",
	Genre:        "Genre",
	Streaming:    "Streaming Service",
	IMDBLink:     "IMDB Link (Optional)",
	Notes:        "Notes (Optional)",
	ActiveSeason: "Active Season",
}
