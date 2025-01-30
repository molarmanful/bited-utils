package bitedscale

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

var reKV = regexp.MustCompile(`^\s*(\w+)\s*(.*)\s*$`)

func Scale(r io.Reader, w io.Writer, scale int, name string) error {
	state := NewState(w, scale, name)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if state.Scale <= 1 {
			if _, err := fmt.Fprintln(w, line); err != nil {
				return err
			}
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

		if err := state.Next(); err != nil {
			return err
		}
	}

	return nil
}
