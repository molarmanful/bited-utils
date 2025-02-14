package bitedimg

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitfield/script"
	"golang.org/x/sync/errgroup"
)

// Build generates images from a [Unit] for a single font.
func (unit *Unit) Build() error {
	log.Println("IMGS", unit.Name)

	if err := os.MkdirAll(unit.OutDir, os.ModePerm); err != nil {
		return err
	}
	if err := unit.genChars(); err != nil {
		return err
	}
	if err := unit.genMap(); err != nil {
		return err
	}
	if err := unit.genGens(); err != nil {
		return err
	}
	if err := unit.drawTxts(); err != nil {
		return err
	}

	return nil
}

// genGens combines existing TXT/CLR pairs into new ones.
func (unit *Unit) genGens() error {
	log.Println("+ GEN gens")
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
			Filter(func(r io.Reader, w io.Writer) error {
				scan := bufio.NewScanner(r)
				for scan.Scan() {
					base := filepath.Join(unit.TxtDir, scan.Text())
					txtScan := bufio.NewScanner(script.File(base + ".txt"))
					clrScan := bufio.NewScanner(script.File(base + ".clr"))
					for txtScan.Scan() {
						if clrScan.Scan() {
							fmt.Fprintln(w, clrScan.Text())
						} else {
							fmt.Fprintln(w)
						}
					}
				}
				return nil
			}).
			WriteFile(filepath.Join(unit.TxtDir, gen.Name+".clr")); err != nil {
			return err
		}

		log.Println("  +", gen.Name)
	}
	return nil
}

// drawTxts draws all TXT/CLR pairs in TxtDir to PNG.
func (unit *Unit) drawTxts() error {
	log.Println("+ IMGS")
	return script.ListFiles(filepath.Join(unit.TxtDir, "*.txt")).
		Filter(func(r io.Reader, w io.Writer) error {
			scan := bufio.NewScanner(r)
			g, _ := errgroup.WithContext(context.Background())
			for scan.Scan() {
				stem, _, _ := strings.Cut(filepath.Base(scan.Text()), ".")
				g.Go(func() error {
					if err := unit.drawTCs(stem); err != nil {
						return err
					}
					log.Println("  +", stem)
					return nil
				})
			}
			return g.Wait()
		}).
		Wait()
}
