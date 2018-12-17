package plex

type Account struct {
	Title string
}

type Metadata struct {
	LibrarySectionType string
	Title              string
	Year               int
	Guid               string
}

type PlexResponse struct {
	Event    string
	Account  Account
	Metadata Metadata
}
