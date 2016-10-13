// Package clix is cui widget library like gtk for your terminal
package clix

import
//	"log"

(
	"flag"
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

var (
	mouse = flag.Bool("mouseui", true, "With mouse support")
)

/*
TODO
channel for menubar to receive updates
fuzzy find
hotkey for menu :
	clix.HotKeys().Set(tcell.Key, clix.MenuItem) {}
	clix.HotKeys().Remove(tcell.Key) error { return nil }

grid buttons
single char buttons
outbar status bar
more widgets
	filepicker
	ascii pic loader

	KEYBOARDLESS KEY ENTRY ( just arrows and mice to enter passwords!?)

*/

// StdOut exits tcell and returns to stdout/stdin terminal
func StdOut(screen tcell.Screen) {
	screen.Clear()
	screen.Fini()
}

// Load Returns new tcell screen
func Load(s tcell.Screen) tcell.Screen {

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
