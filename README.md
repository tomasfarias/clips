# Clips

A Discord bot to look for Twitch clips!

[![Go Report Card](https://goreportcard.com/badge/github.com/tomasfarias/clips)](https://goreportcard.com/report/github.com/tomasfarias/clips)

## Requirements

The bot requires a Twitch Application Client ID and Client Secret, which can be found in the [Applications console](https://dev.twitch.tv/console/apps). From the Discord side, we require a Token for our Bot user which we can obtain by navigating into the bot users of our application in the [Developer Portal](https://discord.com/developers/applications). These credentials correspond to the following flags:
  * `-t`: Discord Bot token
  * `-c`: Twitch Client ID
  * `-s`: Twithc Client Secret

It is recommended to define the credentials in an `.env` instead of directly passing them as command line arguments.

## Running with Docker

Define your credentials in an `.env` file:

```
TOKEN=discord-token
CLIENT_ID=twitch-client-id
CLIENT_SECRET=twitch-client-secret
```

Build and run the container pointing to your `.env` file:

```
docker build --tag clips:1.0 .
docker run --env-file /path/to/.env --name clips clips:1.0
```

Once it's running, your Discord Bot user will need to be invited to your Discord server.
