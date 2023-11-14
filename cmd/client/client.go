package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type Config struct {
	ServerAddr   string `default:"wss://irc-ws.chat.twitch.tv:443"`
	Debug        bool
	Username     string
	ClientID     string
	ClientSecret string
	Token        string
	Channels     []string
	Capabilities []string
	CapReq       string
}

func capReq(caps []string) string {
	capReq := ""
	for _, c := range caps {
		capReq = fmt.Sprintf("%s %s", c, capReq)
	}

	capReq = strings.TrimSpace(fmt.Sprintf("CAP REQ :%s", capReq))

	return capReq
}

func ping(ws *websocket.Conn) error {
	for {
		err := writeWs(ws, "PING")
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 60)
	}
}

func join(ws *websocket.Conn, channels []string) {
	for _, channel := range channels {
		msg := fmt.Sprintf("JOIN #%s", channel)
		err := writeWs(ws, msg)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 1)
	}
}

func readWs(ws *websocket.Conn) (string, error) {
	var err error
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	message := string(msg[:n])

	if err != nil {
		return "", err
	}

	return message, nil
}

func writeWs(ws *websocket.Conn, msg string) error {
	_, err := ws.Write([]byte(fmt.Sprintln(msg)))

	if err != nil {
		return err
	}

	return nil
}

func auth(ws *websocket.Conn, c Config) {
	var err error
	pass := fmt.Sprintf("PASS oauth:%s", c.Token)
	nick := fmt.Sprintf("NICK %s", c.Username)

	authSeq := []string{
		c.CapReq,
		pass,
		nick,
	}
	for _, msg := range authSeq {
		err = writeWs(ws, msg)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func (c *Config) initConfig() {
	debug := strings.ToLower(os.Getenv("DEBUG"))
	if debug == "true" {
		c.Debug = true
	} else {
		c.Debug = false
	}
	c.ServerAddr = "wss://irc-ws.chat.twitch.tv:443"
	c.Username = strings.ToLower(os.Getenv("TWITCH_USERNAME"))
	c.ClientID = os.Getenv("TWITCH_CLIENT_ID")
	c.ClientSecret = os.Getenv("TWITCH_SECRET")
	c.Capabilities = strings.Split(os.Getenv("TWITCH_CAPABILITIES"), ",")
	c.CapReq = capReq(c.Capabilities)
	c.Channels = strings.Split(os.Getenv("TWITCH_CHANNELS"), ",")
	c.Token = os.Getenv("TWITCH_TOKEN")

}

type Message struct {
	Channel   string
	Message   string
	Command   string
	Nick      string
	Timestamp int
}

func parseMsg(msg string) *Message {
	re, _ := regexp.Compile(`(?P<Prefix1>;tmi-sent-ts=)(?P<Timestamp>[0-9]+).*(?P<Prefix>\:.*\!.*\@)(?P<Nick>.*)(?P<Server>.tmi.twitch.tv )(?P<Command>.*)(?P<Extra> #)(?P<Channel>[a-z09]+)`)

	if !re.MatchString(msg) {
		return nil
	}

	matches := re.FindStringSubmatch(msg)

	nickIndex := re.SubexpIndex("Nick")
	timestampIndex := re.SubexpIndex("Timestamp")

	commandIndex := re.SubexpIndex("Command")
	channelIndex := re.SubexpIndex("Channel")

	var message Message
	message.Nick = matches[nickIndex]
	message.Command = matches[commandIndex]
	message.Channel = matches[channelIndex]
	message.Timestamp, _ = strconv.Atoi(matches[timestampIndex])

	return &message

}

func main() {
	var c Config

	c.initConfig()
	origin := "https://www.twitch.tv"
	url := c.ServerAddr

	fmt.Println(origin, url)
	ws, err := websocket.Dial(url, "", origin)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(c)
	auth(ws, c)
	join(ws, c.Channels)

	go ping(ws)

	for {
		msg, err := readWs(ws)

		if err != nil {
			log.Fatal(err)
		}

		message := parseMsg(msg)

		if message != nil {
			fmt.Println(message)
		}

		if c.Debug {
			fmt.Println(msg)
		}

		time.Sleep(time.Millisecond * 10)
	}
}
