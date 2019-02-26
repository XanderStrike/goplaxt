package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/xanderstrike/plexhooks"
)

func hook(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile("({.*})") // not the best way really
	match := re.FindStringSubmatch(string(body))

	response, err := plexhooks.ParseWebhook([]byte(match[0]))
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("%s played an item called %s!", response.Account.Title, response.Metadata.Title))

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/hook", hook).Methods("POST")

	log.Println("Now serving on 0.0.0.0:8000/hook")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
