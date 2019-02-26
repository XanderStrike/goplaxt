package plexhooks

import "encoding/json"

func ParseWebhook(request []byte) (PlexResponse, error) {
	var pr PlexResponse
	err := json.Unmarshal(request, &pr)
	return pr, err
}
