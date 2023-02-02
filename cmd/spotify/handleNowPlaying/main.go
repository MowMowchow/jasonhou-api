package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	httpClient := &http.Client{}
	// redis client setup
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "ec2-34-207-157-163.compute-1.amazonaws.com:6379",
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	// mariadb client setup
	sqlClient, err := sql.Open("mysql", "jasonhouUser:password@tcp(ec2-34-207-157-163.compute-1.amazonaws.com:3306)/jasonhou")
	if err != nil {
		log.Println("SQL CLIENT ERR:", err)
	}
	defer sqlClient.Close()
	err = sqlClient.Ping()
	if err != nil {
		log.Println("SQL CLIENT PING FAIL | err:", err)
	}

	handler := Handler{
		httpClient:  httpClient,
		redisClient: redisClient,
		sqlClient:   sqlClient,
	}

	lambda.Start(handler.HandleRequest)

}
