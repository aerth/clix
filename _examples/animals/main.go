package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aerth/clix"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func main() {
	initialize() // Flags, Logger

	// clix.EnterPrompt
	_ = clix.EnterPrompt("Please Press Enter or Ctrl+C")
	events := clix.NewEventHandler()
	// clix.NewTitleMenu
	mm := clix.NewTitleMenu()
	mm.SetTitle("ANIMALS")
	mm.AddLine("Welcome to Animals")
	mm.AddLine("Press ENTER to continue")

	mm.Present()

	// clix.NewMenuBar: Define a menu with items.
	m := clix.NewMenuBar(nil)
	m.SetTitle("[Demo Onetime Menu]")
	m.NewItem("Cats")
	m.NewItem("Dogs")
	m.NewItem("Bears")
	m.NewItem("Lions")
	m.NewItem("Birds")
	m.NewItem("Lizards")

	events.AddMenuBar(m)
	events.Launch()
	var mainout interface{}
	select {
	// case l := <-menu.GetWidgetController().Output:
	// 	log.Println("MC Out:", l)
	case l := <-events.Output:
		log.Println("Output:", l)
		mainout = l
		switch l.(type) {
		case string:
			if l == "quit" {
				goto Done
			}
			if l == "fun" {

			}

		}

	}
Done:
	m.Destroy()
	fmt.Println(mainout)
}

func initialize() {
	flag.Parse()
	logger()
}
