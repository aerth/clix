package clix

// WIDGET: Entry

import (
	"os"

	"github.com/gdamore/tcell"
)

// Entry is a shell style input bar with a multiline prompt and optional title.
type Entry struct {
	Widget
	title      string
	subtitle   string
	prompt     []string
	IsPassword bool
	screen     tcell.Screen
	scroller   ScrollFrame
	parent     Widget
	maxwidth   int
	Input      chan interface{}
	Output     chan interface{}
	isactive   bool
}

// NewEntry returns a new instance of Entry to be attached to screen s.
func NewEntry(s tcell.Screen) *Entry {
	e := new(Entry)
	e.screen = s
	return e
}

// SetTitle sets a title
func (e *Entry) SetTitle(s string) {
	e.title = s
	return
}

// AttachScroller truncates old buffer and replaces with s
func (e *Entry) AttachScroller(s ScrollFrame) {
	e.scroller.Buffer.Truncate(0)
	e.scroller = s
	e.scroller.parent = e
	return
}

// GetScroller sets a title
func (e *Entry) GetScroller() (s ScrollFrame) {

	return e.scroller
}

// SetSubtitle sets a subtitle
func (e *Entry) SetSubtitle(s string) {
	e.subtitle = s
	return
}

// SetPrompt clears the prompt and replaces with len(s) new lines
func (e *Entry) SetPrompt(s []string) {
	e.prompt = s
	return
}

// SetMaxWidth to an entry
func (e *Entry) SetMaxWidth(width int) {
	e.maxwidth = width
	return
}

// AddPrompt appends len(s) new lines to an Entry window
func (e *Entry) AddPrompt(s ...string) {
	e.prompt = append(e.prompt, s...)
	return
}

// PasswordLabel appends len(s) new lines to an Entry window
func (e *Entry) PasswordLabel(s ...string) {
	e.prompt = append(e.prompt, s...)
	e.IsPassword = true
	return
}

