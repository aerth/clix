package clix

// WIDGET: MenuBar

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell"
)

const quitkey = "quit"

// ToolBar spans all of x and contains ToolButtons
type ToolBar struct {
	Children []ToolButton
}

var style4 = tcell.Style(5)
var style5 = tcell.Style(4)
var styleSelected = style2 + 8
var styleDefault = style2

// MenuBar is a container for MenuItems
type MenuBar struct {
	Widget
	Selection int
	title     string // Set with m.SetTitle(string)
	zindex    int
	message   string
	//multimessage     []string
	Children         []*MenuItem
	Siblings         []*MenuBar // To have menus in parallel
	screen           tcell.Screen
	scroller         *ScrollFrame
	events           *EventHandler
	widgetcontroller *WidgetController
	entrybar         *Entry
	funcmap          FuncMap
	funcmap2         map[tcell.Key]interface{}
	MaxChars         int
	mostitems        int
	drawing          bool
	menumap          map[string]func() interface{}
}

// MenuItem is a selectable menu item
type MenuItem struct {
	Label string
}

// SetMessage at menbar
func (m *MenuBar) SetMessage(s string) (sel string) {
	m.message = s
	return
}

// Connect a function to a menu item, not that different than AddFunc
func (m *MenuBar) Connect(s string, f func() interface{}) (sel string) {
	if m.menumap == nil {
		m.menumap = map[string]func() interface{}{}
	}
	m.menumap[s] = f
	return
}

//AddEntry to a menubar
func (m *MenuBar) AddEntry(menuLabel string, e *Entry) {
	m.NewItem("Entry" + menuLabel)
	m.Connect("Entry"+menuLabel, func() interface{} {
		m.entrybar.PresentMinimal(m.widgetcontroller.Input, m.events.Output)
		return "gold"
	})

	m.entrybar = e
}

// AddFunc to the FuncMap
func (m *MenuBar) AddFunc(key tcell.Key, f func(interface{}) interface{}, data ...interface{}) {
	m.funcmap[key] = f
	m.funcmap2[key] = data
}

// WidgetController is able to close widgets, and receive widget input.
type WidgetController struct {
	MenuBar       MenuBar
	ScrollFrame   ScrollFrame
	TitleMenu     TitleMenu
	x, y          int
	width, height int
	Input, Output chan interface{}
}

// Present returns a new WidgetController and draws the MenuBar.
func (m *MenuBar) Present(clear bool) {
	if clear {
		m.screen.Clear()
	}
	m.drawing = true
	m.draw()
	m.drawing = false
	return
}

// draw the MenuBar once
func (m *MenuBar) draw() {

	if m.events == nil {
		log.Println("not drawing. nil events")
		return
	}
	if m.screen == nil {
		log.Println("not drawing. nil screen")
		return
	}
	if m.widgetcontroller == nil {
		log.Println("not drawing. nil widgetcontroller")
		return

	}

	log.Println("drawing a MenuBar", m.title)
	if m.screen == nil {
		m.screen = load(m.screen)
	}

	xmax, ymax := m.screen.Size()
	var ts = 0
	if m.title != "" {
		Type(m.screen, 2, ymax-m.mostitems-2, style5, fmt.Sprint(m.title))
		ts++
	}
	if m.message != "" {
		Type(m.screen, 2, 0, style5, fmt.Sprint(m.message))

	}
	var itemnum int
	for i, v := range m.Children {
		itemnum++
		runelabel := []rune(v.Label)
		for r, x, y := 0, 1, ymax-m.mostitems+itemnum-2+ts; r < len(runelabel); r++ {
			if r >= len(runelabel) {
				break
			}
			x++
			if x > xmax {
				y++
				x = 20
			}
			if y > ymax {
				break
			}

			if i == m.Selection && m.zindex == 0 {
				m.screen.SetCell(x, y, styleSelected, rune(runelabel[r]))
			} else {
				m.screen.SetCell(x, y, styleDefault, rune(runelabel[r]))
			}
		}

	}
	if len(m.Siblings) > 0 {
		for i := range m.Siblings {
			m.Siblings[i].drawnextto(m, i)
		}
	}

	if m.scroller != nil {
		m.scroller.Present()
	}

	m.screen.Show()

}

// MENU ITEM

