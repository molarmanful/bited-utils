package bitedbuild

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitfield/script"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	bitedutils "github.com/molarmanful/bited-utils"
)

var reCHAR = regexp.MustCompile(`^\s*(STARTCHAR|ENCODING|DWIDTH)\s+`)

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

	fsz, err := bitedutils.GetFsz(unit.Src)
	if err != nil {
		return err
	}
	unit.FontSize = fsz

	chars, err := script.File(unit.Src).MatchRegexp(reCHAR).Slice()
	if err != nil {
		return err
	}
	unit.Widths = make(map[string]int)
	name := ""
	for _, line := range chars {
		kv := strings.Fields(line)
		if len(kv) < 2 {
			return fmt.Errorf("invalid BDF line: %q", line)
		}
		switch kv[0] {
		case "STARTCHAR":
			name = kv[1]
		case "ENCODING":
			uc, err := strconv.Atoi(kv[1])
			if err != nil {
				return err
			}
			name = fmt.Sprintf("U+%04X", uc)
		case "DWIDTH":
			w, err := strconv.Atoi(kv[1])
			if err != nil {
				return err
			}
			unit.Widths[name] = w
			name = ""
		}
	}

	return nil
}
