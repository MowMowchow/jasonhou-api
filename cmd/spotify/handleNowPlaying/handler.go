package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"jasonhou-api/internal/responses"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	httpClient HttpClient
}

type HttpClient interface {
}

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (, error) {
	CLIENT_ID := ""
	CLIENT_SECRET := ""

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", nil)
	if err != nil {
		errMsg := "error building post request for spotify api token"
		log.Println(errMsg)
		return responses.ServerError(err), fmt.Errorf(errMsg)

	}

	queries := req.URL.Query()
	queries.Add("grant_type", "authorization_code")
	queries.Add("json", "true")

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString((bytes.NewBufferString(CLIENT_ID+":"+CLIENT_SECRET)).Bytes()))
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return nil, nil
}
