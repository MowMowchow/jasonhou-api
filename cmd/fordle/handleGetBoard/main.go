package main

import (
	"database/sql"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// mariadb client setup
	sqlClient, err := sql.Open("mysql", "fordleUser:password@tcp(ec2-34-207-157-163.compute-1.amazonaws.com:3306)/fordle")
	if err != nil {
		log.Println("SQL CLIENT ERR:", err)
	}
	defer sqlClient.Close()
	err = sqlClient.Ping()
	if err != nil {
		log.Println("SQL CLIENT PING FAIL | err:", err)
	}

	handler := Handler{
		sqlClient: sqlClient,
	}

	lambda.Start(handler.HandleRequest)
}
