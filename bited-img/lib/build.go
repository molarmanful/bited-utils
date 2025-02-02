package bitedimg

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bitfield/script"
	"github.com/molarmanful/bited-utils/bited-pango/lib"
	"golang.org/x/sync/errgroup"
)

func (unit *Unit) Build() error {
	defer func() {
		if err := os.RemoveAll(unit.TmpDir); err != nil {
			log.Println(err)
		}
	}()

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
	if err := unit.Txts(); err != nil {
		return err
	}
	if err := unit.Imgs(); err != nil {
		return err
	}

	return nil
}

//go:embed fonts.conf
var fontsConf string
var fcTmpl = template.Must(template.New("").Parse(fontsConf))

func (unit *Unit) Pre() error {
	for _, v := range []string{unit.TmpFontDir, unit.TmpTxtDir, unit.OutDir} {
		if err := os.MkdirAll(v, os.ModePerm); err != nil {
			return err
		}
	}

	log.Println("+ FONT")
	if out, err := exec.Command(
		"bitsnpicas", "convertbitmap", "-f", "ttf", "-o", unit.Font, unit.Src).
		CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	log.Println("+ FONTCONFIG")
	var fcB strings.Builder
	if err := fcTmpl.Execute(&fcB, unit.TmpDir); err != nil {
		return err
	}
	if _, err := script.Echo(fcB.String()).WriteFile(unit.FC); err != nil {
		return err
	}

	return nil
}

func (unit *Unit) Txts() error {
	log.Println("+ PANGO")
	if err := unit.EachTxt(unit.Pango); err != nil {
		return err
	}

	log.Println("+ GENS")
	for _, gen := range unit.Gens {
		if _, err := script.Slice(gen.Txts).
			FilterLine(func(stem string) string {
				return filepath.Join(unit.TxtDir, stem+".txt")
			}).
			Concat().
			WriteFile(filepath.Join(unit.TxtDir, gen.Name+".txt")); err != nil {
			return err
		}
		if _, err := script.Slice(gen.Txts).
			FilterLine(func(stem string) string {
				return filepath.Join(unit.TmpTxtDir, stem)
			}).
			Concat().
			WriteFile(filepath.Join(unit.TmpTxtDir, gen.Name)); err != nil {
			return err
		}
		log.Println("  +", gen.Name)
	}

	return nil
}

func (unit *Unit) Pango(stem string) error {
	if _, ok := unit.GensSet[stem]; ok {
		return nil
	}

	txtF := script.File(filepath.Join(unit.TxtDir, stem+".txt"))
	clrF := script.File(filepath.Join(unit.TxtDir, stem+".clr"))
	root := bitedpango.Pango(txtF, clrF, unit.Clrs.Base).
		Font(unit.Name).Size(float64(unit.FontSize) * .75)
	bitedpango.BgFg(root, unit.Clrs.Bg, unit.Clrs.Fg)
	script.Echo(root.String()).WriteFile(filepath.Join(unit.TmpTxtDir, stem))

	log.Println("  +", stem)
	return nil
}

func (unit *Unit) Imgs() error {
	log.Println("+ IMGS")
	return unit.EachTxt(unit.Img)
}

//go:embed magick.bash
var magickBash string
var magickTmpl = template.Must(template.New("").Parse(magickBash))

func (unit *Unit) Img(stem string) error {
	var magickCmd strings.Builder
	if err := magickTmpl.Execute(&magickCmd, MagickPat{
		Pango:    filepath.Join(unit.TmpTxtDir, stem),
		Out:      filepath.Join(unit.OutDir, stem),
		FontSize: unit.FontSize,
		Bg:       unit.Clrs.Bg,
		Fg:       unit.Clrs.Fg,
	}); err != nil {
		return err
	}

	cmd := exec.Command("bash")
	cmd.Stdin = strings.NewReader(magickCmd.String())
	cmd.Env = append(os.Environ(), "FONTCONFIG_FILE="+unit.FC)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	log.Println("  +", stem)
	return nil
}

func (unit *Unit) EachTxt(f func(stem string) error) error {
	return script.ListFiles(filepath.Join(unit.TxtDir, "*.txt")).
		Filter(func(r io.Reader, w io.Writer) error {
			scan := bufio.NewScanner(r)
			g, _ := errgroup.WithContext(context.Background())
			for scan.Scan() {
				stem, _, _ := strings.Cut(filepath.Base(scan.Text()), ".")
				g.Go(func() error {
					return f(stem)
				})
			}
			return g.Wait()
		}).
		Wait()
}
