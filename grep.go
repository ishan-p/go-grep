package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type Input struct {
	searchQuery  string
	rootFilePath string
	inputFiles   []string
}

type Result struct {
	channelInfo  map[string]chan string
	searchResult map[string][]string
	mu           sync.Mutex
}

func Grep(searchQuery string, inputFile string) {
	input := Input{searchQuery, inputFile, []string{}}
	input.parseFiles()
	result := Result{channelInfo: make(map[string]chan string), searchResult: make(map[string][]string)}

	for _, file := range input.inputFiles {
		fileHandler, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fileHandler.Close()
		fileBuffCh := make(chan string)
		result.channelInfo[file] = fileBuffCh
		result.searchResult[file] = []string{}
		go readFileByLine(fileHandler, fileBuffCh)
	}
	var wg sync.WaitGroup
	for _, file := range input.inputFiles {
		wg.Add(1)
		go result.gatherResult(input.searchQuery, file, &wg)
	}
	wg.Wait()
	for f, op := range result.searchResult {
		for _, val := range op {
			fmt.Print(string("\033[35m"), f, string("\033[0m"), ": ", val, "\n")
		}
	}
}

func (input *Input) parseFiles() {
	err := filepath.Walk(input.rootFilePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				input.inputFiles = append(input.inputFiles[:], path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func (result *Result) gatherResult(searchQuery string, fileName string, wg *sync.WaitGroup) {
	defer wg.Done()
	for text := range result.channelInfo[fileName] {
		searchResult := find(searchQuery, text)
		if len(searchResult) > 0 {
			result.mu.Lock()
			result.searchResult[fileName] = append(result.searchResult[fileName], fmt.Sprintf("%s", searchResult))
			result.mu.Unlock()
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

func find(searchQuery string, text string) []byte {
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
