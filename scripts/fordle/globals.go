package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"reflect"
)

/* ~~~~~~~~~ GLOBAL VARS ~~~~~~~~~ */

var moves [][]int

type CurrentBoard struct {
	board    []string
	wordList map[string]bool
	R        int
	C        int
}

type Board struct {
	board     []string
	words     []string
	wordCount int
}

type Batch struct {
	boards    []Board
	size      int
	boardSize int
}

type Coordinate struct {
	x int
	y int
}

var trie *trieNode
var batchList map[string]*Batch

/* ~~~~~~~~~ CONSTANTS ~~~~~~~~~ */
var FORDLE_BOARDS_TABLE_NAMES = map[int]string{
	4: "boards_4",
}

/* ~~~~~~~~~ FUNCTIONS ~~~~~~~~~ */

func trieDfs(word string, ind int, t *trieNode) {
	if ind == len(word) {
		return
	}

	if t == nil {
		return
	}

	if t.children == nil {
		t.children = make(map[string]*trieNode)
	}

	if _, exists := t.children[string(word[ind])]; !(exists) {
		t.children[string(word[ind])] = new(trieNode)
	}

	if ind == len(word)-1 {
		t.children[string(word[ind])].isWord = true
		t.children[string(word[ind])].word = word
	}

	trieDfs(word, ind+1, t.children[string(word[ind])])
}

func printTrie(t *trieNode) {
	log.Println(t)
	for key := range t.children {
		log.Println("GOING TO LET:", key)
		printTrie(t.children[key])
	}
}

func makeTrie(words []string) {
	trie = new(trieNode)
	for _, word := range words {
		trieDfs(word, 0, trie)
	}
}

func Clone(a interface{}) interface{} {
	buff := new(bytes.Buffer)
	v := reflect.New(reflect.TypeOf(a))
	gob.NewEncoder(buff).Encode(a)
	gob.NewDecoder(buff).Decode(v.Interface())
	return v.Elem().Interface()
}

func findWordsDfs(cBoard CurrentBoard, x int, y int, t *trieNode, vis map[Coordinate]bool) bool {
	currCoor := Coordinate{x: x, y: y}
	wordExists := false
	if _, exists := vis[currCoor]; !exists {
		vis[currCoor] = true
		if t.isWord {
			cBoard.wordList[t.word] = true
			wordExists = true
		}
		for _, move := range moves {
			nx := x + move[0]
			ny := y + move[1]
			if (0 <= nx && nx < cBoard.C) && (0 <= ny && ny < cBoard.R) {
				if _, exists := t.children[string(cBoard.board[ny][nx])]; exists {
					newVis := map[Coordinate]bool{}
					for coor, val := range vis {
						newVis[coor] = val
					}
					findWordsDfs(cBoard, nx, ny, t.children[string(cBoard.board[ny][nx])], newVis)
				}
			}
		}
	}
	return wordExists
}

func findWords(cBoard CurrentBoard, trie *trieNode) []string {
	moves = [][]int{{0, 1}, {1, 0}, {-1, 0}, {0, -1}, {1, 1}, {-1, -1}, {-1, 1}, {1, -1}}
	for y := 0; y < cBoard.R; y++ {
		for x := 0; x < cBoard.C; x++ {
			if _, exists := trie.children[string(cBoard.board[y][x])]; exists {
				findWordsDfs(cBoard, x, y, trie.children[string(cBoard.board[y][x])], map[Coordinate]bool{})
			}
		}
	}

	words := []string{}
	for word, _ := range cBoard.wordList {
		words = append(words, word)
	}
	return words
}
