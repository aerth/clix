package main

import (
	"aerthlib/clix"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func main() {
	initialize() // Flags, Logger

	// clix.EnterPrompt
	_ = clix.EnterPrompt("Please Press Enter or Ctrl+C")
	screen := clix.Load() // Get screen

	// clix.NewMainMenu
	mm := clix.NewMainMenu()
	mm.SetTitle("Hello")
	mm.AddLine("Welcome to Galactica RPG")
	mm.AddLine("Press ENTER to continue")

	mm.Present(screen)

	// clix.NewMenuBar: Define a menu with items.
	m := clix.NewMenuBar()
	m.SetTitle("[Demo Onetime Menu]")
	m.NewItem("Cats")
	m.NewItem("Dogs")
	m.NewItem("Bears")
	m.NewItem("Lions")
	m.NewItem("Birds")
	m.NewItem("Lizards")

	// m.Draw: Draw the actual menu over the screen.
	// This gets called in the m.Wait loop quite often.
	m.Draw(screen, false)

	// m.MainLoop() Blocks. See below for a goroutine implementation.
	userSelection := m.MainLoop(screen)
	screen.Clear()
	screen.Fini() // Exit tcell clean

	for _, v := range fmt.Sprintln("User selected:", userSelection) { // Output selection
		fmt.Print(string(v))
		time.Sleep(20 * time.Millisecond)
	}
	for _, v := range fmt.Sprintln("Now for concurrency, adding a quit button, removing title!") {
		fmt.Print(string(v))
		time.Sleep(20 * time.Millisecond)
	}
	for _, v := range fmt.Sprintln("Press enter to continue.") {
		fmt.Print(string(v))
		time.Sleep(20 * time.Millisecond)
	}

	screen = clix.Load()

	m.NewItem("quit")
	chanbutton := make(chan string, 1)
	var loopnum int
	go func() {
		m.Draw(screen, true)
		for {
			chanbutton <- m.MainLoop(screen)
			if 2 < loopnum && loopnum < 10 {
				_, _ = m.NewItem("button " + strconv.Itoa(loopnum) + string(rune(rand.Intn(4000))))
				// set true to clear for new menu item
				m.Draw(screen, true)
			}
			loopnum++
			m.Draw(screen, false)

		}
	}()

	for {

		// Add another menu
		if loopnum == 2 {
			screen.Clear()
			m2 := clix.NewMenuBar()
			m.SetTitle("Original Menu")
			m2.SetTitle("Second Menu")
			m2.NewItem("Humans")
			m2.NewItem("Dolphins")
			m.AddSibling(m2)

		}

		// Add another menu
		if loopnum == 3 {
			screen.Clear()
			m3 := clix.NewMenuBar()
			m3.SetTitle("Third Menu")
			m3.NewItem("Aliens")
			m.AddSibling(m3)
		}

		m.Draw(screen, false)
		screen.Show()
		s := <-chanbutton
		switch s {
		case "quit":
			fmt.Println("Goodbye!")
			screen.Fini()
			os.Exit(0)
		case "dogs":
			log.Println("woof woof")
		}
		clix.Type(screen, 3, 3, 6, "Coming from chan: "+s)
	}
}

func initialize() {
	flag.Parse()
	logger()
}
