package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xanderstrike/goplaxt/lib/plex"
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

func api(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	id := args["id"][0]

	username, accessToken, refreshToken, _ := store.GetUser(id)

	fmt.Println(fmt.Sprintf("%s: %s %s", username, accessToken, refreshToken))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	re := plex.HandleWebhook(body)

	trakt.Handle(re, accessToken)

	json.NewEncoder(w).Encode("success")
}

func main() {
	fmt.Println("Here we go!")
	router := mux.NewRouter()
	router.HandleFunc("/authorize", authorize).Methods("GET")
	router.HandleFunc("/api", api).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Fatal(http.ListenAndServe(":8000", router))
}
