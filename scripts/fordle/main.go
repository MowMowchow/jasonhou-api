package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type trieNode struct {
	children map[string]*trieNode
	isWord   bool
	word     string
}

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
	log.Println("SQL CLIENT PING SUCCESS")

	// word list file setup
	data, err := os.ReadFile("wordList.txt")
	if err != nil {
		panic(err)
	}
	wordList := strings.Fields(string(data))

	// trie and other setup
	batchList = make(map[string]*Batch)
	makeTrie(wordList)

	// event loop
	log.Println("PRESS Q TO QUIT")
	for true {
		inputs := readLineAsList()

		switch inputs[0] {
		case "gb": // generate batch
			log.Println("(batchName) (boardSize) (batchSize) (cropAmount)")
			inputs = readLineAsList()
			batchName := inputs[0]
			boardSize, _ := strconv.Atoi(inputs[1])
			batchSize, _ := strconv.Atoi(inputs[2])
			cropAmount := inputs[3]
			generateBoard(batchName, boardSize, batchSize, cropAmount)

		case "cb": //compress batches
			log.Println("(boardSize) (compressAmount) (newBatchName)")
			inputs = readLineAsList()
			boardSize, _ := strconv.Atoi(inputs[0])
			compressAmount := inputs[1]
			newBatchName := inputs[2]
			compressBatches(boardSize, compressAmount, newBatchName)

		case "db": // delete batch
			log.Println("(batchName)")
			inputs = readLineAsList()
			batchName := inputs[0]
			deleteBatch(batchName)

		case "sb": // save boards (size specific)
			log.Println("(boardSize)")
			inputs = readLineAsList()
			boardSize, _ := strconv.Atoi(inputs[0])
			saveBoards(sqlClient, boardSize)

		case "pb": // print batch
			log.Println("(batchName)")
			inputs = readLineAsList()
			batchName := inputs[0]
			printBatch(batchName)

		case "lb":
			listBatches()

		case "cDB": // clear table of boards of a specific size
			log.Println("(boardSize)")
			inputs = readLineAsList()
			boardSize, _ := strconv.Atoi(inputs[0])
			clearAllBoards(sqlClient, boardSize)

		case "q": // quit
			return

		default:
			log.Println("UNKNOWN COMMAND")
		}
	}
}

// func main() {
// 	board := []string{
// 		"oaan",
// 		"etae",
// 		"ihkr",
// 		"iflv",
// 	}
// 	wordList := []string{"oath", "pea", "eat", "rain", "hklf", "hf"}
// 	makeTrie(wordList)

// 	log.Println(findWords(board, trie))
// }
