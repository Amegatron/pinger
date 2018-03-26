package config

import (
	"bufio"
	"log"
	"os"
)

var remotes []string

func init() {
	remotes = scanRemotes()
}

func scanRemotes() []string {
	var tmpRemotes []string

	file, err := os.Open("./remotes.txt")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tmpRemotes = append(tmpRemotes, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		panic(err)
	}

	return tmpRemotes
}

func Remotes() []string {
	return remotes
}
