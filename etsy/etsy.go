package etsy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/dghubble/oauth1"
)

type config struct {
	APIKey       string `json:"API_KEY"`
	SharedSecret string `json:"SHARED_SECRET"`
}

type accessDetails struct {
	ConsumerKey    string `json:"ConsumerKey"`
	ConsumerSecret string `json:"ConsumerSecret"`
	AccessToken    string `json:"AccessToken"`
	AccessSecret   string `json:"AccessSecret"`
}

//Client type for Etsy API calls.
type Client struct {
	Config  config
	BaseURL string
	Client  *http.Client
}

//Parameters to pass to Etsy api.
type Parameters struct {
	Quantity             int     `csv:"Quantity"`
	Title                string  `csv:"Title"`
	Description          string  `csv:"Description"`
	ImageName            string  `csv:"Image"`
	Price                float64 `csv:"Price"`
	TaxonomyName         string  `csv:"Taxonomy Name"`
	WhoMade              string  `csv:"Who Made"`
	IsSupply             bool    `csv:"Has Supply"`
	WhenMade             string  `csv:"When Made"`
	Username             string  `csv:"Username"`
	ShippingTemplateName string  `csv:"Shipping Template Name"`
}

//Taxonomy is a category from the Etsy api.
type Taxonomy struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Children []Taxonomy `json"children"`
}

//TaxonomyList is a list of taxonomies.
type TaxonomyList struct {
	Count   int        `json:"count"`
	Results []Taxonomy `json:"results"`
}

//User represents a user in the Etsy api.
type User struct {
	ID        int    `json:"user_id"`
	LoginName string `json:"login_name"`
}

//UserList represents a list of users from the Etsy api.
type UserList struct {
	Users []User `json:"results"`
}

//ShippingTemplate represents a shipping template from the Etsy api.
type ShippingTemplate struct {
	ID   int    `json:"shipping_template_id"`
	Name string `json:"title"`
}

//ShippingTemplateList represents a list of shipping templates from the Etsy api.
type ShippingTemplateList struct {
	Templates []ShippingTemplate `json:"results"`
}

//NewClient creates a new Etsy client to make request to the etsy api
func NewClient() Client {
	client := Client{
		BaseURL: "https://openapi.etsy.com/v2",
		Client:  Authenticate(),
	}
	return client
}

//checkCredentials reads credentials from environment variables
func checkCredentials() (bool, *http.Client) {
	authInput, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Println("Error Reading Credentials File: ", err)
		return true, &http.Client{}
	}
	accessDetails := accessDetails{}

	json.Unmarshal([]byte(authInput), &accessDetails)
	if accessDetails.ConsumerKey == "" || accessDetails.ConsumerSecret == "" || accessDetails.AccessToken == "" || accessDetails.AccessSecret == "" {
		log.Println("Missing access information.")
		return true, &http.Client{}
	}
	config := oauth1.NewConfig(accessDetails.ConsumerKey, accessDetails.ConsumerSecret)
	token := oauth1.NewToken(accessDetails.AccessToken, accessDetails.AccessSecret)
	// httpClient will automatically authorize http.Request's
	httpClient := config.Client(oauth1.NoContext, token)
	return false, httpClient
}

