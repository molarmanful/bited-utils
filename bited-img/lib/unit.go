package bitedimg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/bitfield/script"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

var srcT = template.Must(template.New("").Parse("src/{{ .Name }}.bdf"))
var reENC = regexp.MustCompile(`^\s*ENCODING\s+[^-]`)

func FullUnit(pre map[string]any, name string) error {
	k := koanf.New("")
	k.Load(structs.Provider(DUnit, "koanf"), nil)
	k.Load(confmap.Provider(pre, ""), nil)

	var unit Unit
	k.Unmarshal("", &unit)

	if len(unit.Map.LabelClrs) < 2 {
		return fmt.Errorf("%s.map.label_clrs length < 2", name)
	}
	if len(unit.Clrs.Base) < 16 {
		return fmt.Errorf("%s.clrs.base length < 16", name)
	}

	unit.Name = name

	var srcB strings.Builder
	if err := unit.SrcForm.Template.Execute(&srcB, SrcFormPat{Name: name}); err != nil {
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

	return unit.Img()
}
