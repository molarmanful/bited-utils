package bitedbbl

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

const (
	X    = iota // anywhere else not covered by Char
	Prop        // STARTPROPERTIES ... ENDPROPERTIES
	Char        // STARTCHAR ... BITMAP
)

var NerdFont = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x23fb, 0x23fe, 1},
		{0x2b58, 0x2b58, 1},
		{0xe000, 0xe00a, 1},
		{0xe0a0, 0xe0a3, 1},
		{0xe0b0, 0xe0c8, 1},
		{0xe0ca, 0xe0ca, 1},
		{0xe0cc, 0xe0d7, 1},
		{0xe200, 0xe2a9, 1},
		{0xe300, 0xe3e3, 1},
		{0xe5fa, 0xe6b7, 1},
		{0xe700, 0xe8ef, 1},
		{0xea60, 0xec1e, 1},
		{0xed00, 0xefce, 1},
		{0xf000, 0xf381, 1},
		{0xf400, 0xf533, 1},
	},
	R32: []unicode.Range32{
		{0xf0001, 0xf1af0, 1},
	},
}

type _State struct {
	W       io.Writer
	Name    string
	Nerd    bool
	Ceil    int
	Mode    int
	K       string
	V       string
	Char    []string
	Enc     int
	DWIDTH  int
	DWIDTHi int
	BBX     [4]int
	BBXi    int
}

func newState(w io.Writer, name string, nerd bool, ceil bool) *_State {
	var cn int
	if ceil {
		cn = 1
	}
	return &_State{
		W:    w,
		Name: name,
		Nerd: nerd,
		Ceil: cn,
	}
}

func (state *_State) Next() error {
	switch state.Mode {
	case Prop:
		return state.ModeProp()
	case Char:
		return state.ModeChar()
	default:
		return state.ModeX()
	}
}

func (state *_State) ModeX() error {
	if _, err := fmt.Fprint(state.W, state.K); err != nil {
		return err
	}
	switch state.K {

	case "STARTPROPERTIES":
		state.Mode = Prop
		if _, err := fmt.Fprint(state.W, " ", state.V); err != nil {
			return err
		}

	case "STARTCHAR":
		state.Mode = Char
		state.Char = make([]string, 0, 4)
		state.Enc = -1
		state.DWIDTHi = -1
		state.BBXi = -1
		if _, err := fmt.Fprint(state.W, " ", state.V); err != nil {
			return err
		}

	case "FONT":
		if err := state.XLFD(); err != nil {
			return err
		}

	default:
		if state.V != "" {
			if _, err := fmt.Fprint(state.W, " ", state.V); err != nil {
				return err
			}
		}

	}
	if _, err := fmt.Fprintln(state.W); err != nil {
		return err
	}
	return nil
}

func (state *_State) ModeProp() error {
	if _, err := fmt.Fprint(state.W, state.K); err != nil {
		return err
	}
	switch state.K {

	case "ENDPROPERTIES":
		state.Mode = X

	case "FAMILY_NAME":
		if _, err := fmt.Fprint(state.W, ` "`, strings.ReplaceAll(state.Name, `"`, `""`), `"`); err != nil {
			return err
		}

	default:
		if state.V != "" {
			if _, err := fmt.Fprint(state.W, " ", state.V); err != nil {
				return err
			}
		}

	}
	if _, err := fmt.Fprintln(state.W); err != nil {
		return err
	}
	return nil
}

func (state *_State) ModeChar() error {
	var b strings.Builder
	b.WriteString(state.K)
	switch state.K {

	case "BITMAP":
		n := 1
		if state.Nerd && state.Enc >= 0 && unicode.Is(NerdFont, rune(state.Enc)) {
			n = 2
		}
		if state.DWIDTHi >= 0 {
			state.DWIDTH = max(state.DWIDTH*n, state.BBX[0])
			state.Char[state.DWIDTHi] = fmt.Sprintf("DWIDTH %d 0", state.DWIDTH)
		}
		if state.BBXi >= 0 {
			state.BBX[2] = max(0, state.BBX[2]+(state.DWIDTH+state.Ceil)*(n-1)/n/2)
			state.Char[state.BBXi] = fmt.Sprintf(
				"BBX %d %d %d %d",
				state.BBX[0],
				state.BBX[1],
				state.BBX[2],
				state.BBX[3],
			)
		}
		for _, line := range state.Char {
			if _, err := fmt.Fprintln(state.W, line); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(state.W, state.K); err != nil {
			return err
		}
		state.Mode = X

	case "ENCODING":
		n, err := state.Atoi()
		if err != nil {
			return err
		}
		state.Enc = n

	case "DWIDTH":
		state.DWIDTHi = len(state.Char)
		n, err := state.Atoi()
		if err != nil {
			return err
		}
		state.DWIDTH = n

	case "BBX":
		state.BBXi = len(state.Char)
		ns, err := state.Astoi()
		if err != nil {
			return err
		}
		if len(ns) < 4 {
			return fmt.Errorf("BBX entries < 4")
		}
		copy(state.BBX[:], ns)

	}
	b.WriteRune(' ')
	b.WriteString(state.V)
	state.Char = append(state.Char, b.String())
	return nil
}

func (state *_State) Atoi() (int, error) {
	a, _, _ := strings.Cut(state.V, " ")
	return strconv.Atoi(a)
}

func (state *_State) Astoi() ([]int, error) {
	as := strings.Fields(state.V)
	ns := make([]int, len(as))
	for i, a := range as {
		n, err := strconv.Atoi(a)
		if err != nil {
			return ns, err
		}
		ns[i] = n
	}
	return ns, nil
}

func (state *_State) XLFD() error {
	if _, err := fmt.Fprint(state.W, " "); err != nil {
		return err
	}

	xlfd := strings.Split(state.V, "-")
	for i, v := range xlfd {
		if i == 0 {
			continue
		}
		if _, err := fmt.Fprint(state.W, "-"); err != nil {
			return err
		}

		if i == 2 {
			if _, err := fmt.Fprint(state.W, state.Name); err != nil {
				return err
			}
			continue
		}

		if _, err := fmt.Fprint(state.W, v); err != nil {
			return err
		}
	}
	return nil
}
