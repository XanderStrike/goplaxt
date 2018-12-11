package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/peterbourgon/diskv"
)

const clientId string = "c9a8a36c476dcfe72b46b8be2237e8151486af90dac6b94548c89329f2a190c2"
const clientSecret string = "852aa926322f30d54d98d3693a95dfbf13efcaa7ce18f2fc1ad8b21a8463db51"

func flatTransform(s string) []string { return []string{} }

func write(key, value string) error {
	d := diskv.New(diskv.Options{
		BasePath:     "keystore",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	return d.Write(key, []byte(value))
}

func read(key string) (string, error) {
	d := diskv.New(diskv.Options{
		BasePath:     "keystore",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	value, err := d.Read(key)
	return string(value), err
}

func uuid() string {
	f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}

func newUser(username, access_token, refresh_token string) (string, error) {
	id := uuid()
	err := write(fmt.Sprintf("%s.username", id), username)
	if err != nil {
		return id, err
	}
	err = write(fmt.Sprintf("%s.access", id), access_token)
	if err != nil {
		return id, err
	}
	err = write(fmt.Sprintf("%s.refresh", id), refresh_token)
	if err != nil {
		return id, err
	}
	return id, nil
}

func getUser(id string) (string, string, string, error) {
	username, err := read(fmt.Sprintf("%s.username", id))
	if err != nil {
		return "", "", "", err
	}
	access_token, err := read(fmt.Sprintf("%s.access", id))
	if err != nil {
		return "", "", "", err
	}
	refresh_token, err := read(fmt.Sprintf("%s.refresh", id))
	if err != nil {
		return "", "", "", err
	}
	return username, access_token, refresh_token, nil
}

func authRequest(username, code, refreshToken, grantType string) map[string]interface{} {
	values := map[string]string{
		"code":          code,
		"refresh_token": refreshToken,
		"client_id":     clientId,
		"client_secret": clientSecret,
		"redirect_uri":  fmt.Sprintf("http://localhost:8000/authorize?username=%s", string(username)),
		"grant_type":    grantType,
	}
	jsonValue, _ := json.Marshal(values)

	// fmt.Println(string(jsonValue))

	resp, _ := http.Post("https://api.trakt.tv/oauth/token", "application/json", bytes.NewBuffer(jsonValue))

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func authorize(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	username := args["username"][0]
	code := args["code"][0]
	result := authRequest(username, code, "", "authorization_code")
	fmt.Println(result["access_token"])
	fmt.Println(result["refresh_token"])

	id, _ := newUser(username, result["access_token"].(string), result["refresh_token"].(string))

	url := fmt.Sprintf("http://localhost:8000/api?id=%s", id)
	fmt.Println(url)
	json.NewEncoder(w).Encode(url)
}

func api(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	id := args["id"][0]

	// TODO error handling
	username, access_token, refresh_token, _ := getUser(id)

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
