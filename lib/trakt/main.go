package trakt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const clientId string = "c9a8a36c476dcfe72b46b8be2237e8151486af90dac6b94548c89329f2a190c2"
const clientSecret string = "852aa926322f30d54d98d3693a95dfbf13efcaa7ce18f2fc1ad8b21a8463db51"

func AuthRequest(username, code, refreshToken, grantType string) map[string]interface{} {
	values := map[string]string{
		"code":          code,
		"refresh_token": refreshToken,
		"client_id":     clientId,
		"client_secret": clientSecret,
		"redirect_uri":  fmt.Sprintf("http://localhost:8000/authorize?username=%s", string(username)),
		"grant_type":    grantType,
	}
	jsonValue, _ := json.Marshal(values)

	resp, _ := http.Post("https://api.trakt.tv/oauth/token", "application/json", bytes.NewBuffer(jsonValue))

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}
