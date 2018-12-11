package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xanderstrike/plaxt/lib/store"
	"github.com/xanderstrike/plaxt/lib/trakt"
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

	username, access_token, refresh_token, _ := store.GetUser(id)

	fmt.Println(fmt.Sprintf("%s: %s %s", username, access_token, refresh_token))

	json.NewEncoder(w).Encode("")
}

func main() {
	fmt.Println("Here we go!")
	router := mux.NewRouter()
	router.HandleFunc("/authorize", authorize).Methods("GET")
	router.HandleFunc("/api", api).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Fatal(http.ListenAndServe(":8000", router))
}
