package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/xanderstrike/goplaxt/lib/store"
	"github.com/xanderstrike/goplaxt/lib/trakt"
)

const clientId string = "c9a8a36c476dcfe72b46b8be2237e8151486af90dac6b94548c89329f2a190c2"
const clientSecret string = "852aa926322f30d54d98d3693a95dfbf13efcaa7ce18f2fc1ad8b21a8463db51"

func authorize(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	username := args["username"][0]
	code := args["code"][0]
	result := trakt.AuthRequest(username, code, "", "authorization_code")

	id, _ := store.NewUser(username, result["access_token"].(string), result["refresh_token"].(string))

	url := fmt.Sprintf("http://localhost:8000/api?id=%s", id)
	json.NewEncoder(w).Encode(url)
}

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

func api(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	id := args["id"][0]

	username, access_token, refresh_token, _ := store.GetUser(id)

	fmt.Println(fmt.Sprintf("%s: %s %s", username, access_token, refresh_token))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile("({.*})")
	match := re.FindStringSubmatch(string(body))

	var pr PlexResponse
	err = json.Unmarshal([]byte(match[0]), &pr)
	if err != nil {
		panic(err)
	}
	log.Println(pr)

	json.NewEncoder(w).Encode("")
}

func main() {
	fmt.Println("Here we go!")
	router := mux.NewRouter()
	router.HandleFunc("/authorize", authorize).Methods("GET")
	router.HandleFunc("/api", api).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Fatal(http.ListenAndServe(":8000", router))
}
