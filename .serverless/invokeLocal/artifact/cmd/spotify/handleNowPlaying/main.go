package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	httpClient := &http.Client{}
	// redis client setup

	handler := Handler{httpClient: httpClient}

	lambda.Start(handler.HandleRequest)

}
