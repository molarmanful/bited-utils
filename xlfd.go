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
	PxSize   uint64
	Res      [2]uint64
	Spacing  string
	AvgW     uint64
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

	n, err := X2u("px size", xs[7])
	if err != nil {
		return nil, err
	}
	xlfd.PxSize = n

	n, err = X2u("res x", xs[9])
	if err != nil {
		return nil, err
	}
	xlfd.Res[0] = n

	n, err = X2u("res y", xs[10])
	if err != nil {
		return nil, err
	}
	xlfd.Res[1] = n

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
		xlfd.Res[0],
		xlfd.Res[1],
		xlfd.Spacing,
		xlfd.AvgW,
	)
}

func (xlfd *XLFD) Props() *orderedmap.OrderedMap[string, interface{}] {
	return orderedmap.New[string, interface{}](orderedmap.WithInitialData(
		[]orderedmap.Pair[string, interface{}]{
			{Key: "FOUNDRY", Value: xlfd.Foundry},
			{Key: "FAMILY_NAME", Value: xlfd.Family},
			{Key: "WEIGHT_NAME", Value: xlfd.Weight},
			{Key: "SLANT", Value: xlfd.Slant},
			{Key: "SETWIDTH_NAME", Value: xlfd.Setwidth},
			{Key: "ADD_STYLE_NAME", Value: xlfd.AddStyle},
			{Key: "PIXEL_SIZE", Value: xlfd.PxSize},
			{Key: "POINT_SIZE", Value: xlfd.PtSize()},
			{Key: "RESOLUTION_X", Value: xlfd.Res[0]},
			{Key: "RESOLUTION_Y", Value: xlfd.Res[1]},
			{Key: "SPACING", Value: xlfd.Spacing},
			{Key: "AVERAGE_WIDTH", Value: xlfd.AvgW},
			{Key: "CHARSET_REGISTRY", Value: "ISO10646"},
			{Key: "CHARSET_ENCODING", Value: "1"},
		}...,
	))
}

func (xlfd *XLFD) PtSize() uint64 {
	return xlfd.PxSize * 72 / xlfd.Res[1] * 10
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

func X2u(x string, s string) (uint64, error) {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return n, errors.WithMessagef(err, "XLFD %s", x)
	}
	return n, nil
}
