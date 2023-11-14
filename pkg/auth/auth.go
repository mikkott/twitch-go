package Auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Creds struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type TokenValidation struct {
	ClientID  string   `json:"client_id"`
	Login     string   `json:"login"`
	Scopes    []string `json:"scopes"`
	UserID    string   `json:"user_id"`
	ExpiresIn int      `json:"expires_in"`
}

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

type TokenValidationFailed struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (c *Config) SetToken() error {
	creds := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", c.ClientID, c.ClientSecret)
	headers := strings.NewReader(creds)

	resp, err := http.Post("https://id.twitch.tv/oauth2/token", "application/x-www-form-urlencoded", headers)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result Creds

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	c.Token = result.AccessToken

	return nil

}

func (c *Config) ValidateToken() error {
	header := fmt.Sprintf("Authorization: OAuth %s", c.Token)

	headers := strings.NewReader(header)

	resp, err := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", headers)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result TokenValidation
	err = json.Unmarshal(body, &result)
	if err != nil {
		var result TokenValidationFailed
		err = json.Unmarshal(body, &result)
		if err != nil {
			return err
		}
	}

	return nil

}
