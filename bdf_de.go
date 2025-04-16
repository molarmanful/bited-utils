package bitedutils

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/pkg/errors"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

var reKV = regexp.MustCompile(`^\s*(\w+)\s*(.*)\s*$`)

func R2BDF(r io.Reader) (*BDF, error) {
	state := _State{
		BDF: &BDF{
			Props:   orderedmap.New[string, interface{}](),
			Named:   make(map[string]*Glyph),
			Unicode: make(map[rune]*Glyph),
		},
		Defs:      make(map[string]struct{}),
		GlyphDefs: make(map[string]struct{}),
	}

	scan := bufio.NewScanner(r)
	for scan.Scan() {
		if state.Mode == POST {
			break
		}

		line := scan.Text()

		match := reKV.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		state.K = ""
		if len(match) > 1 {
			state.K = match[1]
		}
		if state.K == "COMMENT" {
			continue
		}
		state.V = ""
		if len(match) > 2 {
			state.V = match[2]
		}

		if err := state.Next(); err != nil {
			return nil, errors.WithMessagef(err, "line %d", state.Line)
		}
	}

	if state.Mode != POST {
		return nil, fmt.Errorf("reached end of file without finding ENDFONT")
	}
	return state.BDF, nil
}

type _State struct {
	*BDF
	Line int
	K    string
	V    string
	Mode int
	Defs map[string]struct{}
	*Glyph
	GlyphDefs map[string]struct{}
	Dim       [2]uint64
	Row       int
}

const (
	PRE = iota
	X
	PROPS
	CHARS
	CHAR
	BM
	POST
)

func (state *_State) Next() error {
	state.Line++
	var err error
	switch state.Mode {

	case PRE:
		err = state.ModePRE()
	case X:
		err = state.ModeX()
	case PROPS:
		err = state.ModePROPS()
	case CHARS:
		err = state.ModeCHARS()
	case CHAR:
		err = state.ModeCHAR()
	case BM:
		err = state.ModeBM()

	}
	if err != nil {
		return err
	}
	return nil
}

func (state *_State) ModePRE() error {
	if state.K == "STARTFONT" {
		if err := state.NotDefK(); err != nil {
			return err
		}
		state.Mode = X
	}
	return nil
}

func (state *_State) ModeX() error {
	switch state.K {

	case "FONT":
		if err := state.NotDefK(); err != nil {
			return err
		}
		xlfd, err := ParseXLFD(state.V)
		if err != nil {
			return err
		}
		state.BDF.XLFD = xlfd

	case "STARTPROPERTIES":
		if _, ok := state.Defs["FONT"]; !ok {
			return fmt.Errorf("missing FONT before STARTPROPERTIES")
		}
		state.Mode = PROPS

	case "CHARS":
		if err := state.NotDefK(); err != nil {
			return err
		}
		state.Mode = CHARS

	case "ENDFONT":
		return fmt.Errorf("reached ENDFONT without finding any glyphs")

	case "SIZE",
		"FONTBOUNDINGBOX",
		"CONTENTVERSION",
		"METRICSSET",
		"SWIDTH",
		"SWIDTH1",
		"DWIDTH",
		"DWIDTH1",
		"VVECTOR":
		state.NotDefK()

	default:
		return fmt.Errorf("unknown keyword %s", state.K)
	}
	return nil
}

func (state *_State) ModePROPS() error {
	switch state.K {

	case "ENDPROPERTIES":
		state.Mode = X

	case "FOUNDRY",
		"FAMILY_NAME",
		"WEIGHT_NAME",
		"SLANT",
		"SETWIDTH_NAME",
		"ADD_STYLE_NAME",
		"PIXEL_SIZE",
		"POINT_SIZE",
		"RESOLUTION_X",
		"RESOLUTION_Y",
		"SPACING",
		"AVERAGE_WIDTH",
		"CHARSET_REGISTRY",
		"CHARSET_ENCODING":
		if err := state.NotDefProp(); err != nil {
			return err
		}

	default:
		if err := state.NotDefProp(); err != nil {
			return err
		}
		prop, err := state.FromProp()
		if err != nil {
			return errors.WithMessagef(err, "bad %s value", state.K)
		}
		state.Props.Set(state.K, prop)

	}
	return nil
}

func (state *_State) ModeCHARS() error {
	switch state.K {

	case "STARTCHAR":
		state.Mode = CHAR
		if err := state.NotDefGlyph(); err != nil {
			return err
		}
		state.Glyph = &Glyph{
			Name: state.V,
		}
		clear(state.GlyphDefs)

	case "ENDFONT":
		state.Mode = POST

	default:
		return fmt.Errorf("unknown keyword %s", state.K)
	}
	return nil
}

