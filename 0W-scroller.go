package clix

// WIDGET: Scroller

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
)

// ScrollFrame has a scroller to read through text.
type ScrollFrame struct {
	title       string
	Buffer      *bytes.Buffer
	loc         int // since buffer doesn't know about screen size this is char num
	parent      interface{}
	parentmenu  *MenuBar
	parententry *Entry
}

//NewScrollFrame returns a new one
func NewScrollFrame(t string) *ScrollFrame {
	s := new(ScrollFrame)
	var b bytes.Buffer
	s.Buffer = &b
	s.loc = 0
	s.title = t
	return s
}

//Present to user
func (s *ScrollFrame) Present() {
	b := s.Buffer.Bytes()
	log.Printf("Scroll Buffer:%vb\n", len(b))

	maxx, maxy := s.parentmenu.screen.Size()
	if s.loc > len(b) {
		s.loc = len(b) - maxx*maxy
	}
	if s.loc <= 0 {
		s.loc = 0
	}
	offset := 2 * (s.parentmenu.MaxChars * len(s.parentmenu.Children))
	scrolls, n := formatrune(offset, maxy, s.Buffer.String()[s.loc:])
	ScrollWriter(s.parentmenu, scrolls, n)

}

//MAIN MENU

func (s *ScrollFrame) debug() {
	clearchar(s.parentmenu.screen, 1, 1, 5)
	s.parentmenu.screen.Show()

	Type(s.parentmenu.screen, 1, 1, 2, strconv.Itoa(s.loc)+"/"+strconv.Itoa(s.Buffer.Len()))
	s.parentmenu.screen.Show()
}

// ToolButton is a small rune button.
// that displays a title in the ToolBar when selected.
type ToolButton struct {
	Label  string
	Icon   rune
	Parent *ToolBar
}

// EnterPrompt returns a single ENTER (Not in a screen, just stdout fmt.Println)
func EnterPrompt(prompt ...string) string {
	return EntryOS(prompt...)
}

// EntryOS lol
func EntryOS(prompt ...string) string {
	if len(prompt) > 0 {
		fmt.Println(prompt)
	}
	var str string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return str
}

// AttachScroller truncates old buffer and replaces with s
func (m *MenuBar) AttachScroller(s *ScrollFrame) {

	m.scroller = s
	m.scroller.parent = "MenuBar"
	m.scroller.parentmenu = m
	return
}

// GetScroller sets a title
func (m *MenuBar) GetScroller() (s *ScrollFrame) {

	return m.scroller
}

// ScrollWriter to a screen, all lines with wrap
func ScrollWriter(parent *MenuBar, scrolls []string, chars int) {
	scr := parent.screen
	xmax, ymax := scr.Size()

	ts := 0
	m := 1
	if parent.title != "" {
		ts++
	}

	x, y := 2, m
	for z, v := range scrolls {
		//log.Println("clearing", z)
		if z > ymax-parent.mostitems-ts-3 {
			break
		}
		for i := 2; i < xmax-2; i++ {
			scr.SetCell(i, z+m, style2+3, rune(0))
		}
		x, y = 2, m
		for i := 0; i < len(v); i++ {
			//		time.Sleep(2 * time.Millisecond)
			x++
			if x >= xmax-2 {
				x = 2 * (parent.MaxChars * len(parent.Children))
				y++
			}

			scr.SetCell(x, y+z-1+m, style2+3, rune(v[i]))

			if y+z > ymax-parent.mostitems-ts-3 {
				//	log.Printf("Breaking at %vx%v > %v", y, z, ymax-parent.mostitems-ts-3)
				scr.Show()
				return
			}
		}
	}
	for i := 2; i < xmax-2; i++ {
		scr.SetCell(i, ymax-parent.mostitems-ts-2, style2+3, tcell.RuneHLine)
	}

	scr.Show()
}

// ScrollWriterV to a screen, all lines with wrap
func ScrollWriterV(parent *MenuBar, scrolls []string, chars int) {
	scr := parent.screen
	xmax, ymax := scr.Size()

	style := tcell.StyleDefault + 1
	x, y := 2*(parent.MaxChars*len(parent.Children)), 2

	for i := 0; i < xmax*ymax; i++ {
		//		time.Sleep(2 * time.Millisecond)
		x++
		if x >= xmax-2 {
			x = 2 * (parent.MaxChars * len(parent.Children))
			y++
		}
		if y > ymax-2 {
			y = 2
			x = 2 * (parent.MaxChars * len(parent.Children))
		}
		if i > chars {
			break
		}
		scr.SetCell(x, y, style, rune(0))
	}

	scr.Show()

	for z, v := range scrolls {

		x, y = 2*(parent.MaxChars*len(parent.Children)), 0
		for i := 0; i < len(v); i++ {
			//		time.Sleep(2 * time.Millisecond)
			x++
			if x >= xmax-2 {
				x = 2 * (parent.MaxChars * len(parent.Children))
				y++
			}
			if y > ymax-len(parent.Children) {
				x = 2 * (parent.MaxChars * len(parent.Children))
				break
			}
			// if i >  {
			// 	break
			// }
			scr.SetCell(x, y+z, style3, rune(v[i]))
		}
	}
	scr.Show()
}
func formatrune(xmax, ymax int, s string) ([]string, int) {

	scanner := bufio.NewScanner(strings.NewReader(s))
	// the border
	ymax = ymax - 10
	xmax = xmax - 4
	_, _ = xmax, ymax
	var scanouts []string
	var total int
	for scanner.Scan() {
		// each line
		scanout := scanner.Text()
		total += len(scanout)
		scanouts = append(scanouts, scanout)

	}

	return scanouts, total
}