// AddSibling next to a menu
func (m *MenuBar) AddSibling(s *MenuBar) {

	m.Siblings = append(m.Siblings, s)

	m.mostitems = len(m.Children)
	if len(m.Siblings) > 0 {
		for i := range m.Siblings {
			x := len(m.Siblings[i].Children)
			if x > m.mostitems {
				m.mostitems = x
			}
		}
		m.mostitems++
	}

}

// GetScreen ..
func (m *MenuBar) GetScreen() tcell.Screen {
	return m.screen
}

// GetScreen ..
func (m *TitleMenu) GetScreen() tcell.Screen {
	return m.screen
}

// GetScreen ..
func (e *Entry) GetScreen() tcell.Screen {
	return e.screen
}

// Destroy the screen (Must be called before your program exits.)
func (m *MenuBar) Destroy() {
	m.screen.Fini()
}

// Clear the screen
func (m *MenuBar) Clear() {
	m.screen.Clear()
}

// Length ok
func (m *MenuItem) len() int {
	return len(m.Label)
}

// Length ok
func (m *MenuItem) cap() int {
	return 0
}

// Length ok
func (m *MenuItem) empty() bool {
	return m.len() == 0
}

// NewMenuBar returns a new menubar to add things to. If arg is nil, one is created.
func NewMenuBar(s tcell.Screen) *MenuBar {
	m := new(MenuBar)
	if s != nil {
		m.screen = s
	} else {
		m.screen = load(nil)
	}
	wc := new(WidgetController)
	wc.Input = make(chan interface{})
	wc.Output = make(chan interface{})
	m.widgetcontroller = wc
	m.mostitems = len(m.Children)
	m.menumap = map[string]func() interface{}{}
	m.funcmap = map[tcell.Key]func(interface{}) interface{}{}
	m.funcmap2 = map[tcell.Key]interface{}{}
	m.funcmap2 = map[tcell.Key]interface{}{}
	m.Children = []*MenuItem{}
	return m
}

//SetTitle adds a title for the menu
func (m *MenuBar) SetTitle(title string) {
	m.title = title
}

// NewItem Adds a New MenuItem to the MenuBar.
// It returns a pointer to the menuitem.
// Returns a non-nil error only if the label is not unique.
func (m *MenuBar) NewItem(label string) (*MenuItem, error) {

	for i := 0; i < len(m.Children); i++ {
		if m.Children[i].Label == label {

			return m.Children[i], errors.New("NewItem: non-unique label, pointer exists.")
		}
	}

	l := new(MenuItem)
	l.Label = label
	if len(l.Label) > m.MaxChars {
		m.MaxChars = len(l.Label)
	}
	m.Children = append(m.Children, l)
	m.mostitems = len(m.Children)
	return l, nil

}

// GetChild by label name
func (m *MenuBar) GetChild(label string) (*MenuItem, error) {
	for _, v := range m.Children {
		if v.Label == label {
			return v, nil
		}
	}
	return nil, errors.New("menubar: no child by that name")
}

// GetChildNumber by selection number
func (m *MenuBar) GetChildNumber(sel int) (*MenuItem, error) {
	log.Println("Trying", sel, len(m.Children))
	if sel > len(m.Children) {
		return nil, errors.New("Too far")
	}
	for i, v := range m.Children {
		if i == sel-1 {
			return v, nil
		}
	}
	return nil, errors.New("No id")
}

// Len gth ok
func (m *MenuBar) Len() int {
	return len(m.Children)
}

// Length ok
func (m *MenuBar) cap() int {
	return 0
}

// Length ok
func (m *MenuBar) empty() bool {
	return m.Len() == 0
}

// Clock for a menu
func (m *MenuBar) Clock() {}

//GetWidgetController returns the controller for the menubar
// Sending a m.GetWidgetController().Input <- "end" will end the event loop
func (m *MenuBar) GetWidgetController() *WidgetController {
	return m.widgetcontroller
}

