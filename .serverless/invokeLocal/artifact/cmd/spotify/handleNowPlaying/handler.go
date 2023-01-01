package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jasonhou-api/internal/models"
	"jasonhou-api/internal/responses"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	httpClient HttpClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	CLIENT_ID := os.Getenv("SPOTIFY_CLIENT_ID")
	CLIENT_SECRET := os.Getenv("SPOTIFY_CLIENT_SECRET")

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

	resp, err := h.httpClient.Do(req)
	if err != nil {
		errMsg := "error with request for spotify api token request"
		log.Println(errMsg)
		return responses.ServerError(err), fmt.Errorf(errMsg)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMsg := "Error with ioutil.ReadAll(resp.Body)"
		log.Println(errMsg, err)
		return responses.ServerError(err), fmt.Errorf(errMsg)
	}

	log.Println("respBody PRINTING: ", respBody, " | STRING", string(respBody))

	var currentPlayingResponse models.CurrentPlayingResponse
	err = json.Unmarshal(respBody, &currentPlayingResponse)
	if err != nil {
		errMsg := "error unmarshalling response body from spotify user profile get request"
		log.Println(errMsg, err)
		return responses.ServerError(err), fmt.Errorf(errMsg)
	}

	responseBody, err := json.Marshal(currentPlayingResponse)
	if err != nil {
		errMsg := ("error marshalling currentPlayingResponse to response body")
		log.Println(errMsg)
		return responses.ServerError(err), fmt.Errorf(errMsg)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{},
		Body:       string(responseBody),
	}, nil
}
