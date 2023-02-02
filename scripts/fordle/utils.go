package main

import (
	"bufio"
	"os"
	"strings"
)

func readLineAsList() []string {
	r := bufio.NewReader(os.Stdin)
	inputString, _ := r.ReadString('\n')
	// convert CRLF to LF
	inputString = strings.Replace(inputString, "\n", "", -1)
	return strings.Fields(inputString)
}
