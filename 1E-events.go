// EVENTS: EventHandler

package clix

import (
	"sync"

	"github.com/gdamore/tcell"
)

// EventHandler to be called in a goroutine before drawing the screen.
/*
The EventHandler is the center of clix.

When the user moves the mouse, or types on the keyboard, a stream of *Events move through the Input chan.

When the menu is finished, a single selection will be sent to the Output chan. It may be "quit", or a menu item.

With the Input events coming in, you can make things happen based on user input.

For example, a game.
	A map shows obstacles and an exit. A click moves the player in that direction. Bad click falls down.
	User presses 30,30, receives a free puppy.
	User presses 29,29, falls in a ditch.
	Meanwhile, goroutines can update the visible cells on the terminal, displaying menus and maps.

Or a tool.
	User is presented with 5 options, chooses the 4th one, labeled "new".
	if out == "new" { menu2.Show() }
*/
type EventHandler struct {
	Input        chan interface{}
	Output       chan interface{}
	MenuItems    []*MenuItem
	MenuBars     []*MenuBar
	EntryItems   []*Entry
	ScrollFrames []*ScrollFrame
	Quitchan     chan int
	screen       tcell.Screen
	quitchannels []chan int
	mu           sync.Mutex
	launched     bool
}

// NewEventHandler returns a new one
func NewEventHandler() *EventHandler {
	ev := new(EventHandler)
	ev.Input = make(chan interface{})
	ev.Output = make(chan interface{})
	ev.Quitchan = make(chan int)

	return ev
}

// Launch fires up the attached widgets.
// This should be one of the first things called in your main,
// after attaching to ONE menu with ev.AddMenuBar
func (ev *EventHandler) Launch() chan int {
	donechan := make(chan int, 1)
	if ev.launched {
		return nil
	}
	ev.launched = true
	// Each MenuBar gets a go func loop
	for i := range ev.MenuBars {
		go func(v *MenuBar) {
			for {
				// This goroutine ONLY sends screen input to the "widgetcontroller.Input" channel
				// It only doesn't hang if your application uses it. It can be ignored.
				// The widgetcontroller Input channel can be listened to with: for { select { case in := <- wc.Input: } }
				go func() {
					for {

						req := v.screen.PollEvent()
						if req == nil {
							// req is only nil if screen ended.
							return
						}
						// Send all input to WidgetController Input
						select {
						default:
							v.widgetcontroller.Input <- req

						}
					}
				}()

				s := v.handleEvents()
				if s == "end" || s == "quit" {
					return
				}
				if s == "already" {
					continue
				}

				event := v.GetScreen().PollEvent()
				if event == nil {
					////log.Println("EventHandler is nil, this goroutine is leaving.")
					v.drawing = false
					return
				}
				////log.Println("Got event, sending to ev.Input")
				ev.Input <- event

			}
		}(ev.MenuBars[i])

	}
	return donechan
}

//AddMenuBar to an ev
func (ev *EventHandler) AddMenuBar(m *MenuBar) {
	ev.mu.Lock()
	defer ev.mu.Unlock()
	m.events = ev
	ev.MenuBars = append(ev.MenuBars, m)

}

//AddScrollFrame to an ev
func (ev *EventHandler) AddScrollFrame(s *ScrollFrame) {
	ev.mu.Lock()
	defer ev.mu.Unlock()
	ev.ScrollFrames = append(ev.ScrollFrames, s)
}

//AddEntry to an ev
func (ev *EventHandler) AddEntry(e *Entry) {
	ev.mu.Lock()
	defer ev.mu.Unlock()
	ev.EntryItems = append(ev.EntryItems, e)
}

//Quit to all channels, return to STDOUT
func (ev *EventHandler) Quit() {
	ev.mu.Lock()
	defer ev.mu.Unlock()
	for _, v := range ev.quitchannels {
		v <- 1
	}
}
