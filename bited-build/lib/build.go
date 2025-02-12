package bitedbuild

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/bitfield/script"
	bitedscale "github.com/molarmanful/bited-utils/bited-scale/lib"
	"golang.org/x/sync/errgroup"
)

func (unit *Unit) Build() error {
	log.Println("BUILD", unit.Name)

	if err := os.MkdirAll(unit.OutDir, os.ModePerm); err != nil {
		return err
	}
	if err := unit.BuildSrc(); err != nil {
		return err
	}
	if err := unit.BuildVec(); err != nil {
		return err
	}
	if err := unit.BuildXs(); err != nil {
		return err
	}

	return nil
}

func (unit *Unit) BuildSrc() error {
	log.Println("+ COPY src")
	_, err := script.File(unit.Src).WriteFile(filepath.Join(unit.OutDir, filepath.Base(unit.Src)))
	return err
}

//go:embed fix.py
var fixPy string
var fixTmpl = template.Must(template.New("").Parse(fixPy))

func (unit *Unit) BuildVec() error {
	log.Println("+ TTF")
	if out, err := exec.Command(
		"bitsnpicas", "convertbitmap", "-f", "ttf", "-o", unit.TTF, unit.Src).
		CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	log.Println("  + FIX")
	var fixB strings.Builder
	if err := fixTmpl.Execute(&fixB, unit.TTFix); err != nil {
		return err
	}
	if out, err := exec.Command("fontforge", "-c", fixB.String(), unit.TTF, strconv.Itoa(unit.FontSize)).
		CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	if unit.Nerd {
		log.Println("  + NERD")
		g, _ := errgroup.WithContext(context.Background())
		g.Go(func() error {
			if out, err := exec.Command("nerd-font-patcher", unit.TTF, "-out", unit.OutDir, "--careful", "-c").
				CombinedOutput(); err != nil {
				fmt.Fprintln(os.Stderr, string(out))
				return err
			}
			return nil
		})
		g.Go(func() error {
			if out, err := exec.Command("nerd-font-patcher", unit.TTF, "-out", unit.OutDir, "--careful", "-c", "-s").
				CombinedOutput(); err != nil {
				fmt.Fprintln(os.Stderr, string(out))
				return err
			}
			return nil
		})
		return g.Wait()
	}

	log.Println("  + WOFF2")
	if err := exec.Command("woff2_compress", unit.TTF).Run(); err != nil {
		return err
	}

	return nil
}

func (unit *Unit) BuildXs() error {
	log.Println("+ XS")
	g, _ := errgroup.WithContext(context.Background())
	for _, x := range unit.Xs {
		g.Go(func() error {
			return unit.BuildX(x)
		})
	}
	return g.Wait()
}

func (unit *Unit) BuildX(x int) error {
	if x > 1 {
		var nameB strings.Builder
		if err := unit.XForm.Template.Execute(&nameB, XFormPat{Name: unit.Name, X: x}); err != nil {
			return err
		}
		name := nameB.String()
		base := filepath.Join(unit.OutDir, fmt.Sprintf("%s_%dx", unit.Name, x))
		src := base + ".bdf"

		srcF, err := os.Create(src)
		if err != nil {
			return err
		}
		defer srcF.Close()
		if err := bitedscale.Scale(script.File(unit.Src), srcF, x, name); err != nil {
			return err
		}

		return unit.BuildBit(src, x, base, nameB.String())
	}

	stem, _, _ := strings.Cut(filepath.Base(unit.Src), ".")
	return unit.BuildBit(unit.Src, 1, filepath.Join(unit.OutDir, stem), unit.Name)
}

//go:embed bit.py
var bitPy string
var bitTmpl = template.Must(template.New("").Parse(bitPy))

func (unit *Unit) BuildBit(src string, x int, base string, name string) error {
	var bitB strings.Builder
	if err := bitTmpl.Execute(&bitB, unit.TTFix); err != nil {
		return err
	}
	widthsJSON, err := json.Marshal(unit.Widths)
	if err != nil {
		return nil
	}
	if out, err := exec.Command(
		"fontforge", "-c", bitB.String(),
		src, strconv.Itoa(x), strconv.Itoa(unit.FontSize), string(widthsJSON), base+".", name,
	).CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	if out, err := exec.Command("bdftopcf", "-o", base+".pcf", src).
		CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	log.Println("  + DONE", name)
	return nil
}
