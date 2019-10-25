package devinotele

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	ApiUrl = "https://integrationapi.net/rest/v2/"
)

const (
	DeliveryTypeSms = 1
	DeliveryTypeViber = 2
)

type DevinoTele struct {
	Login    string
	Password string
}

type ErrorResponse struct {
	Code int    `json:"Code"`
	Desc string `json:"Desc"`
}

func NewDevinoTele(login string, password string) (*DevinoTele, error) {
	if login == "" || password == "" {
		return nil, errors.New("login or password can not be empty")
	}

	m := DevinoTele{}
	m.Login = login
	m.Password = password

	return &m, nil
}

func (m *DevinoTele) SendMessage(from string, to string, text string, deliveryType int) (string, error) {
	if from == "" || to == "" || text == "" {
		return "", errors.New("arguments can not be empty")
	}

	urlPart := ""

	switch deliveryType {
	case DeliveryTypeSms:
		urlPart = "/Sms/Send"
	case DeliveryTypeViber:
		urlPart = "/Viber/Send"
	default:
		return "", errors.New("invalid delivery type")
	}

	resp, err := http.PostForm(ApiUrl+urlPart,
		url.Values{
			"Login":              {m.Login},
			"Password":           {m.Password},
			"SourceAddress":      {from},
			"DestinationAddress": {to},
			"Data":               {text},
		},
	)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	case 200:
		var msgId []string
		err := json.Unmarshal(body, &msgId)

		if err != nil {
			return "", err
		}

		return msgId[0], nil
	case 400:
		var response ErrorResponse
		err := json.Unmarshal(body, &response)

		if err != nil {
			return "", err
		}

		return "", errors.New(fmt.Sprintf("Error: %d %s", response.Code, response.Desc))
	case 500:
		return "", errors.New("internal Server Error")
	}

	return "", nil
}