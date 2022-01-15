package config

import (
	"io/ioutil"
	"os"
	"strings"
)

var TraktClientId string = getConfig("TRAKT_ID")
var TraktClientSecret string = getConfig("TRAKT_SECRET")

func getConfig(name string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}

	return ""
}
