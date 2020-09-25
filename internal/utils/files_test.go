package utils

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestReadLinesFromFilePath(t *testing.T) {
	lines := ReadLinesFromFilePath("../../fixtures/nine-letter-words.txt")
	test.AssertDeepEqual(lines[0], "aardvarks", t)
}

func TestReadLinesFromStdin(t *testing.T) {
	tmpFile := writeTmpFile([]byte("aardvarks"))

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin
	os.Stdin = tmpFile

	lines, _ := ReadLinesFromStdin()
	test.AssertDeepEqual(lines[0], "aardvarks", t)
}

func writeTmpFile(content []byte) *os.File {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	return tmpfile
}
