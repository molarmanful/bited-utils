package bitedbuild

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
)

var srcT = template.Must(template.New("").Parse("src/{{ .Name }}.bdf"))
var xFormatT = template.Must(template.New("").Parse("{{ .Name }} {{ .X }}x"))

func FullUnit(md toml.MetaData, pv toml.Primitive, name string, nerd bool) error {
	unit := &Unit{
		OutDir:   "out",
		SFNTLang: "English (US)",
		SrcForm:  SrcForm{srcT},
		XForm:    XForm{xFormatT},
	}
	if err := md.PrimitiveDecode(pv, unit); err != nil {
		return err
	}

	unit.Name = name
	unit.Nerd = nerd
	unit.TTF = filepath.Join(unit.OutDir, name+".ttf")

	var srcB strings.Builder
	if err := unit.SrcForm.Template.Execute(&srcB, SrcFormPat{Name: name}); err != nil {
		return err
	}
	unit.Src = srcB.String()

	var ttfix strings.Builder
	for k, v := range unit.SFNT {
		ttfix.WriteString(fmt.Sprintf("f.appendSFNTName(%q, %q, %q)\n", unit.SFNTLang, k, v))
	}
	ttfix.WriteString(unit.TTFixPre)
	unit.TTFix = ttfix.String()

	xs := map[int]bool{1: true}
	unit.Xs = []int{1}
	for _, x := range unit.XsPre {
		if _, ok := xs[x]; x > 1 && !ok {
			xs[x] = true
			unit.Xs = append(unit.Xs, x)
		}
	}

	return unit.Build()
}
