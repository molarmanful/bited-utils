package bitedimg

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/bitfield/script"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	bitedutils "github.com/molarmanful/bited-utils"
	"github.com/zachomedia/go-bdf"
)

var reENC = regexp.MustCompile(`^\s*ENCODING\s+[^-]`)

func NewUnit(cfg *koanf.Koanf) (Unit, error) {
	var unit Unit
	k := koanf.New("")
	if err := k.Load(structs.Provider(DUnit, "koanf"), nil); err != nil {
		return unit, err
	}
	if err := k.Merge(cfg); err != nil {
		return unit, err
	}
	if err := k.Unmarshal("", &unit); err != nil {
		return unit, err
	}

	err := unit.PostUnit()
	return unit, err
}

func (unit *Unit) PostUnit() error {
	var srcB strings.Builder
	if err := unit.SrcForm.Template.Execute(&srcB, SrcFormPat{Name: unit.Name}); err != nil {
		return err
	}
	unit.Src = srcB.String()

	encs, err := script.File(unit.Src).MatchRegexp(reENC).Column(2).Slice()
	if err != nil {
		return err
	}
	unit.Codes = make([]int, len(encs))
	for i, v := range encs {
		n, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		unit.Codes[i] = n
	}

	fsz, err := bitedutils.GetFsz(unit.Src)
	if err != nil {
		return err
	}
	unit.FontSize = fsz

	bdfB, err := script.File(unit.Src).Bytes()
	if err != nil {
		return err
	}
	bdfP, err := bdf.Parse(bdfB)
	if err != nil {
		return err
	}
	unit.BDF = bdfP.NewFace()

	unit.ClrsMap = make(map[rune]string)
	unit.ClrsMap['.'] = unit.Clrs.Fg
	for i, r := range "0123456789ABCDEF" {
		unit.ClrsMap[r] = unit.Clrs.Base[i]
	}

	return nil
}
