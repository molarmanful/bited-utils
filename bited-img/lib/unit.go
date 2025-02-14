package bitedimg

import (
	"strings"

	"github.com/bitfield/script"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/zachomedia/go-bdf"
)

// NewUnit creates a full-fledged [Unit] from a single-font config and
// populates it with the necessary data for building.
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

	err := unit.postUnit()
	return unit, err
}

// postUnit populates a newly-unmarshaled [Unit] with the necessary data for
// building.
func (unit *Unit) postUnit() error {
	var srcB strings.Builder
	if err := unit.SrcForm.Template.Execute(&srcB, SrcFormPat{Name: unit.Name}); err != nil {
		return err
	}
	unit.Src = srcB.String()

	bdfB, err := script.File(unit.Src).Bytes()
	if err != nil {
		return err
	}
	bdfP, err := bdf.Parse(bdfB)
	if err != nil {
		return err
	}
	unit.BDF = bdfP.NewFace()
	unit.Ascent = bdfP.Ascent

	unit.Codes = make([]int, 0, len(bdfP.Characters))
	for _, g := range bdfP.Characters {
		if g.Encoding > -1 {
			unit.Codes = append(unit.Codes, int(g.Encoding))
		}
	}

	unit.ClrsMap = make(map[rune]string)
	unit.ClrsMap['.'] = unit.Clrs.Fg
	for i, r := range "0123456789ABCDEF" {
		unit.ClrsMap[r] = unit.Clrs.Base[i]
	}

	return nil
}
