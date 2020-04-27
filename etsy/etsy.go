package etsy

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dghubble/oauth1"
)

type config struct {
	APIKey       string `json:"API_KEY"`
	SharedSecret string `json:"SHARED_SECRET"`
}

//Client type for Etsy API calls
type Client struct {
	Config  config
	BaseURL string
	Client  *http.Client
}

//Parameters to pass to Etsy api
type Parameters struct {
	Quantity int64
}

//NewClient creates a new Etsy client to make request to the etsy api
func NewClient() Client {
	client := Client{
		Config:  readConfig(),
		BaseURL: "https://openapi.etsy.com/v2/listings/active",
		Client:  &http.Client{},
	}
	return client
}

func readConfig() config {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln("Error Opening Config File: ", err)
	}
	configData := config{}

	err = json.Unmarshal([]byte(file), &configData)
	if err != nil {
		log.Fatalln("Error Parsing Config File: ", err)
	}
	return configData
}

//Authenticate performs an Oauth1.0 authentication with the Etsy api.
func (c Client) Authenticate() *http.Client {
	endpoint := oauth1.Endpoint{
		RequestTokenURL: "https://openapi.etsy.com/v2/oauth/request_token?scope=listings_w",
		AuthorizeURL:    "https://www.etsy.com/oauth/signin",
		AccessTokenURL:  "https://openapi.etsy.com/v2/oauth/access_token",
	}
	authConfig := oauth1.Config{
		ConsumerKey:    c.Config.APIKey,
		ConsumerSecret: c.Config.SharedSecret,
		CallbackURL:    "",
		Endpoint:       endpoint,
	}

	requestToken, requestSecret, err := authConfig.RequestToken()
	if err != nil {
		log.Fatal("Error Requesting Token: ", err)
	}
	log.Println("Token Details: ", requestToken, requestSecret)

	authorizeURL, err := authConfig.AuthorizationURL(requestToken)
	if err != nil {
		log.Fatal("Error Producing Authorize Url: ", err)
	}
	log.Println("Open this url in your browser: ", authorizeURL)

	fmt.Printf("Choose whether to grant the application access.\nPaste " +
		"the oauth_verifier parameter from the address bar: ")
	var verifier string
	_, err = fmt.Scanf("%s", &verifier)
	if err != nil {
		log.Fatalln("Error Reading Input: ", err)
	}
	accessToken, accessSecret, err := authConfig.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		log.Fatalln("Error Getting Access Token: ", err)
	}

	token := oauth1.NewToken(accessToken, accessSecret)
	log.Println("Access Granted!")
	c.Client = authConfig.Client(context.TODO(), token)
	return c.Client
}

func (c Client) makeRequest(endpoint string, params Parameters) {
	res, err := c.Client.Post("https://openapi.etsy.com/v2/listings")
	if err != nil {
		log.Fatalln("Error Making Request: ", err)
	}
	log.Println(res)
}

//AddListings creates multiple listings on an Etsy account.
func (c Client) AddListings() bool {
	params := Parameters{
		Quantity: 1,
	}
	c.makeRequest("", params)
	return true
}