//Authenticate performs an Oauth1.0 authentication with the Etsy api.
func Authenticate() *http.Client {
	requestAuth, authClient := checkCredentials()
	if !requestAuth {
		return authClient
	}
	configFile := readConfig()

	endpoint := oauth1.Endpoint{
		RequestTokenURL: "https://openapi.etsy.com/v2/oauth/request_token?scope=listings_w%20listings_r",
		AuthorizeURL:    "https://www.etsy.com/oauth/signin",
		AccessTokenURL:  "https://openapi.etsy.com/v2/oauth/access_token",
	}
	authConfig := oauth1.Config{
		ConsumerKey:    configFile.APIKey,
		ConsumerSecret: configFile.SharedSecret,
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
	saveAuthInfo(authConfig.ConsumerKey, authConfig.ConsumerSecret, accessToken, accessSecret)
	//save auth info
	return authConfig.Client(oauth1.NoContext, token)
}

//readConfig reads config.json file to obtain oauth api key and api secret.
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

//saveAuthInfo saves oauth credentils to credentials.json
func saveAuthInfo(consumerKey string, consumerSecret string, accessToken string, accessSecret string) {
	authOutput := accessDetails{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		AccessToken:    accessToken,
		AccessSecret:   accessSecret,
	}
	file, _ := json.MarshalIndent(authOutput, "", " ")
	_ = ioutil.WriteFile("credentials.json", file, 0644)
}

func (c Client) makePostRequest(url string) {
	log.Println(url)
	res, err := c.Client.Post(url, "application/json", bytes.NewBufferString("{}"))
	if err != nil {
		log.Fatalln("Error Making Request: ", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("Error Reading Body: ", err)
	}
	log.Print("Body: ", string(body))
}

func (c Client) makeGetRequest(endpoint string) []byte {
	url := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
	res, err := c.Client.Get(url)
	if err != nil {
		log.Println("Error Getting Active Listings: ", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("Error Reading Body: ", err)
	}
	//log.Println("Body: ", string(body))
	return body
}

//GetActiveListings returns the active listings for an Etsy account.
func (c Client) GetActiveListings() {
	c.makeGetRequest("listings/active")
}

//GetShop gets a shop id given a shope name
func (c Client) GetShop(shopName string) {
	endpoint := fmt.Sprintf("shops?shop_name=%s", shopName)
	c.makeGetRequest(endpoint)
}

//filterTaxonomies searches through list of taxonomies of any depth to find id of matching name.
func (c Client) filterTaxonomies(name string, taxonomies []Taxonomy) (int, bool) {
	for _, v := range taxonomies {
		if v.Name == name {
			return v.ID, true
		}
		if len(v.Children) > 0 {
			recursiveResult, success := c.filterTaxonomies(name, v.Children)
			if success {
				return recursiveResult, true
			}
		}
	}
	return 0, false
}

//FindTaxonomy gets the taxonomy as used by sellers in the listing process.
func (c Client) FindTaxonomy(name string) int {
	endpoint := "taxonomy/seller/get"
	taxonomyList := TaxonomyList{}
	response := c.makeGetRequest(endpoint)
	json.Unmarshal(response, &taxonomyList)
	filteredTaxonomy, success := c.filterTaxonomies(name, taxonomyList.Results)
	if !success {
		log.Fatalf("Taxonomy %s Not Found", name)
	}
	return filteredTaxonomy
}

//FindUser searches for a user by username and returns their id.
func (c Client) FindUser(username string) int {
	endpoint := fmt.Sprintf("users?keywords=%s", username)
	response := c.makeGetRequest(endpoint)
	userList := UserList{}
	json.Unmarshal(response, &userList)
	return userList.Users[0].ID
}

//FindUserShippingTemplate will search the Etsy api for a user's shipping template by name.
func (c Client) FindUserShippingTemplate(userID int, templateName string) int {
	endpoint := fmt.Sprintf("users/%d/shipping/templates", userID)
	response := c.makeGetRequest(endpoint)
	templates := ShippingTemplateList{}
	json.Unmarshal(response, &templates)

	for _, template := range templates.Templates {
		if template.Name == templateName {
			return template.ID
		}
	}
	log.Fatalf("No Template With Name: %s Within User %d Found", templateName, userID)
	return 0
}

//AddListings creates multiple listings on an Etsy account.
func (c Client) AddListings(params Parameters, image []byte) bool {
	taxonomyID := c.FindTaxonomy(params.TaxonomyName)
	userID := c.FindUser(params.Username)
	shippingTemplateID := c.FindUserShippingTemplate(userID, params.ShippingTemplateName)
	query := fmt.Sprintf(
		"quantity=%d&title=%s&description=%s&price=%f&taxonomy_id=%d&who_made=%s&is_supply=%t&when_made=%s&shipping_template_id=%d&image=%s",
		params.Quantity,
		params.Title,
		params.Description,
		params.Price,
		taxonomyID,
		params.WhoMade,
		params.IsSupply,
		params.WhenMade,
		shippingTemplateID,
		image,
	)
	path := fmt.Sprintf("%s/%s?%s", c.BaseURL, "listings", url.PathEscape(query))
	c.makePostRequest(path)
	return true
}
