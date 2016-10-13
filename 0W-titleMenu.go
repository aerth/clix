package clix

// WIDGET: TitleMenu

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

// TitleMenu is a simple "press any key to continue" screen to display information
type TitleMenu struct {
	screen   tcell.Screen
	lines    []string
	title    string
	justify  bool
	scroller ScrollFrame
	funcmap  FuncMap
	funcmap2 map[tcell.Key]interface{}
	Output   chan interface{}
	status   statuscode
}

type statuscode int

const (
	offline statuscode = iota
	online
	minimized
	destroyed
)

// FuncMap maps functions to a tcell.Key
type FuncMap map[tcell.Key]func(interface{}) interface{}

// NewTitleMenu creates a TitleMenu
func NewTitleMenu() *TitleMenu {
	mm := new(TitleMenu)
	mm.funcmap = map[tcell.Key]func(interface{}) interface{}{}
	mm.funcmap[tcell.KeyCtrlC] = func(interface{}) interface{} {
		mm.GetScreen().Fini()
		os.Exit(0)
		return nil
	}
	mm.funcmap[tcell.KeyCtrlQ] = func(interface{}) interface{} {
		mm.GetScreen().Fini()
		os.Exit(0)
		return nil
	}
	mm.funcmap[tcell.KeyCtrlH] = func(interface{}) interface{} {
		mm.GetScreen().Fini()
		flag.Usage()
		os.Exit(0)
		return nil
	}

	mm.funcmap[tcell.KeyEnter] = func(interface{}) interface{} {
		return "done"
	}

	return mm
}

//SetTitle to the titlemenu
func (m *TitleMenu) SetTitle(title string) {
	m.title = title
	return
}

// AddLine to the splash menu
func (m *TitleMenu) AddLine(line string) {
	m.lines = append(m.lines, line)
}

// AddLines to the splash menu (convert ascii art with strings.Split(art, "\n"))s
func (m *TitleMenu) AddLines(line []string) {
	m.lines = append(m.lines, line...)
}

// AddFunc to the FuncMap
func (m *TitleMenu) AddFunc(key tcell.Key, f func(interface{}) interface{}, data ...interface{}) {

	m.funcmap[key] = f
	m.funcmap2[key] = data[0]
}

// SetJustifyLeft left
func (m *TitleMenu) SetJustifyLeft() {
	m.justify = true
}

// SetJustifyRight right
func (m *TitleMenu) SetJustifyRight() {
	m.justify = false
}

// Clock for a menu
func (m *TitleMenu) Clock() {
	go func() {
		var px int
		for {
			time.Sleep(1 * time.Second)
			if m.screen == nil {
				continue
			}
			xmax, ymax := m.screen.Size()

			clearline(m.screen, ymax-2)
			Type(m.screen, px-len(time.Now().String()), ymax-2, 4, time.Now().String())
			px--
			if px < 0 {
				px = xmax
			}

			m.screen.Show()

		}
	}()
}

// Present a TitleMenu to the screen.
// It's contents are made from lines.
// m.AddLines(strings.Split(essay, "\n"))
func (m *TitleMenu) Present() interface{} {
	if m.title != "" {

		log.Printf("Welcome to %q.\n", m.title)
	} else {
		log.Println("Welcome!")
	}
	m.Clock()
	if m.screen == nil {
		m.screen = load(nil)
	}
	xmax, ymax := m.screen.Size()
	if m.title != "" {
		Type(m.screen, 0, 0, 6, m.title)
	}

	if m.justify {
		if len(m.lines) != 0 {
			for i, s := range m.lines { // Print each line
				Type(m.screen, 0, ymax-2-len(m.lines)+i, 4, s)
			}
		}
	} else {
		if len(m.lines) != 0 {
			for i, s := range m.lines { // Print each line
				Type(m.screen, xmax-len(s)-1, ymax-2-len(m.lines)+i, 4, s)
			}
		}
	}

	m.screen.Show()
	defer m.screen.Fini()
	return m.loop()
}

// TitleMenu loop
func (m *TitleMenu) loop() interface{} {
	for {
		ev := m.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey: // User pressed a key or keycombo
			if m.funcmap[ev.Key()] != nil {
				log.Println("Found funcmap match for TitleMenu")
				out := m.funcmap[ev.Key()](m.funcmap2[ev.Key()])
				log.Println("interesting", out)

				return out

			}
			m.screen.Sync()
			log.Println("bailing")
			return ev

		}
	}
}
