package main

import (
	"fmt"
	"os"
	"bufio"
	"encoding/json"
	"github.com/mmcdole/gofeed"
	"crypto/sha1"
	"encoding/hex"
	"net/url"
	"path/filepath"
	"io/ioutil"
)

type Datum struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Version struct {
	Ref string `json:"ref,omitempty"`
}

type Output struct {
	Version  Version `json:"version,omitempty"`
	Metadata []Datum `json:"metadata,omitempty"`
}

type Source struct {
	Url string `json:"url,omitempty"`
}

type Payload struct {
	Source  Source  `json:"source,omitempty"`
	Version Version `json:"version,omitempty"`
}

func main() {
	args := os.Args
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
		panic(err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(payload.Source.Url)
	if err != nil {
		panic(err)
	}

	output, err := json.Marshal(feed)
	if err != nil {
		panic(err)
	}

	algo := sha1.New()
	_, err = algo.Write(output)
	if err != nil {
		panic(err)
	}

	hash := hex.EncodeToString(algo.Sum(nil))
	if payload.Version.Ref != hash {
		fmt.Fprintf(os.Stderr, "invalid hash: %v", hash)
		os.Exit(1)
	}

	u, err := url.Parse(payload.Source.Url)
	if err != nil {
		panic(err)
	}

	file := filepath.Join(args[1], filepath.Base(u.Path))
	err = ioutil.WriteFile(file, output, 0644)
	if err != nil {
		panic(err)
	}

	output, err = json.Marshal(Output{
		Version{hash},
		[]Datum{}},
)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(output))
}
