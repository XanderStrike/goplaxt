package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/xanderstrike/goplaxt/lib/store"
	"github.com/xanderstrike/goplaxt/lib/trakt"
	"github.com/xanderstrike/plexhooks"
)

var storage store.Store

type AuthorizePage struct {
	SelfRoot   string
	Authorized bool
	URL        string
	ClientID   string
}

func SelfRoot(r *http.Request) string {
	u, _ := url.Parse("")
	u.Host = r.Host
	u.Scheme = r.URL.Scheme
	u.Path = ""
	if u.Scheme == "" {
		u.Scheme = "http"

		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "https" {
			u.Scheme = "https"
		}
	}
	return u.String()
}

func authorize(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	username := strings.ToLower(args["username"][0])
	log.Print(fmt.Sprintf("Handling auth request for %s", username))
	code := args["code"][0]
	result := trakt.AuthRequest(SelfRoot(r), username, code, "", "authorization_code")

	user := store.NewUser(username, result["access_token"].(string), result["refresh_token"].(string), storage)

	url := fmt.Sprintf("%s/api?id=%s", SelfRoot(r), user.ID)

	log.Print(fmt.Sprintf("Authorized as %s", user.ID))

	tmpl := template.Must(template.ParseFiles("static/index.html"))
	data := AuthorizePage{
		SelfRoot:   SelfRoot(r),
		Authorized: true,
		URL:        url,
		ClientID:   os.Getenv("TRAKT_ID"),
	}
	tmpl.Execute(w, data)
}

func api(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	id := args["id"][0]
	log.Print(fmt.Sprintf("Webhook call for %s", id))

	user := storage.GetUser(id)

	tokenAge := time.Since(user.Updated).Hours()
	if tokenAge > 1440 { // tokens expire after 3 months, so we refresh after 2
		log.Println("User access token outdated, refreshing...")
		result := trakt.AuthRequest(SelfRoot(r), user.Username, "", user.RefreshToken, "refresh_token")
		user.UpdateUser(result["access_token"].(string), result["refresh_token"].(string))
		log.Println("Refreshed, continuing")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	regex := regexp.MustCompile("({.*})") // not the best way really
	match := regex.FindStringSubmatch(string(body))
	re, err := plexhooks.ParseWebhook([]byte(match[0]))
	if err != nil {
		panic(err)
	}

	// re := plexhooks.ParseWebhook([]byte(match[0]))

	if strings.ToLower(re.Account.Title) == user.Username {
		// FIXME - make everything take the pointer
		trakt.Handle(re, *user)
	} else {
		log.Println(fmt.Sprintf("Plex username %s does not equal %s, skipping", strings.ToLower(re.Account.Title), user.Username))
	}

	json.NewEncoder(w).Encode("success")
}

func allowedHostsHandler(allowedHostnames string) func(http.Handler) http.Handler {
	allowedHosts := strings.Split(regexp.MustCompile("https://|http://|\\s+").ReplaceAllString(strings.ToLower(allowedHostnames), ""), ",")
	log.Println("Allowed Hostnames:", allowedHosts)
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.EscapedPath() == "/healthcheck" {
				h.ServeHTTP(w, r)
				return
			}
			isAllowedHost := false
			lcHost := strings.ToLower(r.Host)
			for _, value := range allowedHosts {
				if lcHost == value {
					isAllowedHost = true
					break
				}
			}
			if !isAllowedHost {
				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Set("Content-Type", "text/plain")
				fmt.Fprintf(w, "Oh no!")
				return
			}
			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

type healthcheckStatus struct {
	Storage string
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	hasError := false
	status := &healthcheckStatus{}

	if err := storage.Ping(); err != nil {
		hasError = true
		status.Storage = err.Error()
	}

	w.Header().Set("Content-Type", "text/json")
	if hasError {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(status)
}

func main() {
	log.Print("Started!")
	if os.Getenv("REDIS_URI") != "" {
		storage = store.NewRedisStore(store.NewRedisClient(os.Getenv("REDIS_URI"), os.Getenv("REDIS_PASSWORD")))
		log.Println("Using redis storage:", os.Getenv("REDIS_URI"))
	} else {
		storage = store.NewDiskStore()
		log.Println("Using disk storage:")
	}

	router := mux.NewRouter()
	// Assumption: Behind a proper web server (nginx/traefik, etc) that removes/replaces trusted headers
	router.Use(handlers.ProxyHeaders)
	// which hostnames we are allowing
	// REDIRECT_URI = old legacy list
	// ALLOWED_HOSTNAMES = new accurate config variable
	// No env = all hostnames
	if os.Getenv("REDIRECT_URI") != "" {
		router.Use(allowedHostsHandler(os.Getenv("REDIRECT_URI")))
	} else if os.Getenv("ALLOWED_HOSTNAMES") != "" {
		router.Use(allowedHostsHandler(os.Getenv("ALLOWED_HOSTNAMES")))
	}
	router.HandleFunc("/authorize", authorize).Methods("GET")
	router.HandleFunc("/api", api).Methods("POST")
	router.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("static/index.html"))
		data := AuthorizePage{
			SelfRoot:   SelfRoot(r),
			Authorized: false,
			URL:        "https://plaxt.astandke.com/api?id=generate-your-own-silly",
			ClientID:   os.Getenv("TRAKT_ID"),
		}
		tmpl.Execute(w, data)
	}).Methods("GET")
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = "0.0.0.0:8000"
	}
	log.Print("Started on " + listen + "!")
	log.Fatal(http.ListenAndServe(listen, router))
}