func (state *_State) ModeCHAR() error {
	switch state.K {
	case "ENCODING":
		if err := state.NotDefGlyphK(); err != nil {
			return err
		}
		n, err := state.V2i()
		if err != nil {
			return errors.WithMessage(err, "ENCODING")
		}
		state.Glyph.Code = n

	case "BBX":
		if err := state.NotDefGlyphK(); err != nil {
			return err
		}
		if err := state.FromBbx(); err != nil {
			return err
		}

	case "DWIDTH":
		if err := state.NotDefGlyphK(); err != nil {
			return err
		}
		n, err := state.V2u()
		if err != nil {
			return errors.WithMessage(err, "DWIDTH")
		}
		state.DWidth = n

	case "BITMAP":
		if err := state.NotDefGlyphK(); err != nil {
			return err
		}
		if state.Glyph.Name == "" {
			return fmt.Errorf("glyph name is empty")
		}
		if _, ok := state.GlyphDefs["ENCODING"]; !ok {
			return fmt.Errorf("missing ENCODING before BITMAP")
		}
		if _, ok := state.GlyphDefs["BBX"]; !ok {
			return fmt.Errorf("missing BBX before BITMAP")
		}
		if _, ok := state.GlyphDefs["DWIDTH"]; !ok {
			return fmt.Errorf("missing DWIDTH before BITMAP")
		}
		if state.Dim[0] > 0 && state.Dim[1] > 0 {
			bm, err := gozxing.NewBitMatrix(int(state.Dim[0]), int(state.Dim[1]))
			if err != nil {
				return errors.WithMessage(err, "BITMAP")
			}
			state.Glyph.Bm = bm
		}
		state.Row = 0
		state.Mode = BM

	case "ENDCHAR":
		return fmt.Errorf("reached ENDCHAR without finding BITMAP")

	case "SWIDTH", "SWIDTH1", "DWIDTH1", "VVECTOR":
		if err := state.NotDefGlyphK(); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown keyword %s", state.K)
	}
	return nil
}

func (state *_State) ModeBM() error {
	switch state.K {

	case "ENDCHAR":
		state.Mode = CHARS
		state.BDF.Glyphs = append(state.BDF.Glyphs, state.Glyph)
		if state.Code < 0 {
			state.BDF.Named[state.Glyph.Name] = state.Glyph
		} else {
			state.BDF.Unicode[rune(state.Glyph.Code)] = state.Glyph
		}

	default:
		if state.V != "" {
			return fmt.Errorf("bad BITMAP row")
		}
		if err := state.Glyph.SetRow(state.Row, state.K); err != nil {
			return errors.WithMessage(err, "BITMAP")
		}
		state.Row++

	}
	return nil
}

func (state *_State) NotDef(k string) error {
	if _, ok := state.Defs[k]; ok {
		return fmt.Errorf("%s is already defined", k)
	}
	state.Defs[k] = struct{}{}
	return nil
}

func (state *_State) NotDefK() error {
	return state.NotDef(state.K)
}

func (state *_State) NotDefProp() error {
	return state.NotDef("prop " + state.K)
}

func (state *_State) NotDefGlyph() error {
	return state.NotDef("glyph " + state.V)
}

func (state *_State) NotDefGlyphK() error {
	if _, ok := state.GlyphDefs[state.K]; ok {
		return fmt.Errorf("%s is already defined in glyph %s", state.K, state.Glyph.Name)
	}
	state.GlyphDefs[state.K] = struct{}{}
	return nil
}

func (state *_State) FromProp() (interface{}, error) {
	if strings.HasPrefix(state.V, `"`) {
		if strings.HasSuffix(state.V, `"`) {
			return nil, fmt.Errorf("string not properly closed")
		}
		s := state.V[1 : len(state.V)-1]
		return strings.ReplaceAll(s, `""`, `"`), nil
	}
	return strconv.Atoi(state.V)
}

func (state *_State) FromBbx() error {
	ss := strings.Fields(state.V)
	if len(ss) < 4 {
		return fmt.Errorf("BBX fields < 4")
	}

	w, err := strconv.ParseUint(ss[0], 10, 64)
	if err != nil {
		return errors.WithMessage(err, "BBX w")
	}
	state.Dim[0] = w

	h, err := strconv.ParseUint(ss[1], 10, 64)
	if err != nil {
		return errors.WithMessage(err, "BBX h")
	}
	state.Dim[1] = h

	x, err := strconv.Atoi(ss[2])
	if err != nil {
		return errors.WithMessage(err, "BBX x")
	}
	state.Glyph.Off[0] = x

	y, err := strconv.Atoi(ss[3])
	if err != nil {
		return errors.WithMessage(err, "BBX y")
	}
	state.Glyph.Off[1] = y

	return nil
}

func (state *_State) V2i() (int, error) {
	s, _, _ := strings.Cut(state.V, " ")
	return strconv.Atoi(s)
}

func (state *_State) V2u() (uint64, error) {
	s, _, _ := strings.Cut(state.V, " ")
	return strconv.ParseUint(s, 10, 64)
}
