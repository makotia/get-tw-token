package main

import (
	"fmt"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"github.com/fatih/color"
	"github.com/skratchdot/open-golang/open"
)

var (
	config         oauth1.Config
	consumerKey    string
	consumerSecret string
)

func main() {
	cyan := color.New(color.FgCyan)
	magenta := color.New(color.FgHiMagenta)
	red := color.New(color.FgRed)
	cyan.Printf("ConsumerKey:    ")
	fmt.Scan(&consumerKey)
	cyan.Printf("ConsumerSecret: ")
	fmt.Scan(&consumerSecret)

	config = oauth1.Config{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		CallbackURL:    "oob",
		Endpoint:       twitter.AuthorizeEndpoint,
	}

	requestToken, err := login()
	if err != nil {
		red.Println("Request Token Phase: %s", err.Error())
	}
	accessToken, err := receivePIN(requestToken)
	if err != nil {
		red.Println("Access Token Phase: %s", err.Error())
	}

	magenta.Printf("AccessToken:  ")
	fmt.Println(accessToken.Token)
	magenta.Printf("AccessSecret: ")
	fmt.Println(accessToken.TokenSecret)
}

func login() (requestToken string, err error) {
	requestToken, _, err = config.RequestToken()
	if err != nil {
		return "", err
	}
	authorizationURL, err := config.AuthorizationURL(requestToken)
	if err != nil {
		return "", err
	}
	if err := open.Run(authorizationURL.String()); err != nil {
		fmt.Printf("Open:\n%s\n", authorizationURL.String())
	}
	return requestToken, err
}

func receivePIN(requestToken string) (*oauth1.Token, error) {
	cyan := color.New(color.FgCyan)
	cyan.Printf("PIN: ")
	var verifier string
	_, err := fmt.Scanf("%s", &verifier)
	if err != nil {
		return nil, err
	}
	accessToken, accessSecret, err := config.AccessToken(requestToken, "", verifier)
	if err != nil {
		return nil, err
	}
	return oauth1.NewToken(accessToken, accessSecret), err
}
