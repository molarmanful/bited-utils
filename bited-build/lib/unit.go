package bitedbuild

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

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
	unit.PostUnit()
	return unit, nil
}

func (unit *Unit) PostUnit() error {
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

	return nil
}
