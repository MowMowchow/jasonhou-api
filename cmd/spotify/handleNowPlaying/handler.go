package main

import (
	"bytes"
	"context"
	"database/sql"
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
	"github.com/go-redis/redis/v8"
)

type Handler struct {
	httpClient  HttpClient
	redisClient RedisClient
	sqlClient   SqlClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RedisClient interface {
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
}

type SqlClient interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

var ctx = context.Background()

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	CLIENT_ID := os.Getenv("SPOTIFY_CLIENT_ID")
	CLIENT_SECRET := os.Getenv("SPOTIFY_CLIENT_SECRET")
	REFRESH_TOKEN := os.Getenv("SPOTIFY_REFRESH_TOKEN")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", nil)
	if err != nil {
		errMsg := "error building post request for spotify api refresh token"
		log.Println(errMsg)
		return responses.ServerError(err), fmt.Errorf(errMsg, err)

	}

	queries := req.URL.Query()
	queries.Add("grant_type", "refresh_token")
	queries.Add("refresh_token", REFRESH_TOKEN)
	req.Header.Add("json", "true")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString((bytes.NewBufferString(CLIENT_ID+":"+CLIENT_SECRET)).Bytes()))
	req.URL.RawQuery = queries.Encode()

	tokenResp, err := h.httpClient.Do(req)
	if err != nil {
		errMsg := "error with request for spotify api refresh token request"
		log.Println(errMsg)
		return responses.ServerError(err), fmt.Errorf(errMsg, err)
	}

	var tokenRespBody models.SpotifyTokenResponseBody
	err = json.NewDecoder(tokenResp.Body).Decode(&tokenRespBody)
	if err != nil {
		errMsg := "error decoding response body from spotify refresh token post request"
		log.Println(errMsg, err)
		return responses.ServerError(err), fmt.Errorf(errMsg, err)
	}

	// prepare another request to get currently playing song
	req, err = http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	req.Header.Add("Authorization", "Bearer "+tokenRespBody.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		errMsg := "error with get request response for user current track"
		log.Println(errMsg, resp.StatusCode, err)
		return responses.ServerError(err), fmt.Errorf(errMsg, resp.StatusCode, err)
	}

	// process response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMsg := "Error with ioutil.ReadAll(resp.Body)"
		log.Println(errMsg, err)
		return responses.ServerError(err), fmt.Errorf(errMsg, err)
	}

	// extract response body
	var currentPlayingResponse models.CurrentPlayingResponse
	err = json.Unmarshal(respBody, &currentPlayingResponse)
	if err != nil {
		errMsg := "error unmarshalling response body from spotify user profile get request"
		log.Println(errMsg, err)
		return responses.ServerError(err), fmt.Errorf(errMsg)
	}

	// re-package request body
	responseBody, err := json.Marshal(currentPlayingResponse)
	if err != nil {
		errMsg := ("error marshalling currentPlayingResponse to response body")
		log.Println(errMsg)
		return responses.ServerError(err), fmt.Errorf(errMsg)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Headers":     "*",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: string(responseBody),
	}, nil
}
