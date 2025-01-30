package bitedimg

import (
	"text/template"
)

type Unit struct {
	Name  string `koanf:"-"`
	Src   string `koanf:"-"`
	Codes []int  `koanf:"-"`

	SrcForm     SrcForm             `koanf:"src"`
	OutDir      string              `koanf:"out_dir"`
	TxtDir      string              `koanf:"txt_dir"`
	HideAccents bool                `koanf:"hide_accents"`
	FontSize    int                 `koanf:"font_size"`
	Chars       Chars               `koanf:"chars"`
	Map         Map                 `koanf:"map"`
	Clrs        Clrs                `koanf:"clrs"`
	Gens        map[string][]string `koanf:"gens"`
}

var DUnit = Unit{
	SrcForm:     SrcForm{srcT},
	OutDir:      "out",
	TxtDir:      "txt",
	HideAccents: true,
	FontSize:    16,
	Chars: Chars{
		Out:   "chars",
		Width: 48,
	},
	Map: Map{
		Out:       "map",
		LabelClrs: []string{"1", "5"},
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
	Out       string   `koanf:"out"`
	LabelClrs []string `koanf:"label_clrs"`
	BorderClr string   `koanf:"border_clr"`
}

type Clrs struct {
	Bg   string   `koanf:"bg"`
	Fg   string   `koanf:"fg"`
	Base []string `koanf:"base"`
}