// PresentMinimal to entry while inside a MenuBar (as an "other" field)
func (e *Entry) PresentMinimal(input chan interface{}) string {
	e.Input = input

	e.isactive = true
	for {
		mx, my := e.screen.Size()
		xmax, ymax := e.screen.Size()

		var prompt string
		if e.prompt != nil {
			prompt = e.prompt[0]
		}
		////log.Println("Entry: Drawing prompt rune: >")
		e.screen.SetContent(xmax-21, ymax-1, '>', nil, color1)
		e.screen.Show()
		ev := <-input
		if ev == nil {
			return ""
		}
		//log.Printf("Entry: got input: %T \nl", ev)
		var f func()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			mx, my = e.screen.Size()
			_, ymax = mx, my
			e.screen.Clear()
			f()
			e.screen.Show()
		case *tcell.EventKey: // User pressed a key or keycombo
			switch ev.Key() {
			case tcell.KeyUp, tcell.KeyDown:

			case tcell.KeyCtrlC, tcell.KeyEscape, tcell.KeyCtrlQ:
				e.screen.Fini()
				os.Exit(0)
			case tcell.KeyTAB:
				e.screen.Sync()
				// Eat(e.screen, 0)
				// UnEat(e.screen, 0)
				continue
			case tcell.KeyInsert:

				inserting = !inserting
				drawInputBar(e.screen, 1)
				e.screen.Show()
			case tcell.KeyEnter:

				/* Here we need to clean up before we leave */

				str := getTextMin(e.screen)
				////log.Println(str)
				if str == "" {
					continue
				}
				if e.title != "" {
					offset := 1
					if e.subtitle != "" {
						offset = 2
					}
					clearline(e.screen, ymax-len(prompt)-offset-2)
				}
				if e.subtitle != "" {
					clearline(e.screen, ymax-len(prompt)-1-1)
				}
				if len(prompt) > 0 {
					for i := range prompt {
						clearline(e.screen, my-3-i)
					}
				}
				e.screen.SetContent(xmax-21, ymax-1, ' ', nil, color1)
				e.screen.Show()
				clearline(e.screen, my-1)
				numchar = 0
				selchar = 0
				return str

			case tcell.KeyLeft:
				if selchar < 1 {
					continue
				}
				selchar--
				e.screen.ShowCursor(xmax-20+1+selchar, ymax-1)
				e.screen.Show()
			case tcell.KeyRight:
				if selchar >= numchar {
					continue
				}
				selchar++
				e.screen.ShowCursor(xmax-20+1+selchar, ymax-1)
				e.screen.Show()
			case tcell.KeyDelete:
				if numchar < 1 {
					continue
				}
				ognumchar := numchar
				ogselchar := selchar
				if selchar == numchar {
					continue
				}
				peek := peekInputMin(e.screen)
				////log.Println("Peek:", peek)
				var news []byte
				for i := 0; i < numchar; i++ {
					if i != selchar {
						news = append(news, peek[i])
					}
				}

				////log.Println("testing del:", peek, string(news))
				drawInputBar(e.screen, 1)
				selchar = ogselchar
				numchar = ognumchar - 1
				TypeUI(e.screen, style2, 2, ymax-1, string(news))
				//TypeUI(style, startx, starty, s)
				e.screen.ShowCursor(xmax-20+1+selchar, ymax-1)
			case tcell.KeyBackspace, tcell.KeyClear, 127:
				if numchar == 0 {
					//	drawInputBar(e.screen, 1)
					continue
				}
				if numchar == -1 {
					//	drawInputBar(screen, 1)
					continue
				}
				numchar, selchar = numchar-1, selchar-1
				e.screen.SetContent(xmax-20+1+numchar, ymax-1, rune(0), nil, color1) // Delete Char
				e.screen.ShowCursor(xmax-20+1+numchar, ymax-1)                       // Show Cursor

				e.screen.Show()
				continue
			case tcell.KeyRune:
				xmax, ymax := e.screen.Size()
				if numchar >= xmax-1 {
					// if *verbose {
					// 	////log.Println("Can't type past window border in this version.")
					// }
					continue
				}
				if !inserting || numchar == selchar {
					numchar++
				}
				selchar++
				if ymax-1 >= ymax {
					numchar = 1
				}
				if selchar != numchar {
					if inserting {
						e.screen.SetContent(xmax-20+1+selchar-1, ymax-1, ev.Rune(), nil, color1)
					} else {
						peek := peekInput(e.screen)
						if len(peek) < 3 {

						}
						if e.IsPassword {
							TypeUI(e.screen, color1, xmax-20+1+selchar-1, ymax-1, "*"+peek[selchar-2:])
						} else {
							TypeUI(e.screen, color1, xmax-20+1+selchar-1, ymax-1, string(ev.Rune())+peek[selchar-2:])
						}
						//screen.SetContent(1+selchar-1, ymax-1+line, ev.Rune(), nil, color1)
					}
				} else { // dont matter if inserting or not
					e.screen.SetContent(xmax-20+1+selchar-1, ymax-1, ev.Rune(), nil, color1)
				}
				e.screen.ShowCursor(xmax-20+1+selchar, ymax-1)
				e.screen.Show()
			}
		}
	}
}

