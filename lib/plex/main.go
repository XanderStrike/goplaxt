package plex

import (
	"encoding/json"
	"regexp"
)

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
