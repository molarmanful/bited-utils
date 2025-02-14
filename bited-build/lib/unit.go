package bitedbuild

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitfield/script"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/zachomedia/go-bdf"
)

// NewUnit creates a full-fledged [Unit] from a single-font config and
// populates it with the necessary data for building.
func NewUnit(cfg *koanf.Koanf, nerd bool) (Unit, error) {
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

	unit.Nerd = nerd
	unit.postUnit()
	return unit, nil
}

// postUnit populates a newly-unmarshaled [Unit] with the necessary data for
// building.
func (unit *Unit) postUnit() error {
	unit.TTF = filepath.Join(unit.OutDir, unit.Name+".ttf")

	var srcB strings.Builder
	if err := unit.SrcForm.Template.Execute(&srcB, SrcFormPat{Name: unit.Name}); err != nil {
		return err
	}
	unit.Src = srcB.String()

	var ttfix strings.Builder
	for k, v := range unit.SFNT {
		ttfix.WriteString(fmt.Sprintf("f.appendSFNTName(%q, %q, %q)\n", unit.SFNTLang, k, v))
	}
	ttfix.WriteString(unit.TTFixPre)
	unit.TTFix = ttfix.String()

	xs := map[int]struct{}{1: {}}
	unit.Xs = []int{1}
	for _, x := range unit.XsPre {
		if _, ok := xs[x]; x > 1 && !ok {
			xs[x] = struct{}{}
			unit.Xs = append(unit.Xs, x)
		}
	}

	bdfB, err := script.File(unit.Src).Bytes()
	if err != nil {
		return err
	}
	bdfP, err := bdf.Parse(bdfB)
	if err != nil {
		return err
	}
	unit.FontSize = bdfP.PixelSize

	unit.Widths = make(map[string]int)
	for _, g := range bdfP.Characters {
		name := g.Name
		if g.Encoding >= 0 {
			name = fmt.Sprintf("U+%04X", g.Encoding)
		}
		unit.Widths[name] = g.Advance[0]
	}

	return nil
}
