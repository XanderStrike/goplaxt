# Plaxt

This is a full rewrite of my somewhat popular previous iteration. This time it's written in Go
and deployable with Docker so I can run it on my own infrastructure instead of Heroku.

To start scrobbling today, head to [plaxt.astandke.com](https://plaxt.astandke.com) and enter your Plex username!
It's as easy as can be!

If you experience any problems or have any suggestions, please don't hesitate to create an issue on this repo.

### Deploying For Yourself

Goplaxt is designed to be run in Docker. You can host it right on your Plex server!

To run it yourself, first create an API application through Trakt [here](https://trakt.tv/oauth/applications). Set the
Redirect URI to be the URI you will hit to access Plaxt, plus `/authorize`. So if you're exposing your server at
`http://10.20.30.40:8000`, you'll set it to `http://10.20.30.40:8000/authorize`. For the CORS Origin, just use the URI.
Bare IP addresses and ports are totally fine, but keep in mind your Plaxt instance _must_ be available to the public 
internet in order for Plex to send it data.

Once you have that, creating your container is a snap:

    docker create \
      --name=plaxt \
      --restart always \
      -v <path to configs>:/app/keystore \
      -e TRAKT_ID=<trakt_id> \
      -e TRAKT_SECRET=<trakt_secret> \
      -e ALLOWED_HOSTNAME=<your public hostname> \
      -p 8000:8000 \
      xanderstrike/goplaxt:latest

Then go ahead and start it with:

    docker start plaxt

### Contributing

Please do! I accept any and all PRs. My golang is not the best currently, so I'd love some thoughts on worthwhile
refactors. I sort of blew through this without adding any tests, so testing won't be a hard requirement for
contributions until I add some (though they're always welcome, of course).

### Security PSA

You should know that by using the instance I host, I perminantly retain your Plex username, and an API key that
allows me to add plays to your Trakt account (but not your username). Also, I log the title and year of films
you watch and the title, season, and episode of shows you watch. These logs are temporary and are rotated every
24 hours with older logs perminantly deleted.

I promise to Do No Harm with this information. It will never leave my server, and I won't look at it unless I'm
troubleshooting bugs. Frankly, I couldn't care less. However, I believe it's important to disclose my access to
your information. If you are not comfortable sharing I encourage you to host the application on your own hardware.

[I have never been served with any government requests for data](https://en.wikipedia.org/wiki/Warrant_canary).

### License

MIT
