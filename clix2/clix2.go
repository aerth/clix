package clix2

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell"
)

var mouse = flag.Bool("mouse", true, "Use mouse")

func init() {

}

// Window to the next dimension
type Window struct {
	title      string
	children   []Widget
	Input      chan interface{}
	Output     chan interface{}
	screen     tcell.Screen
	buffer     *bytes.Buffer
	Xmax, Ymax int
}

// Widget primarily occupy all the space on the window.
type Widget struct {
	X         int
	Y         int
	MaxWidth  int // -1 for auto
	MaxHeight int // -1 for auto
	Width     int // -1 for auto
	Height    int // -1 for auto
	Content   []Coordinates
	screen    tcell.Screen
}

// Coordinates to draw (or clear)
type Coordinates struct {
	X, Y, Style int
	Rune        rune
	Screen      tcell.Screen
}

// NewWindow first step to app.
func NewWindow(title string) *Window {
	w := new(Window)
	w.screen = load(nil)
	w.Input = make(chan interface{})
	w.Output = make(chan interface{})
	w.title = title
	buf := new(bytes.Buffer)
	w.buffer = buf

	return w
}

// GetScreen method to return screen
func (w *Window) GetScreen() tcell.Screen {
	return w.screen
}

// GetBuffer method to return buffer
func (w *Window) GetBuffer() *bytes.Buffer {
	return w.buffer
}

// Close a window (for good)
func (w *Window) Close() {
	w.screen.Fini()
}

// Clear a window
func (w *Window) Clear() {
	w.screen.Clear()
}

// Draw the rune on coordinates
func (c Coordinates) Draw() {
	c.Screen.SetCell(c.X, c.Y, tcell.Style(c.Style), c.Rune)
	//	s.Show()
}

// Load Returns new tcell screen
func load(s tcell.Screen) tcell.Screen {
	if s != nil {
		if s.Colors() != 0 {
			return s
		}
		s.Fini()
	}
	screen := s
	if s == nil {
		var err error
		screen, err = tcell.NewScreen()
		if err != nil {
			screen.Fini()
			fmt.Fprintf(os.Stderr, "%v\n", err)
			fmt.Fprintf(os.Stdout, "%v\n", err)
			os.Exit(1)
		}
	}
	var err error
	if err = screen.Init(); err != nil {
		screen.Fini()
		fmt.Fprintf(os.Stderr, "%v\n", err)
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(1)
	}
	screen.SetStyle(tcell.StyleDefault.
		Foreground(6).
		Background(0)) // Sweet colors
	if *mouse {
		screen.EnableMouse()
	}
	screen.RegisterRuneFallback(tcell.RuneDegree, "*")
	screen.RegisterRuneFallback(tcell.RuneDiamond, "*")

	return screen
}
func clearlin(s tcell.Screen, y int) {
	mx, _ := s.Size()
	for i := 0; i <= mx; i++ {
		s.SetCell(i, y, tcell.StyleDefault, rune(0))
	}
}
func clearchar(s tcell.Screen, x, x2, y int) {
	mx, _ := s.Size()
	if mx < x2 {
		x2 = mx
	}
	for i := x; i <= x2; i++ {
		s.SetCell(i, y, tcell.StyleDefault, rune(0))
	}
}

// Type a phrase
func (w *Window) Type(coords Coordinates, content string) Coordinates {
	orig := coords
	coords.Screen = w.screen
	lines := strings.Split(content, "\n")
	for _, l := range lines {
		coords.Y++
		coords.X = orig.X
		clearchar(w.screen, orig.X, orig.X+len(l), coords.Y)
		for i, r := range []rune(l) {
			coords.Rune = r
			coords.Draw()
			coords.X++
			w.screen.ShowCursor(coords.X, coords.Y)
			_ = i
		}
	}
	w.screen.ShowCursor(coords.X, coords.Y)
	w.screen.Show()
	return coords
}

