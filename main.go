package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/xanderstrike/goplaxt/lib/plex"
	"github.com/xanderstrike/goplaxt/lib/store"
	"github.com/xanderstrike/goplaxt/lib/trakt"
)

type AuthorizePage struct {
	Authorized bool
	URL        string
}

func authorize(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	username := strings.ToLower(args["username"][0])
	log.Print(fmt.Sprintf("Handling auth request for %s", username))
	code := args["code"][0]
	result := trakt.AuthRequest(username, code, "", "authorization_code")

	user := store.NewUser(username, result["access_token"].(string), result["refresh_token"].(string))

	url := fmt.Sprintf("%s/api?id=%s", os.Getenv("REDIRECT_URI"), user.ID)

	log.Print(fmt.Sprintf("Authorized as %s", user.ID))

	tmpl := template.Must(template.ParseFiles("static/index.html"))
	data := AuthorizePage{
		Authorized: true,
		URL:        url,
	}
	tmpl.Execute(w, data)
}

func api(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	id := args["id"][0]
	log.Print(fmt.Sprintf("Webhook call for %s", id))

	user := store.GetUser(id)

	tokenAge := time.Since(user.Updated).Hours()
	if tokenAge > 1440 { // tokens expire after 3 months, so we refresh after 2
		log.Println("User access token outdated, refreshing...")
		result := trakt.AuthRequest(user.Username, "", user.RefreshToken, "refresh_token")
		user = store.UpdateUser(user, result["access_token"].(string), result["refresh_token"].(string))
		log.Println("Refreshed, continuing")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	re := plex.HandleWebhook(body)

	if strings.ToLower(re.Account.Title) == user.Username {
		trakt.Handle(re, user)
	} else {
		log.Println(fmt.Sprintf("Plex username %s does not equal %s, skipping", strings.ToLower(re.Account.Title), user.Username))
	}

	json.NewEncoder(w).Encode("success")
}

func main() {
	log.Print("Started!")
	router := mux.NewRouter()
	router.HandleFunc("/authorize", authorize).Methods("GET")
	router.HandleFunc("/api", api).Methods("POST")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("static/index.html"))
		data := AuthorizePage{
			Authorized: false,
			URL:        "https://plaxt.astandke.com/api?id=generate-your-own-silly",
		}
		tmpl.Execute(w, data)
	}).Methods("GET")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
