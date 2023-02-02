package main

import (
	"database/sql"
	"fmt"
	"log"
)

func clearAllBoards(sqlClient *sql.DB, boardSize int) {
	if _, exists := FORDLE_BOARDS_TABLE_NAMES[boardSize]; !exists {
		fmt.Printf("NO TABLE FOR BOARD OF SIZE %d EXISTS\n", boardSize)
		return
	}
	_, err := sqlClient.Query(fmt.Sprintf("TRUNCATE %s;", FORDLE_BOARDS_TABLE_NAMES[boardSize]))
	if err != nil {
		log.Println("ERROR WITH CLEARING BOARDS OF SIZE", boardSize, "| ERROR:", err)
		return
	}
	log.Println("SUCCESSFULLY CLEARED BOARDS OF SIZE", boardSize)
}

func AddBoards(sqlClient *sql.DB, boards []Board, boardSize int) {
	// get cols
	cols := ""
	for i := 1; i <= boardSize; i++ {
		cols += fmt.Sprintf("r%d", i)
		cols += ", "
	}
	cols += "words"

	for boardNum := 0; boardNum < len(boards); boardNum++ {
		// get cols to inserts
		insCols := ""
		for i := 0; i < boardSize; i++ {
			insCols += fmt.Sprintf("'%s'", boards[boardNum].board[i])
			insCols += ", "
		}

		wordsForInsCol := "'"
		for _, word := range boards[boardNum].words {
			wordsForInsCol += word
			wordsForInsCol += " "
		}
		wordsForInsCol += "'"
		insCols += wordsForInsCol

		_, err := sqlClient.Query(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", FORDLE_BOARDS_TABLE_NAMES[boardSize], cols, insCols))
		if err != nil {
			log.Println("ERROR WHEN UPLOADING BOARDS | ERR", err)
		}
	}
}