// TypeUI a phrase with no cursor
func (w *Window) TypeUI(coords Coordinates, content string) {
	orig := coords
	coords.Screen = w.screen
	for _, l := range strings.Split(content, "\n") {
		clearchar(w.screen, orig.X, orig.X+len(l), coords.Y)
		runes := []rune(l)
		for _, r := range runes {
			coords.Rune = r
			coords.Draw()
			coords.X++
			//w.screen.ShowCursor(coords.X, coords.Y)
			w.screen.Show()

		}
		coords.Y++
		coords.X = orig.X
	}
	w.screen.Show()
}

// Click is a mouse click
type Click struct {
	Button int
	X      int
	Y      int
}

/*
Present the window loop,

**All** input gets relayed to the win.Input channel. It can be ignored, but *can* be used by your program logic.
Typing in the window sends runes to the buffer, redrawing.
Pressing ENTER adds a "\n" to the buffer, and sends the entire buffer to win.Output, redrawing.
Backspace/Delete deletes a byte/rune from the buffer, redrawing.
CtrlC: win.Output channel gets sent "quit" if an exit key combo is pushed. Buffer is available in win.Buffer
*/
func (w *Window) Present() {
	var loopnum int
	for {
		loopnum++
		log.Println("Loop", loopnum)
		e := w.screen.PollEvent()

		switch e.(type) {
		case *tcell.EventMouse:

			b := e.(*tcell.EventMouse).Buttons()
			if b != 0 {
				var click Click
				click.X, click.Y = e.(*tcell.EventMouse).Position()
				click.Button = int(b)
				log.Println(fmt.Sprintf("Mouse button \"%v\" at %v,%v", click.Button, click.X, click.Y))
				go func() {
					w.Input <- click
					log.Println(click, "sent to Input channel")
				}()
			}
			// End case: Mouse
		case *tcell.EventKey:
			log.Println("Key:", e.(*tcell.EventKey).Name())
			switch e.(*tcell.EventKey).Key() {
			case tcell.KeyRune:
				w.buffer.WriteRune(e.(*tcell.EventKey).Rune())
				//w.Type(Coordinates{X: 1, Y: 2}, w.buffer.String())
			case tcell.KeyCtrlJ: // Ctrl ENTER
				mod := e.(*tcell.EventKey).Modifiers()
				log.Println("Key Mod:", mod)
			case tcell.KeyEnter:
				mod := e.(*tcell.EventKey).Modifiers()
				log.Println("Key Mod:", mod)
				w.buffer.WriteString("\n")
				w.Output <- w.buffer.String()
				//w.Type(Coordinates{X: 1, Y: 2}, w.buffer.String())
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				log.Println("Buffer length:", w.buffer.Len())
				if w.buffer.Len() > 0 {
					w.buffer.Truncate(w.buffer.Len() - 1)
					//w.Type(Coordinates{X: 1, Y: 2}, w.buffer.String())
				}
			// case tcell.KeyCtrlS:
			// 	filer.Touch("saveout.log")
			// 	filer.Write("saveout.log", w.buffer.Bytes())
			case tcell.KeyCtrlC, tcell.KeyCtrlQ:
				log.Println("Ctrl+C sending 'quit' to Output channel")
				w.Output <- "quit"
				return
			default:
				go func() {
					w.Input <- e.(*tcell.EventKey).Name()
					log.Println(e.(*tcell.EventKey).Name(), "sent to Input channel")
				}()
			}
			// End case: Key
		case *tcell.EventResize:
			x, y := w.screen.Size()
			w.Xmax, w.Ymax = x, y
			log.Printf("Window resized. New size: %vx%v\n", x, y)
			w.Input <- "resize"
			w.screen.Sync()

		case nil:
			return // Silently exit stage left
		default:
			// We dont know what this was.
			log.Println("Unknown user input.", e)
			w.Output <- "quit"
			return
			// End case: Resize
		}

	}

}
