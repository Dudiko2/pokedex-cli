package main

import (
	"bufio"
	"strings"
)

type parsedInput struct {
	command   string
	arguments []string
}

func parseInput(input string) parsedInput {
	words := strings.Fields(input)
	wordsLen := len(words)
	if wordsLen < 1 {
		return parsedInput{}
	}
	if wordsLen < 2 {
		return parsedInput{
			command:   words[0],
			arguments: nil,
		}
	}
	p := parsedInput{
		command:   words[0],
		arguments: words[1:],
	}
	return p
}

func prepareInput(scanner *bufio.Scanner) parsedInput {
	scanner.Scan()
	input := scanner.Text()
	sanitized := strings.TrimSpace(input)
	parsed := parseInput(sanitized)
	return parsed
}
