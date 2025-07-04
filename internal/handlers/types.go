package handlers

// YearRange for form inputs
type YearRange struct {
	Min int
	Max int
}

// PageData represents the data passed to templates
type PageData struct {
	Title             string
	Count             int
	Genres            []string
	StreamingServices []string
	YearRange         YearRange
	Navigation        []NavItem
	FeatureCards      []FeatureCard
}

// MovieBoardData represents the data for the movie board page
type MovieBoardData struct {
	Title             string
	Movies            []Movie
	MovieCount        int
	Genres            []string
	StreamingServices []string
	YearRange         YearRange
	Navigation        []NavItem
	FormText          FormText
	Colors            ColorScheme
	BadgeColors       map[string]string
}

// TVShowBoardData represents the data for the TV show board page
type TVShowBoardData struct {
	Title             string
	TVShows           []TVShow
	TVShowCount       int
	Genres            []string
	StreamingServices []string
	YearRange         YearRange
	Navigation        []NavItem
	FormText          FormText
	Colors            ColorScheme
	BadgeColors       map[string]string
}

// NavItem represents a navigation item
type NavItem struct {
	URL      string
	Label    string
	Icon     string // Optional SVG icon class
	IsActive bool
}

// FeatureCard represents a feature card on the homepage
type FeatureCard struct {
	Title       string
	Description string
	URL         string
	Icon        string
	IconColor   string // Tailwind color classes
	ButtonColor string // Tailwind color classes
	ButtonText  string
}

// ColorScheme defines consistent colors for different UI elements
type ColorScheme struct {
	Primary   string // Main brand color
	Secondary string // Secondary brand color
	Success   string // Success states
	Warning   string // Warning states
	Error     string // Error states
	Info      string // Info states
	Neutral   string // Neutral/gray colors
}

// FormText defines consistent form labels and text
type FormText struct {
	AddNew       string
	Edit         string
	Cancel       string
	Save         string
	Delete       string
	Title        string
	Year         string
	Genre        string
	Streaming    string
	IMDBLink     string
	Notes        string
	ActiveSeason string
	AvailableNow string
}
