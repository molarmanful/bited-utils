package bitedbuild

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/bitfield/script"
)

var srcT = template.Must(template.New("").Parse("src/{{ .Name }}.bdf"))
var xFormatT = template.Must(template.New("").Parse("{{ .Name }} {{ .X }}x"))

func FullUnit(md toml.MetaData, pv toml.Primitive, name string, nerd bool) error {
	unit := &Unit{
		OutDir:   "out",
		SFNTLang: "English (US)",
		SrcForm:  SrcForm{srcT},
		XForm:    XForm{xFormatT},
	}
	if err := md.PrimitiveDecode(pv, unit); err != nil {
		return err
	}

	var srcB strings.Builder
	if err := unit.SrcForm.Template.Execute(&srcB, SrcFormPat{Name: name}); err != nil {
		return err
	}

	unit.Name = name
	unit.Src = srcB.String()
	unit.Nerd = nerd
	unit.TTF = filepath.Join(unit.OutDir, name+".ttf")

	var ttfix strings.Builder
	for k, v := range unit.SFNT {
		ttfix.WriteString(fmt.Sprintf("f.appendSFNTName(%q, %q, %q)\n", unit.SFNTLang, k, v))
	}
	ttfix.WriteString(unit.TTFixPre)
	unit.TTFix = ttfix.String()

	xs := map[int]bool{1: true}
	unit.Xs = []int{1}
	for _, x := range unit.XsPre {
		if _, ok := xs[x]; x > 1 && !ok {
			xs[x] = true
			unit.Xs = append(unit.Xs, x)
		}
	}

	return unit.Build()
}

func (unit *Unit) Build() error {
	if err := os.MkdirAll(unit.OutDir, os.ModePerm); err != nil {
		return err
	}

	log.Println("BUILD", unit.Name)

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
	if out, err := exec.Command("bitsnpicas", "convertbitmap", "-f", "ttf", "-o", unit.TTF, unit.Src).CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	log.Println("  + FIX ttf")
	var fixB strings.Builder
	if err := fixTmpl.Execute(&fixB, unit.TTFix); err != nil {
		return err
	}
	if out, err := exec.Command("fontforge", "-c", fixB.String(), unit.TTF).CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	log.Println("  + PATCH nerd")
	if unit.Nerd {
		if out, err := exec.Command("nerd-font-patcher", unit.TTF, "-out", unit.OutDir, "--careful", "-c").CombinedOutput(); err != nil {
			fmt.Fprintln(os.Stderr, string(out))
			return err
		}
		if out, err := exec.Command("nerd-font-patcher", unit.TTF, "-out", unit.OutDir, "--careful", "-c", "-s").CombinedOutput(); err != nil {
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
		scaleCmd := exec.Command("bited-scale", "-n", strconv.Itoa(x), "--name", name)
		scaleCmd.Stdin = script.File(unit.Src)
		scaleCmd.Stdout = srcF
		if err = scaleCmd.Run(); err != nil {
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
	if out, err := exec.Command("fontforge", "-c", bitB.String(), src, base+".", name).CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	log.Println("  + BUILD pcf")
	if out, err := exec.Command("bdftopcf", "-o", base+".pcf", src).CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	return nil
}
