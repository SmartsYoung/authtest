package main


import (
	"fmt"
	"golang.org/x/oauth2/github"

	"gopkg.in/yaml.v2"

	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

const htmlIndex = `<html><body>
<a href="/login">Log in with Github</a>
</body></html>
`



//func (a &oauth2.Config)


//var accouts=authentication.WithAccount("5227d176177edcdcb5e0")
//var secret =authentication.WithSecret("6f5dae1fe00eb0aa0af931e8e249de8fa76fdacd")

var str  string = "5227d176177edcdcb5e0"
// Auth.Authorization(authOptions ...ConfigOption) (string, error)

var githubOauthConfig = &oauth2.Config{
	ClientID:    str,
	ClientSecret: "6f5dae1fe00eb0aa0af931e8e249de8fa76fdacd",
	RedirectURL:  "http://localhost:9094/oauth2",
	Scopes: []string{"user","project"},

}

const oauthStateString = "random"           //状态码

func main() {

	var c auth
	auth := c.getAuth()
	clientID :=auth.account.clientID
	fmt.Println(clientID)


	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGithubLogin)
	http.HandleFunc("/oauth2", handleGithubCallback)
	fmt.Println(http.ListenAndServe(":9094", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

//https://github.com/login

func handleGithubLogin(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL(oauthStateString)
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGithubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Println(state)

	code := r.FormValue("code")
	fmt.Println(code)
	token, err := githubOauthConfig.Exchange(oauth2.NoContext, code)
	fmt.Println(token)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	header := http.Header{}
	header.Set("Accept", "application/json")
	header.Set("Content-Type", "application/json")
	header.Set("Authorization","token "+token.AccessToken)
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		fmt.Println("new request failed", err)
		return
	}

	req.Header = header
	response, err := http.DefaultClient.Do(req)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Fprintf(w, "Content: %s\n", contents)
}

type auth struct {
   account
}


type account struct {
	clientID string `yaml:"host"`
	clientSecret string `yaml:"user"`
	redirectURL string `yaml:"pwd"`
	scopes []string `yaml:"pwd"`
	endpoint string `yaml:"dbname"`
}

func (ac *account)ClientID() string{
	return ac.clientID
}

func (ac *account)ClientSecret() string{
	return ac.clientSecret
}

func (ac *account)RedirectURL() string{
	return ac.redirectURL 
}

func (ac *account)Scopes() []string{
	return ac.scopes
}

func (ac *account)Endpoint() string{
	return ac.endpoint
}

func (c *auth) getAuth() *auth {
	yamlFile, err := ioutil.ReadFile("auth.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

