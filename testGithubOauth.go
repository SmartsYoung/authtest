package main

import (
	"authtest/authentication"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

func main() {

	var c authtest.auth
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

