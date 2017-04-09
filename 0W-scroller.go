package clix

// WIDGET: Scroller

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
)

// ScrollFrame has a scroller to read through text.
type ScrollFrame struct {
	Widget
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

// ScrollToByte number
func (s *ScrollFrame) ScrollToByte(n int) {
	s.loc = n
}

// ScrollToByte number
func (s *ScrollFrame) ScrollToEnd() {
	maxx, maxy := s.parentmenu.screen.Size()
	s.loc = s.Buffer.Len() - maxx*maxy
}

//Present to user
func (s *ScrollFrame) Present() {
	b := s.Buffer.Bytes()
	//log.Printf("Scroll Buffer:%vb\n", len(b))

	maxx, maxy := s.parentmenu.screen.Size()
	if s.loc > len(b) {
		s.loc = len(b) - maxx*maxy
	}
	if s.loc <= 0 {
		s.loc = 0
	}
	//offset := 2 * (s.parentmenu.MaxChars * len(s.parentmenu.Children))
	var scrolls [][]rune
	if s.Buffer.Len() < s.loc {
		scrolls = [][]rune{[]rune(s.Buffer.String())}
	} else {
		scrolls, _ = formatrune(maxx-2, maxy, s.Buffer.String()[s.loc:])
	}
	ScrollWriter(s.parentmenu, scrolls)

}

//MAIN MENU

func (s *ScrollFrame) debug() {
	clearchar(s.parentmenu.screen, 1, 1, 2)
	s.parentmenu.screen.Show()

	Type(s.parentmenu.screen, 1, 1, tcell.StyleDefault, strconv.Itoa(s.loc)+"/"+strconv.Itoa(s.Buffer.Len()))
	s.parentmenu.screen.Show()
}

// ToolButton is a small rune button.
// that displays a title in the ToolBar when selected.
type ToolButton struct {
	Label  string
	Icon   rune
	Parent *ToolBar
}

// EnterPrompt returns a single ENTER (Not in a 'screen', just stdin fmt.Println)
func EnterPrompt(prompt ...string) string {
	return EntryOS(prompt...)
}

// EntryOS returns a single ENTER (Not in a 'screen', just stdin fmt.Println)
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
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	return str
}

// AttachScroller truncates old buffer and replaces with s,
// a menubar can have at most one scroller, a scroller can have at most one parent.
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
func ScrollWriter(parent *MenuBar, scrolls [][]rune) {
	scr := parent.screen
	xmax, ymax := scr.Size()

	ts := 0
	m := 1
	if parent.title != "" {
		ts = 4
	}

	if parent.message != "" {
		m = strings.Count(parent.message, "\n") + 1
	}
	x, y := 2, m
	for _, v := range scrolls {
		if y > ymax-parent.mostitems-ts-3 {
			break // thats all for this slice of the buffer
		}
		for i := 2; i < xmax-2; i++ {
			scr.SetCell(i, y+m, style2+4, rune(0))

		}
		x = 2
		for i := 0; i < len(v); i++ {
			//		time.Sleep(2 * time.Millisecond)
			if x >= xmax-8 {
				x = 2
				y++
			}
			x++

			scr.SetCell(x, y+m, style2, v[i])

			if y > ymax-parent.mostitems-ts-3 {
				//	log.Printf("Breaking at %vx%v > %v", y, z, ymax-parent.mostitems-ts-3)
				scr.Show()
				return
			}

		}
		y++
	}
	for i := 2; i < xmax-2; i++ {
		scr.SetCell(i, ymax-parent.mostitems-ts-2, style2+3, tcell.RuneHLine)
	}

	scr.Show()
}

// Format a string into a []string of suitable length lines, using xmax and ymax.
func formatrune(xmax, ymax int, s string) ([][]rune, int) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	// the border
	ymax = ymax - 10

	_, _ = xmax, ymax
	var scanouts [][]rune
	for scanner.Scan() {
		// each line

		scanout := scanner.Text()

		//runes := []rune(scanout)

		var i, x int
		line := []rune(scanout)
		if len(line) <= xmax {
			scanouts = append(scanouts, line)
			continue
		}
		var good []rune
		for _, char := range line {

			// if i >= len(runes) {
			// 	scanouts = append(scanouts, string(line))
			// 	break
			// }
			good = append(good, char)
			if x == xmax {
				x = 0

				scanouts = append(scanouts, good)
				good = nil
			}
			x++
			i++
		}

		//scanouts = append(scanouts, strconv.Itoa(x)+" "+line)
	}

	return scanouts, len(scanouts)
}
