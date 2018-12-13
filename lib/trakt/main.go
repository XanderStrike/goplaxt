package trakt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/xanderstrike/goplaxt/lib/plex"
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

func Handle(pr plex.PlexResponse) {
	if pr.Metadata.LibrarySectionType == "show" {
		HandleShow(pr)
	} else if pr.Metadata.LibrarySectionType == "movie" {
		HandleMovie(pr)
	}
}

// [{"type":"show","score":1000,"show":{"title":"Disenchantment","year":2018,"ids":{"trakt":126558,"slug":"disenchantment","tvdb":340234,"imdb":"tt5363918","tmdb":73021,"tvrage":null}}}]

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

// [{"season":1,"number":1,"title":"A Princess, an Elf, and a Demon Walk Into a Bar","ids":{"trakt":2907713,"tvdb":6699345,"imdb":"tt6561570","tmdb":1494890,"tvrage":0}
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

func HandleShow(pr plex.PlexResponse) {
	fmt.Println("handling show")

	re := regexp.MustCompile("thetvdb://(\\d*)/(\\d*)/(\\d*)")
	showID := re.FindStringSubmatch(pr.Metadata.Guid)

	url := fmt.Sprintf("https://api.trakt.tv/search/tvdb/%s?type=show", showID[1])

	resp_body := makeRequest(url)
	fmt.Println(string(resp_body))

	var showInfo []ShowInfo
	_ = json.Unmarshal(resp_body, &showInfo)

	fmt.Println(showInfo[0].Show.Ids.Trakt)

	url = fmt.Sprintf("https://api.trakt.tv/shows/%d/seasons?extended=episodes", showInfo[0].Show.Ids.Trakt)

	resp_body = makeRequest(url)
	var seasons []Season
	_ = json.Unmarshal(resp_body, &seasons)

	fmt.Printf("%+v\n", seasons)

	sweetBytes, _ := json.Marshal(seasons[0].Episodes[0])
	fmt.Printf(string(sweetBytes))
}

func HandleMovie(pr plex.PlexResponse) {
	fmt.Println("handling movie")
}

func makeRequest(url string) []byte {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("trakt-api-version", "2")
	req.Header.Add("trakt-api-key", clientId)

	resp, err := client.Do(req)
	handleErr(err)

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)
	return resp_body
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
