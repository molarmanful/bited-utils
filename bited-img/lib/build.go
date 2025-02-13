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
	if err := unit.GenGens(); err != nil {
		return err
	}
	if err := unit.DrawTxts(); err != nil {
		return err
	}

	return nil
}

func (unit *Unit) Pre() error {
	return os.MkdirAll(unit.OutDir, os.ModePerm)
}

func (unit *Unit) GenGens() error {
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

func (unit *Unit) DrawTxts() error {
	log.Println("+ IMGS")
	return script.ListFiles(filepath.Join(unit.TxtDir, "*.txt")).
		Filter(func(r io.Reader, w io.Writer) error {
			scan := bufio.NewScanner(r)
			g, _ := errgroup.WithContext(context.Background())
			for scan.Scan() {
				stem, _, _ := strings.Cut(filepath.Base(scan.Text()), ".")
				g.Go(func() error {
					if err := unit.DrawTCs(stem); err != nil {
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
