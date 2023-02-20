package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"lukechampine.com/frand"
)

func generateBoard(batchName string, boardSize int, batchSize int, cropAmount string) {
	startTime := time.Now()
	if _, exists := batchList[batchName]; exists {
		log.Println("BATCH NAME ALREADY EXISTS")
		return
	}

	asciiLowerBound := 97
	asciiUpperBound := 123
	tempBoardList := []Board{}
	batchList[batchName] = &Batch{}
	batchList[batchName].boardSize = boardSize

	var mutex sync.Mutex
	var wg sync.WaitGroup

	for batchNum := 0; batchNum < batchSize; batchNum++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			board := CurrentBoard{}
			for row := 0; row < boardSize; row++ {
				rowS := ""
				for col := 0; col < boardSize; col++ {
					rowS += string(asciiLowerBound + (frand.Intn(asciiUpperBound - asciiLowerBound)))
				}
				board.board = append(board.board, rowS)
			}
			board.R = boardSize
			board.C = boardSize
			board.wordList = map[string]bool{}
			tempBoard := Board{}
			tempBoard.board = board.board
			tempBoard.words = findWords(board, trie)
			if len(tempBoard.words) > 0 {
				tempBoard.wordCount = len(tempBoard.words)
				mutex.Lock()
				tempBoardList = append(tempBoardList, tempBoard)
				mutex.Unlock()
			}
		}()
	}
	wg.Wait()

	// sort tempBoardList
	sort.Slice(tempBoardList, func(i, j int) bool {
		return tempBoardList[i].wordCount > tempBoardList[j].wordCount
	})

	var raw bool
	if strings.Contains(cropAmount, "%") {
		raw = false
	} else {
		raw = true
	}

	if raw {
		saveAmt, _ := strconv.Atoi(cropAmount)
		if batchSize < saveAmt {
			log.Println("BAD INPUT")
			return
		}
		batchList[batchName].boards = tempBoardList[:saveAmt]
		batchList[batchName].size = saveAmt
	} else {
		saveAmt, _ := strconv.Atoi(strings.Replace(cropAmount, "%", "", -1))
		if saveAmt < 0 || 100 < saveAmt {
			log.Println("BAD INPUT")
			return
		}
		saveAmt = int(float64(batchSize) * (.01 * float64(saveAmt)))
		batchList[batchName].boards = tempBoardList[:saveAmt]
		batchList[batchName].size = saveAmt
	}
	elapsed := time.Now().Sub(startTime)
	log.Println("DONE GENERATING BOARDS IN:", elapsed)
}

func deleteBatch(batch string) {
	if _, exists := batchList[batch]; !exists {
		log.Println("COULD NOT FIND BATCH:", batch)
		return
	}
	delete(batchList, batch)
	log.Println("DONE DELETEING BATCH: ", batch)
}

func getBoardsOfSize(boardSize int) []Board {
	if _, exists := FORDLE_BOARDS_TABLE_NAMES[boardSize]; !exists {
		fmt.Printf("NO TABLE FOR BOARD OF SIZE %d EXISTS\n", boardSize)
		return nil
	}

	boardsToGet := []Board{}
	for _, batch := range batchList {
		if batch.boardSize == boardSize {
			boardsToGet = append(boardsToGet, batch.boards...)
		}
	}
	return boardsToGet
}

func compressBatches(boardSize int, compressAmount string, newBatchName string) {
	if _, exists := batchList[newBatchName]; exists {
		fmt.Printf("BATCH NAME ALREADY EXISTS: ", newBatchName)
		return
	}

	startTime := time.Now()
	boardsToCompress := getBoardsOfSize(boardSize)
	NboardsToCompress := len(boardsToCompress)
	sort.Slice(boardsToCompress, func(i, j int) bool {
		return boardsToCompress[i].wordCount > boardsToCompress[j].wordCount
	})

	var raw bool
	if strings.Contains(compressAmount, "%") {
		raw = false
	} else {
		raw = true
	}
	batchList[newBatchName] = &Batch{}
	batchList[newBatchName].boardSize = boardSize

	if raw {
		saveAmt, _ := strconv.Atoi(compressAmount)
		if NboardsToCompress < saveAmt {
			log.Println("BAD INPUT")
			return
		}
		batchList[newBatchName].boards = boardsToCompress[:saveAmt]
		batchList[newBatchName].size = saveAmt
	} else {
		saveAmt, _ := strconv.Atoi(strings.Replace(compressAmount, "%", "", -1))
		if saveAmt < 0 || 100 < saveAmt {
			log.Println("BAD INPUT")
			return
		}
		saveAmt = int(float64(NboardsToCompress) * (.01 * float64(saveAmt)))
		batchList[newBatchName].boards = boardsToCompress[:saveAmt]
		batchList[newBatchName].size = saveAmt
	}
	for batchName, batch := range batchList {
		if batch.boardSize == boardSize && batchName != newBatchName {
			deleteBatch(batchName)
		}
	}
	elapsed := time.Now().Sub(startTime)
	log.Println("DONE COMPRESSING BOARDS OF SIZE:", boardSize, "IN:", elapsed)
}

func saveBoards(sqlClient *sql.DB, boardSize int) {
	startTime := time.Now()
	boardsToSave := getBoardsOfSize(boardSize)
	AddBoards(sqlClient, boardsToSave, boardSize)
	elapsed := time.Now().Sub(startTime)
	log.Println("DONE UPLOADING BOARDS OF SIZE:", boardSize, "IN:", elapsed)
}

func clearTable(boardSize int) {
	startTime := time.Now()
	if _, exists := FORDLE_BOARDS_TABLE_NAMES[boardSize]; !exists {
		fmt.Printf("NO TABLE FOR BOARD OF SIZE %d EXISTS\n", boardSize)
		return
	}
	clearTable(boardSize)
	elapsed := time.Now().Sub(startTime)
	log.Println("DONE CLEARING TABLE OF SIZE:", boardSize, "IN:", elapsed)
}

func printBatch(batch string) {
	if _, exists := batchList[batch]; !exists {
		log.Println("COULD NOT FIND BATCH:", batch)
		return
	}

	maxWordCount := 0
	totalCharCnt := 0
	for i, board := range batchList[batch].boards {
		charCnt := 0
		for _, word := range board.words {
			charCnt += len(word)
		}
		fmt.Println("Board:", i, "| Total Words:", board.wordCount, " | Total Chars:", charCnt)
		maxWordCount = int(math.Max(float64(maxWordCount), float64(board.wordCount)))
		for _, row := range board.board {
			fmt.Println(row)
		}
		totalCharCnt += charCnt
	}

	log.Println("TOTAL BOARDS:", batchList[batch].size, "| MAX WORD COUNT:", maxWordCount, "| TOTAL CHAR CNT:", totalCharCnt, "| SIZE OF BATCH (mem):", fmt.Sprintf("%d", unsafe.Sizeof(batchList[batch])))
}

func listBatches() {
	for batchName, batch := range batchList {
		fmt.Println("batch:", batchName, "| size:", batch.size, "| board size:", batch.boardSize)
	}
	log.Println("DONE LISTING BATCHES")
}
