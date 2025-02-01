package bitedpango

import (
	"bufio"
	"io"

	"barista.run/colors"
	"barista.run/pango"
)

func Pango(txtR io.Reader, clrR io.Reader, clrs []string) *pango.Node {
	clrsMap := make(map[rune]string)
	var ks = []rune("0123456789ABCDEF")
	for i, clr := range clrs[:min(16, len(clrs))] {
		clrsMap[ks[i]] = clr
	}

	state := NewState()
	txtScan := bufio.NewScanner(txtR)
	clrScan := bufio.NewScanner(clrR)
	for txtScan.Scan() {
		txtLine := []rune(txtScan.Text())
		var clrLine []rune
		if clrScan.Scan() {
			clrLine = []rune(clrScan.Text())
		}

		for i, c := range txtLine {
			if len(clrLine) > i {
				k := clrLine[i]
				switch k {
				case '.':
					state.Blank()
				default:
					if clr, ok := clrsMap[k]; ok {
						state.Clr(clr)
					}
				}
			}
			state.Content.WriteRune(c)
		}
		state.Content.WriteRune('\n')
	}
	state.Blank()

	return state.Root
}

func BgFg(node *pango.Node, bg string, fg string) {
	node.Background(colors.Hex(bg)).Color(colors.Hex(fg))
}
