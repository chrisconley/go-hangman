package utils

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

func ReadLinesFromFilePath(path string) (lines []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func ReadLinesFromStdin() ([]string, error) {
	var lines []string
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		return lines, errors.New("The command is intended to work with pipes: cat ./fixtures/nine-letter-words.txt | go run cmd/player.go")
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("ERROR")
	}

	return lines, nil
}
