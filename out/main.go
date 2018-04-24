package main

import (
	"fmt"
	"os"
	"bufio"
	"encoding/json"
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
	Secret   string    `json:"secret,omitempty"`
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

	output, err := json.Marshal(Output{})
	if err != nil {
		panic(err)
	}
	fmt.Print(string(output))
}
