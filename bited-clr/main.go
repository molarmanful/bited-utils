package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/molarmanful/bited-utils"
	"github.com/rivo/tview"
)

func main() {
	if len(os.Args) < 2 {
		panic("missing arg")
	}
	base := os.Args[1]

	state := &State{}

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

var clrMap = map[rune]string{
	'0': "[black]",
	'1': "[maroon]",
	'2': "[green]",
	'3': "[olive]",
	'4': "[navy]",
	'5': "[purple]",
	'6': "[teal]",
	'7': "[silver]",
	'8': "[gray]",
	'9': "[red]",
	'A': "[lime]",
	'B': "[yellow]",
	'C': "[blue]",
	'D': "[fuchsia]",
	'E': "[aqua]",
	'F': "[white]",
	'.': "[-]",
}

type State struct {
	View *tview.TextView
	Txt  [][]rune
	Clr  [][]rune
	ClrF *os.File
	X    int
	Y    int
}

func (state *State) MkView(app *tview.Application) {
	state.View = tview.NewTextView().
		SetDynamicColors(true).SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	state.View.
		SetBackgroundColor(tcell.ColorDefault).
		SetBorder(true).SetBorderColor(tcell.ColorGray).
		SetTitleColor(tcell.ColorGray)

	state.View.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ec := event.Rune()
		if _, ok := clrMap[ec]; ok {
			state.ClrC(ec)
			return nil
		}
		switch event.Rune() {
		case 'h':
			state.MoveC(-1, 0)
			return nil
		case 'j':
			state.MoveC(0, 1)
			return nil
		case 'k':
			state.MoveC(0, -1)
			return nil
		case 'l':
			state.MoveC(1, 0)
			return nil
		}
		return event
	})

	state.Gen()
}

func (state *State) Gen() {
	var res strings.Builder
	var ir strings.Builder
	for y, l := range state.Txt {
		for x, c := range l {
			if ca, ok := state.ClrAt(x, y); ok {
				res.WriteString(tview.Escape(ir.String()))
				ir.Reset()
				res.WriteString(clrMap[ca])
			}
			isC := x == state.X && y == state.Y
			if isC {
				res.WriteString(tview.Escape(ir.String()))
				ir.Reset()
				res.WriteString(`["c"]`)
				res.WriteRune(c)
				res.WriteString(`[""]`)
			} else {
				ir.WriteRune(c)
			}
		}
		ir.WriteRune('\n')
	}

	res.WriteString(tview.Escape(ir.String()))
	state.View.SetText(res.String()).Highlight("c").ScrollToHighlight()
}

func (state *State) MoveC(x int, y int) {
	y0 := state.Y
	state.X = max(0, min(len(state.Txt[y0])-1, state.X+x))
	state.Y = max(0, min(len(state.Txt)-1, state.Y+y))
	state.Gen()
}

func (state *State) ClrC(c rune) {
	state.SetClr(c, state.X, state.Y)
}

func (state *State) ClrAt(x int, y int) (rune, bool) {
	if len(state.Clr) > y && len(state.Clr[y]) > x && state.Clr[y][x] != ' ' {
		return state.Clr[y][x], true
	}
	return rune(0), false
}

func (state *State) SetClr(c rune, x int, y int) {
	if len(state.Clr) <= y {
		state.Clr = append(state.Clr, make([][]rune, y-len(state.Clr)+1)...)
	}
	if len(state.Clr[y]) <= x {
		xs := make([]rune, x-len(state.Clr[y])+1)
		for i := range xs {
			xs[i] = ' '
		}
		state.Clr[y] = append(state.Clr[y], xs...)
	}
	state.Clr[y][x] = c
	state.WriteClr()
	state.Gen()
}

func (state *State) WriteClr() {
	var res strings.Builder
	var ir strings.Builder
	clr := '.'
	for _, l := range state.Clr {
		for _, c := range l {
			if c == clr {
				ir.WriteRune(' ')
				continue
			}
			if _, ok := clrMap[c]; !ok {
				ir.WriteRune(' ')
				continue
			}
			clr = c
			ir.WriteRune(c)
		}
		res.WriteString(strings.TrimRight(ir.String(), " "))
		ir.Reset()
		res.WriteRune('\n')
	}

	err := state.ClrF.Truncate(0)
	bitedutils.Check(err)
	_, err = state.ClrF.Seek(0, 0)
	bitedutils.Check(err)
	_, err = state.ClrF.WriteString(strings.TrimRight(res.String(), "\n"))
	bitedutils.Check(err)
}
