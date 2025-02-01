package bitedbuild

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
	"github.com/molarmanful/bited-utils/bited-scale/lib"
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
	for _, x := range unit.Xs {
		if err := unit.BuildX(x); err != nil {
			return err
		}
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
	log.Println("+ BUILD ttf")
	if out, err := exec.Command(
		"bitsnpicas", "convertbitmap", "-f", "ttf", "-o", unit.TTF, unit.Src).
		CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	log.Println("  + FIX ttf")
	var fixB strings.Builder
	if err := fixTmpl.Execute(&fixB, unit.TTFix); err != nil {
		return err
	}
	if out, err := exec.Command("fontforge", "-c", fixB.String(), unit.TTF).
		CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	log.Println("  + PATCH nerd")
	if unit.Nerd {
		if out, err := exec.Command("nerd-font-patcher", unit.TTF, "-out", unit.OutDir, "--careful", "-c").
			CombinedOutput(); err != nil {
			fmt.Fprintln(os.Stderr, string(out))
			return err
		}
		if out, err := exec.Command("nerd-font-patcher", unit.TTF, "-out", unit.OutDir, "--careful", "-c", "-s").
			CombinedOutput(); err != nil {
			fmt.Fprintln(os.Stderr, string(out))
			return err
		}
	}

	log.Println("+ BUILD woff2")
	if err := exec.Command("woff2_compress", unit.TTF).Run(); err != nil {
		return err
	}

	return nil
}

func (unit *Unit) BuildX(x int) error {
	if x > 1 {
		log.Println("+ SCALE", x)
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

		return unit.BuildBit(src, base, unit.Name)
	}

	log.Println("+ SCALE 1")
	stem, _, _ := strings.Cut(filepath.Base(unit.Src), ".")
	return unit.BuildBit(unit.Src, filepath.Join(unit.OutDir, stem), unit.Name)
}

//go:embed bit.py
var bitPy string
var bitTmpl = template.Must(template.New("").Parse(bitPy))

func (unit *Unit) BuildBit(src string, base string, name string) error {
	log.Println("  + BUILD otb,dfont,pcf")
	var bitB strings.Builder
	if err := bitTmpl.Execute(&bitB, unit.TTFix); err != nil {
		return err
	}
	if out, err := exec.Command("fontforge", "-c", bitB.String(), src, base+".", name).
		CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	log.Println("  + BUILD pcf")
	if out, err := exec.Command("bdftopcf", "-o", base+".pcf", src).
		CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	return nil
}
