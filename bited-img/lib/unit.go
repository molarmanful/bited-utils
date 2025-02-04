package bitedimg

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitfield/script"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
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

	unit.GensSet = make(map[string]struct{})
	for _, gen := range unit.Gens {
		unit.GensSet[gen.Name] = struct{}{}
	}
	tmpd, err := os.MkdirTemp("", "bited-img-")
	if err != nil {
		return err
	}
	unit.TmpDir = tmpd
	unit.TmpTxtDir = filepath.Join(unit.TmpDir, "txts")
	unit.TmpFontDir = filepath.Join(unit.TmpDir, "fonts")
	unit.Font = filepath.Join(unit.TmpFontDir, "tmp.ttf")
	unit.FC = filepath.Join(unit.TmpDir, "fonts.conf")

	return nil
}
