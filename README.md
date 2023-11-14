# twitch-go
Twitch chatbot written in go.

## How to to get started
Sign up, register application and generate client secret at: https://dev.twitch.tv/console.

Name: anything, URL: http://localhost:3000, Category: Chat bot.

Install twitch-cli: `go install github.com/twitchdev/twitch-cli@latest`

Get token with scopes: `twitch-cli token -u -s "chat:read chat:edit"`

### Configuration file .env

Create configuration file `.env` at the same directory as Makefile.

Populate with required key-values:
```
DEBUG=true
TWITCH_USERNAME=<yourtwitchusername>
TWITCH_CLIENT_ID=<from dev console>
TWITCH_SECRET=<from dev console>
TWITCH_TOKEN=<from twitch-cli>
TWITCH_CAPABILITIES=twitch.tv/membership,twitch.tv/tags,twitch.tv/commands
TWITCH_CHANNELS=<whatever channels separated by comma>
```