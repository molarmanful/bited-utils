package bitedimg

import (
	"text/template"

	"golang.org/x/image/font"
)

// A Unit holds the data necessary to build a single font.
type Unit struct {
	// Src is the processed path to the source bited BDF.
	// Derived from SrcForm.
	Src string `koanf:"-"`
	// Codes is the list of all defined codepoints in the source BDF.
	Codes []int `koanf:"-"`
	// BDF is a representation of the source BDF for use with drawing.
	BDF font.Face `koanf:"-"`
	// Ascent is the font's ascent.
	Ascent int `koanf:"-"`
	// ClrsMap maps CLR codes to hex colors.
	// Derived from Clrs.
	ClrsMap map[rune]string `koanf:"-"`

	// Name is the font's family name.
	Name string `koanf:"name"`
	// SrcForm is the template for Src.
	SrcForm SrcForm `koanf:"src"`
	// OutDir is the output directory for generated images.
	OutDir string `koanf:"out_dir"`
	// OutDir is the source directory for TXT/CLR pairs.
	TxtDir string `koanf:"txt_dir"`
	// PadZWs determines whether to add a zero-width space when generating map.txt
	// and chars.txt.
	PadZWs bool `koanf:"pad_zws"`

	Chars Chars `koanf:"chars"`
	Map   Map   `koanf:"map"`
	Clrs  Clrs  `koanf:"clrs"`
	Gens  []Gen `koanf:"gens"`
}

// SrcT is the default value for [Unit] SrcForm.
var SrcT = template.Must(template.New("").Parse("src/{{ .Name }}.bdf"))

// DUnit specifies default values for [Unit].
var DUnit = Unit{
	SrcForm: SrcForm{SrcT},
	OutDir:  "img",
	TxtDir:  "txt",
	Chars: Chars{
		Out:   "chars",
		Width: 48,
	},
	Map: Map{
		Out:       "map",
		UClr:      "5",
		XClr:      "1",
		BorderClr: "8",
	},
	Clrs: Clrs{
		Bg: "#161616",
		Fg: "#ffffff",
		Base: []string{
			"#222222",
			"#e84f4f",
			"#b7ce42",
			"#fea63c",
			"#66aabb",
			"#b7416e",
			"#6d878d",
			"#dddddd",
			"#666666",
			"#d23d3d",
			"#bde077",
			"#ffe863",
			"#aaccbb",
			"#e16a98",
			"#42717b",
			"#cccccc",
		},
	},
}

// SrcForm is a wrapper type to enable [koanf] to unmarshal strings into
// [template.Template] for [Unit] SrcForm.
type SrcForm struct{ *template.Template }

// A SrcFormPat specifies fields passed to [Unit] SrcForm templates.
type SrcFormPat struct {
	Name string
}

func (src *SrcForm) UnmarshalText(text []byte) error {
	var err error
	src.Template, err = template.New("").Parse(string(text))
	return err
}

// Chars is the subconfig for generating chars.txt.
type Chars struct {
	// Out is the base filename to output to inside TxtDir.
	Out string `koanf:"out"`
	// Width is the maximum number of glyphs to fit per line.
	Width int `koanf:"width"`
}

// Map is the subconfig for generating map.txt.
type Map struct {
	// Map is the base filename to output to inside TxtDir.
	Out       string `koanf:"out"`
	UClr      string `koanf:"u_clr"`
	XClr      string `koanf:"x_clr"`
	BorderClr string `koanf:"border_clr"`
}

// Clrs specifies the colorscheme to use during image generation.
type Clrs struct {
	// Bg is the background color (hex).
	Bg string `koanf:"bg"`
	// Fg is the foreground color (hex).
	Fg string `koanf:"fg"`
	// Base is a list of 16 hex colors corresponding to Base16 colors.
	Base []string `koanf:"base"`
}

// Gen specifies a TXT/CLR pair to generate by combining individual TXT/CLR
// pairs.
type Gen struct {
	// Name is the base filename to output to inside TxtDir.
	Name string `koanf:"name"`
	// Txts is a list of TXT/CLR pair basenames to combine in order.
	Txts []string `koanf:"txts"`
}