/*

already = this menubar is already looping
ended = this menubar has been destroyed so we exited properly
nil event = nil event mystery



*/
func (m *MenuBar) handleEvents() string {
	for {
		xmax, ymax := m.screen.Size()

		// if m.drawing {
		// 	return "already"
		// }

		defer func() {
			m.screen.Clear()
			m.events.Input <- quitkey
			m.drawing = false
		}()

		m.drawing = true
		m.draw()
		// We need a widgetcontroller to continue
		if m.widgetcontroller == nil {
			log.Println("nil WidgetController, waiting one second:", m.title)
			time.Sleep(1 * time.Second)
			continue
		}

		select {
		case <-time.After(10 * time.Second):
			log.Println("Refreshing menu (10 sec)")
			m.draw()
			continue
		case in := <-m.widgetcontroller.Input:
			switch in.(type) {
			case string:
				log.Println("WidgetController: Received signal:", in)

				if in == "continue" {
					continue
				}
				if in == "end" || in == "stop" {
					log.Println("ENDING IT")
					m.screen.Clear()
					m.screen.Show()
					m.events.Output <- "ended"
					return "ended"
				}
				log.Println("")
			case tcell.Event:
				log.Println("WidgetController: Received tcell.Event")

				ev := in.(tcell.Event)

				if ev.When().IsZero() {
					log.Println(" no event? ")
					continue
				}

				switch ev := ev.(type) {

				case *tcell.EventResize:
					log.Println("WidgetController: Event is Resize")

					xmax, ymax = m.screen.Size()
					_ = ymax
					m.screen.Clear()
					m.screen.Show()

				case *tcell.EventKey: // User pressed a key or keycombo
					log.Println("WidgetController: Event is Key")
					if m.entrybar != nil && m.entrybar.isactive {
						m.entrybar.Input <- ev
					}
					switch ev.Key() {
					case tcell.KeyEnd:
						m.scroller.loc = m.scroller.Buffer.Len() - xmax
						m.Clear()
						continue
					case tcell.KeyHome:
						m.scroller.loc = 0
						m.Clear()
						continue
					case tcell.KeyPgDn:

						if m.scroller.loc+xmax+xmax+m.scroller.loc >= m.scroller.Buffer.Len()-5 {
							m.scroller.loc = m.scroller.Buffer.Len() - (xmax * 5)
							if m.scroller.loc < 0 {
								m.scroller.loc = 0
							}

							continue
						}
						m.scroller.loc = m.scroller.loc + (xmax * 5)
						continue
						//log.Println(m.scroller.loc)
					case tcell.KeyPgUp:
						if m.scroller.loc-xmax <= (xmax - 5) {
							m.scroller.loc = 0
							m.Clear()
							continue
						}
						m.scroller.loc = m.scroller.loc - (xmax * 5) - 5
						//log.Println(m.scroller.loc)
					case tcell.KeyEnter:
						if m.zindex != 0 {
							short := m.Siblings[m.zindex-1]
							if 0 > short.Selection || short.Selection > len(short.Children) {
								continue
							}
							s := short.Children[m.Siblings[m.zindex-1].Selection].Label
							m.Siblings[m.zindex-1].Selection = len(m.Siblings[m.zindex-1].Children) + 1
							m.screen.Clear()

							log.Println("Got enter:", s)

							if m.menumap[s] != nil {
								log.Println("ENTER NOT NIL")
								m.menumap[s]()
								continue
							} else {
								log.Println("not in menumap:", m.menumap)
							}
							m.events.Output <- s

							return s
						}
						if len(m.Children) == 0 {
							m.events.Output <- quitkey
							m.events.Input <- quitkey
							return quitkey

						}
						if m.Selection > len(m.Children) || m.Selection < 0 {
							//m.draw(screen, false)// hack
							continue
						}

						if 0 <= m.Selection && m.Selection <= len(m.Children) {
							s := m.Children[m.Selection].Label
							if m.menumap[s] != nil {
								log.Println("ENTER NOT NIL")
								m.menumap[s]()
								continue
							} else {
								log.Println("not in menumap:", m.menumap)
							}
							m.Selection = len(m.Children) - 1
							m.screen.Clear()
							m.events.Output <- s
							return s
							//
						}
						continue
					case tcell.KeyUp:
						if m.zindex != 0 {
							if m.Siblings[m.zindex-1].Selection < 1 {
								m.Siblings[m.zindex-1].Selection = len(m.Siblings[m.zindex-1].Children) - 1
								//m.draw(screen, false)// hack
								continue
							}

							m.Siblings[m.zindex-1].Selection--
							//m.draw(screen, false)// hack
							continue
						}

						if m.Selection < 1 {
							m.Selection = len(m.Children) - 1
							//m.draw(screen, false)// hack
							continue
						}
						m.Selection--
						//m.draw(screen, false)// hack
						continue
					case tcell.KeyDown:

						log.Println(m.zindex)
						if m.zindex != 0 {
							if m.Siblings[m.zindex-1].Selection > len(m.Siblings[m.zindex-1].Children)-2 {
								m.Siblings[m.zindex-1].Selection = 0
								//m.draw(screen, false)// hack
								continue
							}
							m.Siblings[m.zindex-1].Selection++
							//m.draw(screen, false)// hack
							continue
						}
						if m.Selection > len(m.Children)-2 {
							m.Selection = 0
							//m.draw(screen, false)// hack
							continue
						}

						m.Selection++
						//m.draw(screen, false)// hack
						continue
					case tcell.KeyLeft:
						if len(m.Siblings) == 0 {
							log.Println("Left: Continuing because no menu Siblings")
							continue
						}

						if m.zindex <= 0 {
							// Loop around to the left
							m.zindex = len(m.Siblings)
							log.Println("Left: looping right", m.zindex)
							continue
						}

						log.Println("Left", m.zindex, m.Selection, m.Siblings[0].Selection)
						if m.Selection < len(m.Siblings[0].Children) {

							m.Selection = m.Siblings[0].Selection
						} else {
							m.Siblings[0].Selection = len(m.Siblings[0].Children) - 1
						}
						m.zindex--

						continue
					case tcell.KeyRight:
						if len(m.Siblings) == 0 {
							log.Println("Right: Continuing because no menu Siblings")
							continue
						}

						if m.zindex >= len(m.Siblings) {
							// Loop around to the left
							m.zindex = 0
							log.Println("Right: looping left", m.zindex)
							continue
						}

						log.Println("Right", m.zindex, m.Selection, m.Siblings[0].Selection)
						if m.Selection < len(m.Siblings[0].Children) {
							m.Siblings[0].Selection = m.Selection
						} else {
							m.Siblings[0].Selection = len(m.Siblings[0].Children) - 1
						}
						m.zindex++

						continue
					case tcell.KeyEscape, tcell.KeyCtrlQ, tcell.KeyCtrlC:
						m.events.Output <- quitkey
						return quitkey
					}

					// Since the key has not been assigned,
					// lets see if its in the application's defined funcmap

					if m.funcmap[ev.Key()] != nil {
						log.Println("Got a match.", ev.Name())
						var out interface{}
						switch m.funcmap2[ev.Key()].(type) {
						case func():
							out = m.funcmap2[ev.Key()].(func() interface{})()

						}
						m.funcmap[ev.Key()](out)
						continue
					}
					log.Println("Got a key that has not been assigned:", ev.Name())
					//m.events.Input <- ev.Key()
				case *tcell.EventMouse:
					log.Println("WidgetController: Event is Mouse")
					x, y := ev.Position()
					_, ymax := m.screen.Size()
					// switch {
					// case x < 15:
					// 	log.Println("col 1", x, ymax-y)
					// case 16 < x && x < 32:
					// 	log.Println("col 2", x, ymax-y)
					// case 32 < x && x < 60:
					// 	log.Println("col 3", x, ymax-y)
					// default:
					// 	continue
					// }
					switch ev.Buttons() {
					case 4: // right
					//	log.Println("Right Click", x, y)
					case 1: // left
						switch {
						case x < 15:
							//		log.Println("*col 1", x, y)
							if y+2 > ymax-len(m.Children) {
								//		log.Println(len(m.Children) - (ymax - y) + 2)
								form := len(m.Children) - (ymax - y) + 2
								z, err := m.GetChildNumber(form)
								if err != nil {
									//				log.Println(err)
									continue
								}
								//			log.Println("Clicked", z.Label, x, y)

								m.events.Output <- z.Label

							}
						case 16 < x && x < 32:
							//			log.Println("*col 2", x, y)
							if len(m.Siblings) < 1 {
								continue
							}
							form := len(m.Children) - (ymax - y) + 2
							z, err := m.Siblings[0].GetChildNumber(form)
							if err != nil {
								//	log.Println(err)
								continue
							}
							//		log.Println("Clicked", z.Label, x, y)
							m.events.Output <- z.Label

						case 32 < x && x < 60:
							//		log.Println("*col 3", x, y)
							if len(m.Siblings) < 2 {
								continue
							}
							form := len(m.Children) - (ymax - y) + 2
							z, err := m.Siblings[1].GetChildNumber(form)
							if err != nil {
								//			log.Println(err)
								continue
							}
							//		log.Println("Clicked", z.Label, x)
							m.events.Output <- z.Label

						}

					}

					//goto Start
				}
			}
		}
	}
}
