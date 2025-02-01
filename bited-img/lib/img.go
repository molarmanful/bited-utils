package bitedimg

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bitfield/script"
	"github.com/molarmanful/bited-utils/bited-pango/lib"
)

func (unit *Unit) Build() error {
	log.Println("IMGS", unit.Name)

	if err := unit.Pre(); err != nil {
		return err
	}
	if err := unit.GenChars(); err != nil {
		return err
	}
	if err := unit.GenMap(); err != nil {
		return err
	}
	if err := unit.CorrectTxts(); err != nil {
		return err
	}

	return nil
}

var fcTmpl = template.Must(template.New("").Parse(fontsConf))

func (unit *Unit) Pre() error {
	for _, v := range []string{"fonts", "txts"} {
		if err := os.MkdirAll(filepath.Join(unit.TmpDir, v), os.ModePerm); err != nil {
			return err
		}
	}

	log.Println("+ TTF")
	if out, err := exec.Command(
		"bitsnpicas", "convertbitmap", "-f", "ttf", "-o", unit.TTF, unit.Src).
		CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	log.Println("+ FONTCONFIG")
	var fcB strings.Builder
	if err := fcTmpl.Execute(&fcB, unit.TmpDir); err != nil {
		return err
	}
	if _, err := script.Echo(fcB.String()).WriteFile(filepath.Join(unit.TmpDir, "fonts.conf")); err != nil {
		return err
	}
	if err := os.Setenv("FONTCONFIG_FILE", unit.FC); err != nil {
		return err
	}

	return nil
}

func (unit *Unit) CorrectTxts() error {
	log.Println("+ TXTCORRECT")
	if err := script.ListFiles(filepath.Join(unit.TxtDir, "*.txt")).
		ExecForEach(`perl -pi -e 'chomp if eof' {{.}}`).
		Wait(); err != nil {
		return err
	}
	return nil
}

var magick = `magick \
  -background "{{ .Bg }}" -fill "{{ .Fg }}" +antialias \
  pango:@"{{ .Pango }}" \
  -bordercolor "{{ .Bg }}" -border "{{ .FontSize }}" \
	"{{ .Out }}.png"`
var magickTmpl = template.Must(template.New("").Parse(magick))

func (unit *Unit) Img(stem string, gen bool) error {
	if _, ok := unit.Gens[stem]; ok != gen {
		return nil
	}

	txtF := script.File(filepath.Join(unit.TxtDir, stem+".txt"))
	clrF := script.File(filepath.Join(unit.TxtDir, stem+".clr"))
	root := bitedpango.Pango(txtF, clrF, unit.Clrs.Base).
		Font(unit.Name).Size(float64(unit.FontSize) * 0.768)
	bitedpango.BgFg(root, unit.Clrs.Bg, unit.Clrs.Fg)
	pango := filepath.Join(unit.TmpTxtsDir, stem)
	script.Echo(root.String()).WriteFile(pango)

	out := filepath.Join(unit.OutDir, stem)
	var magickCmd strings.Builder
	if err := magickTmpl.Execute(&magickCmd, MagickPat{
		Pango:    pango,
		Out:      out,
		FontSize: unit.FontSize,
		Bg:       unit.Clrs.Bg,
		Fg:       unit.Clrs.Fg,
	}); err != nil {
		return err
	}

	cmd := exec.Command("bash")
	cmd.Stdin = strings.NewReader(magickCmd.String())
	cmd.Env = append(os.Environ(), "FONTCONFIG_FILE="+unit.FC)

	log.Println("+ DONE", stem)
	return nil
}
