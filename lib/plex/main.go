package plex

import (
	"encoding/json"
	"regexp"
)

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

func HandleWebhook(body []byte) PlexResponse {
	re := regexp.MustCompile("({.*})")
	match := re.FindStringSubmatch(string(body))

	var pr PlexResponse
	err := json.Unmarshal([]byte(match[0]), &pr)
	if err != nil {
		panic(err)
	}
	return pr
}
