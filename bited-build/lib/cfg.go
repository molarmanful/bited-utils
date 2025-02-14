package bitedbuild

import (
	"text/template"
)

// A Unit holds the data necessary to build a single font.
type Unit struct {
	// Src is the processed path to the source bited BDF.
	// Derived from SrcForm.
	Src string `koanf:"-"`
	// Nerd determines whether to compile Nerd Fonts variants.
	Nerd bool `koanf:"-"`
	// TTF is the path to the output TTF.
	TTF string `koanf:"-"`
	// TTFix is a processed string of custom Python code to be injected when
	// building via Fontforge.
	// Derived from TTFixPre.
	TTFix string `koanf:"-"`
	// Xs is a processed list of integer scales to build.
	// Derived from XsPre.
	Xs []int `koanf:"-"`
	// FontSize is the base font size of the font.
	FontSize int `koanf:"-"`
	// Widths maps glyph names defined in Src to glyph widths.
	Widths map[string]int `koanf:"-"`

	// Name is the font's family name.
	Name string `koanf:"name"`
	// SrcForm is the template for Src.
	SrcForm SrcForm `koanf:"src"`
	// OutDir is the output directory for generated images.
	OutDir string `koanf:"out_dir"`
	// Xs is an unprocessed list of integer scales to build.
	XsPre []int `koanf:"xs"`
	// XForm is the template for each scaled font's name.
	XForm XForm `koanf:"x_format"`
	// SFNTLang is the language to use for SFNT table data.
	SFNTLang string `koanf:"sfnt_lang"`
	// TTFix is an unprocessed string of custom Python code to be injected when
	// building via Fontforge.
	TTFixPre string `koanf:"ttfix"`
	// SFNT represents user-defined key-value pairings for the font's SFNT table.
	SFNT map[string]string `koanf:"sfnt"`
}

// SrcT is the default value for [Unit] SrcForm.
var SrcT = template.Must(template.New("").Parse("src/{{ .Name }}.bdf"))

// XFormT is the default value for [Unit] XForm.
var XFormatT = template.Must(template.New("").Parse("{{ .Name }}{{ .X }}x"))

// DUnit specifies default values for [Unit].
var DUnit = Unit{
	SrcForm:  SrcForm{SrcT},
	OutDir:   "out",
	XForm:    XForm{XFormatT},
	SFNTLang: "English (US)",
}

// SrcForm is a wrapper type to enable [koanf] to unmarshal strings into
// [template.Template] for [Unit] SrcForm.
type SrcForm struct{ *template.Template }

// A SrcFormPat specifies fields passed to [Unit] SrcForm templates.
type SrcFormPat struct {
	// Name is the font's family name.
	Name string
}

func (src *SrcForm) UnmarshalText(text []byte) error {
	var err error
	src.Template, err = template.New("").Parse(string(text))
	return err
}

// XForm is a wrapper type to enable [koanf] to unmarshal strings into
// [template.Template] for [Unit] XForm.
type XForm struct{ *template.Template }

// XFormPat specifies fields passed to [Unit] XForm templates.
type XFormPat struct {
	// Name is the font's family name.
	Name string
	// X is the scale.
	X int
}

func (xf *XForm) UnmarshalText(text []byte) error {
	var err error
	xf.Template, err = template.New("").Parse(string(text))
	return err
}
