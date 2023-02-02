package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"jasonhou-api/internal/models"
	"jasonhou-api/internal/responses"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	sqlClient SqlClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type SqlClient interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

var ctx = context.Background()

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	boardSizeStr, exists := request.PathParameters["boardSize"]

	if !exists {
		errMsg := "no board size provided"
		log.Println(errMsg)
		return responses.ServerError(errMsg), fmt.Errorf(errMsg)
	}

	boardSize, err := strconv.Atoi(boardSizeStr)
	if err != nil {
		errMsg := "invalid board size | board size must be int"
		log.Println(errMsg)
		return responses.ServerError(err), fmt.Errorf(errMsg)
	}

	rows, err := h.sqlClient.Query(fmt.Sprintf("SELECT * FROM boards_%d ORDER BY RAND() LIMIT 1;", boardSize))
	if err != nil {
		errMsg := fmt.Sprintf("ERROR WITH QUERY TO FETCH BOARD OF SIZE: %d", boardSize)
		log.Println(errMsg)
		return responses.ServerError(err), fmt.Errorf(errMsg)
	}

	defer rows.Close()

	rows.Next()
	board := make([]string, boardSize+1)
	pBoard := make([]interface{}, boardSize+1)
	for i, _ := range board {
		pBoard[i] = &board[i]
	}
	err = rows.Scan(pBoard...)
	if err != nil {
		errMsg := fmt.Sprintf("ERROR WITH QUERY TO FETCH BOARD OF SIZE %d", boardSize)
		log.Println(errMsg)
		return responses.ServerError(err), fmt.Errorf(errMsg)
	}

	words := board[boardSize:]
	board = board[:boardSize]

	fordleResponse := models.FordleResponse{
		Board: board,
		Words: words,
	}

	// re-package request body
	responseBody, err := json.Marshal(fordleResponse)
	if err != nil {
		errMsg := ("error marshalling fordleResponse to response body")
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
