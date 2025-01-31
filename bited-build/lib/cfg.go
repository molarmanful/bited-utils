package bitedbuild

import (
	"text/template"
)

type Unit struct {
	Name  string `koanf:"-"`
	Src   string `koanf:"-"`
	Nerd  bool   `koanf:"-"`
	TTF   string `koanf:"-"`
	TTFix string `koanf:"-"`
	Xs    []int  `koanf:"-"`

	SrcForm  SrcForm           `koanf:"src"`
	OutDir   string            `koanf:"out_dir"`
	XsPre    []int             `koanf:"xs"`
	XForm    XForm             `koanf:"x_format"`
	SFNTLang string            `koanf:"sfnt_lang"`
	TTFixPre string            `koanf:"ttfix"`
	SFNT     map[string]string `koanf:"sfnt"`
}

var srcT = template.Must(template.New("").Parse("src/{{ .Name }}.bdf"))
var xFormatT = template.Must(template.New("").Parse("{{ .Name }} {{ .X }}x"))

var DUnit = Unit{
	SrcForm:  SrcForm{srcT},
	OutDir:   "out",
	XForm:    XForm{xFormatT},
	SFNTLang: "English (US)",
}

type SrcForm struct{ *template.Template }

type SrcFormPat struct {
	Name string
}

func (src *SrcForm) UnmarshalText(text []byte) error {
	var err error
	src.Template, err = template.New("").Parse(string(text))
	return err
}

type XForm struct{ *template.Template }

type XFormPat struct {
	Name string
	X    int
}

func (xf *XForm) UnmarshalText(text []byte) error {
	var err error
	xf.Template, err = template.New("").Parse(string(text))
	return err
}
