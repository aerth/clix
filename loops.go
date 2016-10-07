package clix

import (
	"fmt"
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

// Launch a eventhandler (in two go routines)
func (ev *EventHandler) Launch() {

	for i, v := range ev.MenuBars {

		log.Printf("Launching Event Handler %v of %v\n", i+1, len(ev.MenuBars))

		go func(v MenuBar) {
			defer func() {
				v.drawing = true
				v.Draw()
				v.drawing = false
			}()
			for {

				go func() {
					log.Println("EventHandler: Listening to key/mouse events")

					for {
						req := v.screen.PollEvent()
						if req == nil {
							log.Println("EventHandler: screen PollEvent is nil, leaving.")
							return
						}
						log.Println("EventHandler: Sending event to widgetcontroller:" + fmt.Sprintf("%T", req))

						v.widgetcontroller.Input <- req

					}
				}()

				log.Println("Playing Event Handler:", v.title)
				s := v.handleEvents()
				log.Printf("%q returned %q\n", v.title, s)
				if s == "end" {
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

//Quit to all channels
func (ev *EventHandler) Quit() {
	for _, v := range ev.quitchannels {
		v <- 1
	}
}
