// EVENTS: EventHandler

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

package clix

import (
	"log"

	"github.com/gdamore/tcell"
)

// EventHandler to be called in a goroutine before drawing the screen.
type EventHandler struct {
	Input        chan interface{}
	Output       chan interface{}
	MenuItems    []MenuItem
	MenuBars     []MenuBar
	EntryItems   []Entry
	ScrollFrames []ScrollFrame
	Quitchan     chan int
	screen       tcell.Screen
	quitchannels []chan int
}

// NewEventHandler returns a new one
func NewEventHandler() *EventHandler {
	ev := new(EventHandler)
	ev.Input = make(chan interface{})
	ev.Output = make(chan interface{})
	ev.Quitchan = make(chan int)

	return ev
}

// Launch starts the attached widgets.
func (ev *EventHandler) Launch() {
	for _, v := range ev.MenuBars {

		//log.Printf("Launching Event Handler %v of %v\n", i+1, len(ev.MenuBars))

		go func(v MenuBar) {
			defer func() {
				v.drawing = true
				v.draw()
				v.drawing = false
			}()
			for {

				go func() {
					//	log.Println("EventHandler: Listening to key/mouse events")

					for {
						req := v.screen.PollEvent()
						if req == nil {
							//	log.Println("EventHandler: screen PollEvent is nil, leaving.")
							return
						}
						//						log.Println("EventHandler: Sending event to widgetcontroller:" + fmt.Sprintf("%T", req))

						v.widgetcontroller.Input <- req

					}
				}()

				log.Println("Playing Event Handler:", v.title)
				s := v.handleEvents()
				log.Printf("%q returned %q\n", v.title, s)
				if s == "end" || s == "quit" {
					log.Println("EventHandler got END")
					return
				}
				if s == "already" {

					continue
					//					return
				}

				//v.screen.Show()
				event := v.GetScreen().PollEvent()
				if event == nil {
					log.Println("EventHandler is nil, this goroutine is leaving.")
					v.drawing = false
					return
				}
				log.Println("Got event, sending to ev.Input")
				ev.Input <- event

			}
		}(v)

	}

}

//AddMenuBar to an ev
func (ev *EventHandler) AddMenuBar(m *MenuBar) {
	m.events = ev
	ev.MenuBars = append(ev.MenuBars, *m)
}

//AddScrollFrame to an ev
func (ev *EventHandler) AddScrollFrame(s *ScrollFrame) {
	ev.ScrollFrames = append(ev.ScrollFrames, *s)
}

//AddEntry to an ev
func (ev *EventHandler) AddEntry(e *Entry) {
	ev.EntryItems = append(ev.EntryItems, *e)
}

//Quit to all channels, return to STDOUT
func (ev *EventHandler) Quit() {
	for _, v := range ev.quitchannels {
		v <- 1
	}
}
