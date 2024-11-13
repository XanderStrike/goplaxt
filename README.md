**Notice** I will not be doing any more development on this project as Trakt now has [official support](https://blog.trakt.tv/plex-scrobbler-52db9b016ead) for Plex Webhooks. I'm using em, you should too.

# Plaxt

[![CircleCI](https://circleci.com/gh/XanderStrike/goplaxt.svg?style=svg)](https://circleci.com/gh/XanderStrike/goplaxt) ![Docker Cloud Build](https://img.shields.io/docker/cloud/build/xanderstrike/goplaxt.svg)

Plex provides webhook integration for all Plex Pass subscribers, and users of their servers. A webhook is a request that the Plex application sends to third party services when a user takes an action, such as watching a movie or episode.

You can ask Plex to send these webhooks to this tool, which will then log those plays in your Trakt account.

### Deploying For Yourself

Goplaxt is designed to be run in Docker. You can host it right on your Plex server!

To run it yourself, first create an API application through Trakt [here](https://trakt.tv/oauth/applications). Set the
Allowed Hostnames to be the URI you will hit to access Plaxt, plus `/authorize`. So if you're exposing your server at
`http://10.20.30.40:8000`, you'll set it to `http://10.20.30.40:8000/authorize`. Bare IP addresses and ports are
totally fine, but keep in mind your Plaxt instance _must_ be accessible to _all_ the Plex servers you intend to 
play media from.

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

Alternatively you can use `docker-compose`:

```yaml
version: "3.4" # This will probably also work with version 2
services:
  plaxt:
    container_name: plaxt
    environment:
    - TRAKT_ID=<trakt_id>
    - TRAKT_SECRET=<trakt_secret>
    - ALLOWED_HOSTNAMES=<your public hostname(s) comma or space seperated>
    image: xanderstrike/goplaxt
    ports:
    - 8000:8000
    restart: unless-stopped
    volumes:
    - <path to configs>:/app/keystore
```

### Contributing

This repository is now a public archive but I encourage forking!

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
