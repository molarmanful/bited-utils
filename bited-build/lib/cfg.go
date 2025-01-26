package bitedbuild

import (
	"text/template"
)

type Unit struct {
	Name  string
	Src   string
	Nerd  bool
	TTF   string
	TTFix string
	Xs    []int

	SrcForm  SrcForm           `toml:"src"`
	OutDir   string            `toml:"out_dir"`
	XsPre    []int             `toml:"xs"`
	XForm    XForm             `toml:"x_format"`
	SFNTLang string            `toml:"sfnt_lang"`
	TTFixPre string            `toml:"ttfix"`
	SFNT     map[string]string `toml:"sfnt"`
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
