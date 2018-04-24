package main

import (
	"encoding/json"
	"fmt"
	"os"
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"github.com/mmcdole/gofeed"
)

type Version struct {
	Ref string `json:"ref,omitempty"`
}

type Source struct {
	Url string `json:"url,omitempty"`
}

type Payload struct {
	Source  Source  `json:"source,omitempty"`
	Version Version `json:"version,omitempty"`
}

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if stat.Mode()&os.ModeNamedPipe == 0 {
		panic("stdin is empty")
	}

	var payload Payload
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &payload); err != nil {
			panic(err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(payload.Source.Url)

	output, err := json.Marshal(feed)
	if err != nil {
		panic(err)
	}

	algo := sha1.New()
	algo.Write(output)
	output, err = json.Marshal([]Version{{hex.EncodeToString(algo.Sum(nil))}})
	if err != nil {
		panic(err)
	}
	fmt.Print(string(output))
}
