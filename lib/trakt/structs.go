package trakt

type Ids struct {
	Trakt  int    `json:"trakt"`
	Tvdb   int    `json:"tvdb"`
	Imdb   string `json:"imdb"`
	Tmdb   int    `json:"tmdb"`
	Tvrage int    `json:"tvrage"`
}

type Show struct {
	Ids Ids
}

type ShowInfo struct {
	Show Show
}

type Episode struct {
	Season int    `json:"season"`
	Number int    `json:"number"`
	Title  string `json:"title"`
	Ids    Ids    `json:"ids"`
}

type Season struct {
	Number   int
	Episodes []Episode
}

type Movie struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	Ids   Ids    `json:"ids"`
}

type MovieSearchResult struct {
	Movie Movie
}

type ShowScrobbleBody struct {
	Episode  Episode `json:"episode"`
	Progress int     `json:"progress"`
}

type MovieScrobbleBody struct {
	Movie    Movie `json:"movie"`
	Progress int   `json:"progress"`
}
