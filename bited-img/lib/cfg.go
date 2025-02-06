package bitedimg

import (
	"text/template"
)

type Unit struct {
	Src        string              `koanf:"-"`
	Codes      []int               `koanf:"-"`
	TmpDir     string              `koanf:"-"`
	TmpTxtDir  string              `koanf:"-"`
	TmpFontDir string              `koanf:"-"`
	Font       string              `koanf:"-"`
	FC         string              `koanf:"-"`
	GensSet    map[string]struct{} `koanf:"-"`

	Name        string  `koanf:"name"`
	SrcForm     SrcForm `koanf:"src"`
	OutDir      string  `koanf:"out_dir"`
	TxtDir      string  `koanf:"txt_dir"`
	HideAccents string  `koanf:"hide_accents"`
	FontSize    int     `koanf:"font_size"`
	Chars       Chars   `koanf:"chars"`
	Map         Map     `koanf:"map"`
	Clrs        Clrs    `koanf:"clrs"`
	Gens        []Gen   `koanf:"gens"`
}

var srcT = template.Must(template.New("").Parse("src/{{ .Name }}.bdf"))

var DUnit = Unit{
	SrcForm:     SrcForm{srcT},
	OutDir:      "img",
	TxtDir:      "txt",
	HideAccents: ".",
	FontSize:    16,
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

type SrcForm struct{ *template.Template }

type SrcFormPat struct {
	Name string
}

func (src *SrcForm) UnmarshalText(text []byte) error {
	var err error
	src.Template, err = template.New("").Parse(string(text))
	return err
}

type Chars struct {
	Out   string `koanf:"out"`
	Width int    `koanf:"width"`
}

type Map struct {
	Out       string `koanf:"out"`
	UClr      string `koanf:"u_clr"`
	XClr      string `koanf:"x_clr"`
	BorderClr string `koanf:"border_clr"`
}

type Clrs struct {
	Bg   string   `koanf:"bg"`
	Fg   string   `koanf:"fg"`
	Base []string `koanf:"base"`
}

type Gen struct {
	Name string   `koanf:"name"`
	Txts []string `koanf:"txts"`
}

type MagickPat struct {
	Pango    string
	Out      string
	FontSize int
	Bg       string
	Fg       string
}
