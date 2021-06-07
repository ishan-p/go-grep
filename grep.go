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

const maxOpenFileDescriptors = 1000

func Grep(searchQuery string, inputFile string) {
	input := Input{searchQuery, inputFile, []string{}}
	input.parseFiles()

	result := Result{channelInfo: make(map[string]chan string), searchResult: make(map[string][]string)}

	var wg sync.WaitGroup

	resultSyncChannel := make(chan string)
	quit := make(chan int)
	go result.writeToStdout(resultSyncChannel, quit)

	fileDescriptorBuffer := make(chan int, maxOpenFileDescriptors)
	for _, file := range input.inputFiles {
		fileBuffCh := make(chan string)
		result.mu.Lock()
		result.channelInfo[file] = fileBuffCh
		result.searchResult[file] = []string{}
		result.mu.Unlock()
		go readFileByLine(file, fileBuffCh, fileDescriptorBuffer)
		wg.Add(1)
		go result.gatherResult(input.searchQuery, file, resultSyncChannel, &wg)
	}

	wg.Wait()

	quit <- 1
	close(resultSyncChannel)
	close(quit)
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

func (result *Result) writeToStdout(resultSyncChannel chan string, quit chan int) {
	for {
		select {
		case fileName := <-resultSyncChannel:
			result.mu.Lock()
			for _, op := range result.searchResult[fileName] {
				fmt.Print(string("\033[35m"), fileName, string("\033[0m"), ": ", op, "\n")
			}
			result.mu.Unlock()
		case <-quit:
			return
		}
	}
}

func (result *Result) gatherResult(searchQuery string, fileName string, resultSyncChannel chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	result.mu.Lock()
	fileReadChannel := result.channelInfo[fileName]
	result.mu.Unlock()
	for text := range fileReadChannel {
		searchResult := find(searchQuery, text)
		if len(searchResult) > 0 {
			result.mu.Lock()
			result.searchResult[fileName] = append(result.searchResult[fileName], fmt.Sprintf("%s", searchResult))
			result.mu.Unlock()
		}
	}
	// Send message to result sync channel denoting current file search is complete
	resultSyncChannel <- fileName
}

func readFileByLine(filePath string, channel chan string, fileDescriptorBuffer chan int) {
	fileDescriptorBuffer <- 1
	fileHandler, err := os.Open(filePath)
	defer func() {
		fileHandler.Close()
		close(channel)
		<-fileDescriptorBuffer
	}()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(fileHandler)
	for scanner.Scan() {
		channel <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(filePath, err)
	}
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
