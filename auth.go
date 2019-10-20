package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)





type auth struct {
	Account struct
{
	clientID string `yaml:"clientID"`
	clientSecret string `yaml:"clientSecret"`
	redirectURL string `yaml:"redirectURL"`
	scopes []string `yaml:"scopes"`
	endpoint string `yaml:"endpoint"`
}
}
var c = auth.Account
fmt.Println(c)

type account struct {
	clientID string `yaml:"clientID"`
	clientSecret string `yaml:"clientSecret"`
	redirectURL string `yaml:"redirectURL"`
	scopes []string `yaml:"scopes"`
	endpoint string `yaml:"endpoint"`
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

