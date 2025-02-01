package bitedimg

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/bitfield/script"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

var reENC = regexp.MustCompile(`^\s*ENCODING\s+[^-]`)

//go:embed fonts.conf
var fontsConf string

func NewUnit(pre map[string]any, name string) (Unit, error) {
	k := koanf.New("")
	k.Load(structs.Provider(DUnit, "koanf"), nil)
	k.Load(confmap.Provider(pre, ""), nil)

	var unit Unit
	if err := k.Unmarshal("", &unit); err != nil {
		return unit, err
	}

	if len(unit.Map.LabelClrs) < 2 {
		return unit, fmt.Errorf("%s.map.label_clrs length < 2", name)
	}
	if len(unit.Clrs.Base) < 16 {
		return unit, fmt.Errorf("%s.clrs.base length < 16", name)
	}

	unit.Name = name

	var srcB strings.Builder
	if err := unit.SrcForm.Template.Execute(&srcB, SrcFormPat{Name: name}); err != nil {
		return unit, err
	}
	unit.Src = srcB.String()

	encs, err := script.File(unit.Src).MatchRegexp(reENC).Column(2).Slice()
	if err != nil {
		return unit, err
	}
	unit.Codes = make([]int, len(encs))
	for i, v := range encs {
		n, err := strconv.Atoi(v)
		if err != nil {
			return unit, err
		}
		unit.Codes[i] = n
	}

	tmpd, err := os.MkdirTemp("", "")
	if err != nil {
		return unit, err
	}
	defer func() {
		if err := os.RemoveAll(tmpd); err != nil {
			log.Println(err)
		}
	}()
	unit.TmpDir = tmpd
	unit.TmpTxtsDir = filepath.Join(tmpd, "txts")
	unit.TTF = filepath.Join(tmpd, "fonts/tmp.ttf")
	unit.FC = filepath.Join(tmpd, "fonts.conf")

	return unit, nil
}
