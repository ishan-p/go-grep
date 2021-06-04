package gogrep

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func Grep(searchQuery string, inputFile string) {
	fileHandler, err := os.Open(inputFile)

	if err != nil {
		log.Fatal(err)
	}

	defer fileHandler.Close()
	fileBuffCh := make(chan string)
	go readFileByLine(fileHandler, fileBuffCh)
	result := []string{}
	for text := range fileBuffCh {
		searchResult := Find(searchQuery, text)
		if len(searchResult) > 0 {
			result = append(result, fmt.Sprintf("%s", searchResult))
			fmt.Printf("%s\n", searchResult)
		}
	}
}

func readFileByLine(fileHandler *os.File, channel chan string) {
	scanner := bufio.NewScanner(fileHandler)
	for scanner.Scan() {
		channel <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	close(channel)
}

func Find(searchQuery string, text string) []byte {
	re := regexp.MustCompile(searchQuery)
	input := []byte(text)
	result := re.FindAll([]byte(text), -1)
	var output []byte
	if len(result) > 0 {
		for _, val := range result {
			coloredOp := fmt.Sprintf("\033[31m%s\033[0m", val)
			output = re.ReplaceAll(input, []byte(coloredOp))
		}
	}
	return output
}
