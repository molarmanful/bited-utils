package main

import (
	"bufio"
	"flag"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	bitedutils "github.com/molarmanful/bited-utils"
	bitedimg "github.com/molarmanful/bited-utils/bited-img/lib"
	"github.com/rivo/tview"
)

func main() {
	name := flag.String("name", "", "font name to retrieve colors from")
	stem := flag.String("stem", "", "name of txt/clr pair to retrieve")
	flag.Parse()

	k := koanf.New("")
	err := k.Load(file.Provider("bited-img.toml"), toml.Parser())
	bitedutils.Check(err)

	clrs := bitedimg.DUnit.Clrs
	txtDir := bitedimg.DUnit.TxtDir
	for _, cfg := range k.Slices("fonts") {
		if v := cfg.String("name"); v == *name {
			cfg.Unmarshal("clrs", &clrs)
			cfg.Unmarshal("txt_dir", &txtDir)
			break
		}
	}

	state := &State{}
	state.MkClrs(clrs)
	base := filepath.Join(txtDir, *stem)

	txtF, err := os.Open(base + ".txt")
	bitedutils.Check(err)
	defer txtF.Close()

	txtScan := bufio.NewScanner(txtF)
	for txtScan.Scan() {
		state.Txt = append(state.Txt, []rune(txtScan.Text()))
	}

	clrF, err := os.OpenFile(base+".clr", os.O_CREATE|os.O_RDWR, 0644)
	bitedutils.Check(err)
	defer clrF.Close()
	state.ClrF = clrF

	clrScan := bufio.NewScanner(clrF)
	for clrScan.Scan() {
		state.Clr = append(state.Clr, []rune(clrScan.Text()))
	}

	app := tview.NewApplication()

	state.MkView(app)
	state.View.SetTitle(" COLORING: " + base + " ")

	err = app.SetRoot(state.View, true).SetFocus(state.View).Run()
	bitedutils.Check(err)
}