// Present an Entry widget
func (e *Entry) Present() string {
	if e.screen == nil {
		e.screen = load(nil)
	} else {
		////log.Println("Not nil screen")
	}
	defer e.screen.Fini()

	prompt := e.prompt

	//Start:
	numchar, selchar = 1, 1
	mx, my := e.screen.Size()
	_, ymax := mx, my

	var f func()
	f = func() {
		if e.title != "" {
			offset := 1
			if e.subtitle != "" {
				offset = 2
			}
			clearline(e.screen, ymax-len(prompt)-offset-2)
			Type(e.screen, mx-len(e.title), ymax-len(prompt)-offset-2, tcell.StyleDefault, e.title)
		}
		if e.subtitle != "" {
			clearline(e.screen, ymax-len(prompt)-1-1)
			Type(e.screen, mx-len(e.title), ymax-len(prompt)-1-1, tcell.StyleDefault, e.subtitle)
		}
		if len(prompt) > 0 {
			for i := range prompt {
				clearline(e.screen, my-2-i)
				Type(e.screen, mx-len(prompt[i]), my-2-i, tcell.StyleDefault, prompt[i])
			}
		}
		numchar, selchar = 0, 0
		drawInputBar(e.screen, 1)
		e.screen.Show()
	}
	f()

	for {

		ev := e.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			mx, my = e.screen.Size()
			_, ymax = mx, my
			e.screen.Clear()
			f()
			e.screen.Show()
		case *tcell.EventKey: // User pressed a key or keycombo
			switch ev.Key() {
			case tcell.KeyUp, tcell.KeyDown:

			case tcell.KeyCtrlC, tcell.KeyEscape, tcell.KeyCtrlQ:
				e.screen.Fini()
				os.Exit(0)
			case tcell.KeyTAB:
				e.screen.Sync()
				// Eat(e.screen, 0)
				// UnEat(e.screen, 0)
				continue
			case tcell.KeyInsert:

				inserting = !inserting
				drawInputBar(e.screen, 1)
				e.screen.Show()
			case tcell.KeyEnter:

				/* Here we need to clean up before we leave */

				str := getText(e.screen)
				if e.title != "" {
					offset := 1
					if e.subtitle != "" {
						offset = 2
					}
					clearline(e.screen, ymax-len(prompt)-offset-2)
				}
				if e.subtitle != "" {
					clearline(e.screen, ymax-len(prompt)-1-1)
				}
				if len(prompt) > 0 {
					for i := range prompt {
						clearline(e.screen, my-3-i)
					}
				}

				drawInputBar(e.screen, 1)

				return str
			case tcell.KeyLeft:
				if selchar < 1 {
					continue
				}
				selchar--
				e.screen.ShowCursor(1+selchar, ymax-1)
				e.screen.Show()
			case tcell.KeyRight:
				if selchar >= numchar {
					continue
				}
				selchar++
				e.screen.ShowCursor(1+selchar, ymax-1)
				e.screen.Show()
			case tcell.KeyDelete:
				if numchar < 1 {
					continue
				}
				ognumchar := numchar
				ogselchar := selchar
				if selchar == numchar {
					continue
				}
				peek := peekInput(e.screen)
				var news []byte
				for i := 0; i < numchar-1; i++ {
					if i != selchar-1 {
						news = append(news, peek[i])
					}
				}

				////log.Println("testing del:", peek, string(news))
				drawInputBar(e.screen, 1)
				selchar = ogselchar
				numchar = ognumchar - 1
				TypeUI(e.screen, style2, 2, ymax-1, string(news))
				//TypeUI(style, startx, starty, s)
				e.screen.ShowCursor(1+selchar, ymax-1)
			case tcell.KeyBackspace, tcell.KeyClear, 127:
				if numchar == 1 {
					//	drawInputBar(e.screen, 1)
					continue
				}
				if numchar == 0 {
					//	drawInputBar(screen, 1)
					continue
				}
				numchar, selchar = numchar-1, selchar-1
				e.screen.SetContent(1+numchar, ymax-1, rune(0), nil, color1) // Delete Char
				e.screen.ShowCursor(1+numchar, ymax-1)                       // Show Cursor

				e.screen.Show()
				continue
			case tcell.KeyRune:
				xmax, ymax := e.screen.Size()
				if numchar >= xmax-1 {
					// if *verbose {
					// 	////log.Println("Can't type past window border in this version.")
					// }
					continue
				}
				if !inserting || numchar == selchar {
					numchar++
				}
				selchar++
				if ymax-1 >= ymax {
					numchar = 1
				}
				if selchar != numchar {
					if inserting {
						e.screen.SetContent(1+selchar-1, ymax-1, ev.Rune(), nil, color1)
					} else {
						peek := peekInput(e.screen)
						if len(peek) < 3 {

						}
						if e.IsPassword {
							TypeUI(e.screen, color1, 1+selchar-1, ymax-1, "*"+peek[selchar-2:])
						} else {
							TypeUI(e.screen, color1, 1+selchar-1, ymax-1, string(ev.Rune())+peek[selchar-2:])
						}
						//screen.SetContent(1+selchar-1, ymax-1+line, ev.Rune(), nil, color1)
					}
				} else { // dont matter if inserting or not
					e.screen.SetContent(1+selchar-1, ymax-1, ev.Rune(), nil, color1)
				}
				e.screen.ShowCursor(1+selchar, ymax-1)
				e.screen.Show()
			}
		}

	}
}
