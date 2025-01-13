package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

var reKV = regexp.MustCompile(`^\s*(\w+)\s*(.*)\s*$`)
var scale = flag.Int("n", 2, "scaling factor")
var name = flag.String("name", "", "scaling factor")

func main() {
	flag.Parse()
	state := NewState(*scale, *name)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if state.Scale <= 1 {
			fmt.Println(line)
			continue
		}

		match := reKV.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		state.K = ""
		if len(match) > 1 {
			state.K = match[1]
		}
		state.V = ""
		if len(match) > 2 {
			state.V = match[2]
		}

		err := state.Next()
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERR:", err)
			os.Exit(1)
		}
	}
}
