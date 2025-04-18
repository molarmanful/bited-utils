package bitedutils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type XLFD struct {
	Foundry  string
	Family   string
	Weight   string
	Slant    string
	Setwidth string
	AddStyle string
	PxSize   int
	Res      struct {
		X int
		Y int
	}
	Spacing string
	avgW    int
}

func ParseXLFD(s string) (*XLFD, error) {
	xlfd := &XLFD{}

	xs := strings.Split(s, "-")
	if len(xs) == 0 || xs[0] != "" {
		return nil, fmt.Errorf("XLFD doesn't start with '-'")
	}
	if len(xs) < 15 {
		return nil, fmt.Errorf("XLFD fields < 14")
	}

	xlfd.Foundry = xs[1]
	xlfd.Family = xs[2]
	xlfd.Weight = xs[3]

	xlfd.Slant = xs[4]
	if err := xlfd.ValidateSlant(); err != nil {
		return nil, err
	}

	xlfd.Setwidth = xs[5]
	xlfd.AddStyle = xs[6]

	n, err := x2i("px size", xs[7])
	if err != nil {
		return nil, err
	}
	xlfd.PxSize = n

	n, err = x2i("res x", xs[9])
	if err != nil {
		return nil, err
	}
	xlfd.Res.X = n

	n, err = x2i("res y", xs[10])
	if err != nil {
		return nil, err
	}
	xlfd.Res.Y = n

	xlfd.Spacing = xs[11]
	if err := xlfd.ValidateSpacing(); err != nil {
		return nil, err
	}

	return xlfd, nil
}

func (xlfd *XLFD) String() string {
	return fmt.Sprintf(
		"-%s-%s-%s-%s-%s-%s-%d-%d-%d-%d-%s-%d-ISO10646-1",
		xlfd.Foundry,
		xlfd.Family,
		xlfd.Weight,
		xlfd.Slant,
		xlfd.Setwidth,
		xlfd.AddStyle,
		xlfd.PxSize,
		xlfd.PtSize(),
		xlfd.Res.X,
		xlfd.Res.X,
		xlfd.Spacing,
		xlfd.avgW,
	)
}

func (xlfd *XLFD) Props() *orderedmap.OrderedMap[string, any] {
	return orderedmap.New[string, any](orderedmap.WithInitialData(
		[]orderedmap.Pair[string, any]{
			{Key: "FOUNDRY", Value: xlfd.Foundry},
			{Key: "FAMILY_NAME", Value: xlfd.Family},
			{Key: "WEIGHT_NAME", Value: xlfd.Weight},
			{Key: "SLANT", Value: xlfd.Slant},
			{Key: "SETWIDTH_NAME", Value: xlfd.Setwidth},
			{Key: "ADD_STYLE_NAME", Value: xlfd.AddStyle},
			{Key: "PIXEL_SIZE", Value: xlfd.PxSize},
			{Key: "POINT_SIZE", Value: xlfd.PtSize()},
			{Key: "RESOLUTION_X", Value: xlfd.Res.X},
			{Key: "RESOLUTION_Y", Value: xlfd.Res.Y},
			{Key: "SPACING", Value: xlfd.Spacing},
			{Key: "AVERAGE_WIDTH", Value: xlfd.avgW},
			{Key: "CHARSET_REGISTRY", Value: "ISO10646"},
			{Key: "CHARSET_ENCODING", Value: "1"},
		}...,
	))
}

func (xlfd *XLFD) PtSize() int {
	return xlfd.PxSize * 72 / xlfd.Res.Y * 10
}

func (xlfd *XLFD) ValidateSlant() error {
	switch xlfd.Slant {
	case "R", "I", "O", "RI", "RO":
		return nil
	}
	return fmt.Errorf("XLFD slant '%s' is not one of (R, I, O, RI, RO)", xlfd.Slant)
}

func (xlfd *XLFD) ValidateSpacing() error {
	switch xlfd.Spacing {
	case "M", "P", "C":
		return nil
	}
	return fmt.Errorf("XLFD spacing '%s' is not one of (M, P, C)", xlfd.Spacing)
}

func x2i(x string, s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return n, errors.WithMessagef(err, "XLFD %s", x)
	}
	return n, nil
}
