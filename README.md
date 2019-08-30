**Notice** I will not be doing any more development on this project as Trakt now has [official support](https://blog.trakt.tv/plex-scrobbler-52db9b016ead) for Plex Webhooks. I'm using em, you should too. 

This tool should continue to work if you're not willing to pay for Trakt VIP, and I'll continue merging PRs, but don't expect your issues to be answered.

# Plaxt

[![CircleCI](https://circleci.com/gh/XanderStrike/goplaxt.svg?style=svg)](https://circleci.com/gh/XanderStrike/goplaxt) ![Docker Cloud Build](https://img.shields.io/docker/cloud/build/xanderstrike/goplaxt.svg)

Plex provides webhook integration for all Plex Pass subscribers, and users of their servers. A webhook is a request that the Plex application sends to third party services when a user takes an action, such as watching a movie or episode.

You can ask Plex to send these webhooks to this tool, which will then log those plays in your Trakt account.

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
Bare IP addresses and ports are totally fine, but keep in mind your Plaxt instance _must_ be accessible to _all_ the Plex servers you intend to play media from.

Once you have that, creating your container is a snap:

    docker create \
      --name=plaxt \
      --restart always \
      -v <path to configs>:/app/keystore \
      -e TRAKT_ID=<trakt_id> \
      -e TRAKT_SECRET=<trakt_secret> \
      -e ALLOWED_HOSTNAMES=<your public hostname(s) comma or space seperated> \
      -p 8000:8000 \
      xanderstrike/goplaxt:latest

If you are using a Raspberry Pi or other ARM based device, simply use
`xanderstrike/goplaxt:latest-arm7`.

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
