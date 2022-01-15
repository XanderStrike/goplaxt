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
	} else if os.Getenv(name+"_FILE") != "" {
		file, err := ioutil.ReadFile(os.Getenv(name + "_FILE"))

		if err != nil {
			panic(err)
		}

		return strings.TrimSpace(string(file))
	}

	return ""
}
