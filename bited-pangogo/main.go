package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var txtPath = flag.String("txt", "", "base text file (.txt)")
var clrPath = flag.String("clr", "", "color overlay file (.clr)")
var clrsStr = flag.String("clrs", "", "newline-separated string of colors 0-15")

func main() {
	flag.Parse()

	txt, err := os.Open(*txtPath)
	if err != nil {
		log.Fatal(err)
	}
	defer txt.Close()

	clr, err := os.Open(*clrPath)
	if err != nil {
		log.Fatal(err)
	}
	defer clr.Close()

	clrs := make(map[rune]string)
	var ks = []rune("0123456789ABCDEF")
	for i, clr := range strings.SplitN(*clrsStr, "\n", 16) {
		clrs[ks[i]] = clr
	}

	state := NewState()
	txtScan := bufio.NewScanner(txt)
	clrScan := bufio.NewScanner(clr)
	nl := false
	for txtScan.Scan() {
		if nl {
			state.Content.WriteRune('\n')
		}
		nl = true

		txtLine := []rune(txtScan.Text())
		var clrLine []rune
		if clrScan.Scan() {
			clrLine = []rune(clrScan.Text())
		}

		for i, c := range txtLine {
			if len(clrLine) > i {
				k := clrLine[i]
				switch k {
				case '.':
					state.Next(nil)
				default:
					if clr, ok := clrs[k]; ok {
						state.Next(&clr)
					}
				}
			}
			state.Content.WriteRune(c)
		}
	}

	if state.Content.Len() > 0 {
		state.Next(nil)
	}

	res, err := json.Marshal(state.Tags)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(res))
}

type State struct {
	Tags    []*Tag
	Span    *Tag
	Content strings.Builder
}

func NewState() *State {
	return &State{Span: NewSpan()}
}

func (state *State) Next(clr *string) {
	state.Span.Content = []string{state.Content.String()}
	state.Tags = append(state.Tags, state.Span)
	state.Span = NewSpan()
	if clr != nil {
		state.Span.Attributes["color"] = *clr
	}
	state.Content.Reset()
}

type Tag struct {
	Tag        string            `json:"tag"`
	Attributes map[string]string `json:"attributes"`
	Content    []string          `json:"content"`
}

func NewSpan() *Tag {
	return &Tag{Tag: "span", Attributes: make(map[string]string)}
}
