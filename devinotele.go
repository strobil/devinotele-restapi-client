package devinotele

import (
	"encoding/json"
	"errors"
	_ "fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const API_URL = "https://integrationapi.net/rest/v2/"

type DevinoTele struct {
	login    string
	password string
}

func NewDevinoTele(login string, password string) (*DevinoTele, error) {
	if login == "" || password == "" {
		return nil, errors.New("login or password can not be empty")
	}

	m := DevinoTele{}
	m.login = login
	m.password = password

	return &m, nil
}

func (m *DevinoTele) SendSms(from string, to string, text string) ([]string, error) {

	resp, err := http.PostForm(API_URL+"/Sms/Send",
		url.Values{
			"Login":              {m.login},
			"Password":           {m.password},
			"SourceAddress":      {from},
			"DestinationAddress": {to},
			"Data":               {text},
		},
	)

	if err != nil {
		return []string{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var messageIds []string
	json.Unmarshal(body, &messageIds)

	return messageIds, nil
}
