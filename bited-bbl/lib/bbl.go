package bitedbbl

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

var reKV = regexp.MustCompile(`^\s*(\w+)\s*(.*)\s*$`)

// Bbl proportionalizes a bited BDF.
// It reads/writes via streams.
func Bbl(r io.Reader, w io.Writer, name string, nerd bool, ceil bool) error {
	if name == "" {
		return fmt.Errorf("name is empty")
	}

	state := newState(w, name, nerd, ceil)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
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

		if err := state.Next(); err != nil {
			return err
		}
	}

	return nil
}
