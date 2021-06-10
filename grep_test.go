package main

import (
	"fmt"
	"testing"
)

func TestListFilesInDirLen(t *testing.T) {
	inputDir := "./tests"
	countFilesInDir := 3
	val := listFilesInDir(inputDir)
	if len(val) != countFilesInDir {
		t.Fatalf(`listFilesInDir(%v) = %v, want %v`, inputDir, len(val), countFilesInDir)
	}
}

type FindTestCase struct {
	searchQuery string
	inputText   string
	output      string
}

func TestFind(t *testing.T) {
	testCases := []FindTestCase{
		{"test", "Hello world this is a test case", "Hello world this is a \033[31mtest\033[0m case"},
		{"in the", "All words in the world", "All words \033[31min the\033[0m world"},
		{"in ", "All words in the world", "All words \033[31min \033[0mthe world"},
		{"wor.?", "All words in the world", "All \033[31mword\033[0ms in the \033[31mworl\033[0md"},
		{"import", "All words in the world", ""},
	}
	for _, test := range testCases {
		output := find(test.searchQuery, test.inputText)
		outputStr := fmt.Sprintf("%s", output)
		if outputStr != test.output {
			t.Fatalf(`Test find(%v), expected - %v, got - %v`, test.searchQuery, test.output, outputStr)
		}
	}

}
