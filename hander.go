package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)


var ac = auth.account.getAuth();
Auth := c.getAuth()
clientID :=auth.account.clientID
fmt.Println(clientID)

var websiteOauthConfig = &oauth2.Config{
	ClientID:     str,
	ClientSecret: "6f5dae1fe00eb0aa0af931e8e249de8fa76fdacd",
	RedirectURL:  "http://localhost:9094/oauth2",
	Scopes:       []string{"user","project"},

}


func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, main.htmlIndex)
}

//https://github.com/login

func handleWebsiteLogin(w http.ResponseWriter, r *http.Request) {
	url := websiteOauthConfig.AuthCodeURL(main.oauthStateString)
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleWebsiteCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != main.oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", main.oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Println(state)

	code := r.FormValue("code")
	fmt.Println(code)
	token, err := main.githubOauthConfig.Exchange(context.Background(), code)
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

