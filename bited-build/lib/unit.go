package bitedbuild

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

func NewUnit(pre map[string]any, name string, nerd bool) (Unit, error) {
	k := koanf.New("")
	k.Load(structs.Provider(DUnit, "koanf"), nil)
	k.Load(confmap.Provider(pre, ""), nil)

	var unit Unit
	k.Unmarshal("", &unit)

	unit.Name = name
	unit.Nerd = nerd
	unit.TTF = filepath.Join(unit.OutDir, name+".ttf")

	var srcB strings.Builder
	if err := unit.SrcForm.Template.Execute(&srcB, SrcFormPat{Name: name}); err != nil {
		return unit, err
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

	return unit, nil
}
