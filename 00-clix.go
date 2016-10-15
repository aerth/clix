// Package clix is cui widget library like gtk for your terminal
/*

With clix library, you can quickly create programs with a command-line user interface.
Your program can accept user input with the available Widget types.

Widgets that are currently active:

	Entry - similar to a shell, accepts a line of user input
	MenuBar - up, down, left, right selection menu
	Scroller - a large buffer of text that is able to be scrolled with PGUP/PGDN
	TitleMenu - a "press any key to continue" style menu

To get started, have a look at the _examples directory, and the 99_test.go file.


Testing


To see what clix looks like, try the test suite.

run:
	go test -v

Select "HUMAN" for the interactive test.


Library

Some goals of the clix library:

	Easy to use, write, read, and learn
	Easy to implement into existing command line programs
	A solid array of widgets, which are components that accept user input
	Not to break your software

Compatibility promise

When clix 1 is released, there will be a compatibility promise
to make sure your app doesn't break.

*/
package clix

import (
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
mutex locks
thread control (minimal num)
ascii pic loader
KEYBOARDLESS KEY ENTRY ( just arrows and mice to enter passwords!?)


PR welcome

*/

// StdOut returns to stdout/stdin terminal, closing screen.
func StdOut(screen tcell.Screen) {
	screen.Clear()
	screen.Fini()
}

// load Returns new tcell screen if necessary
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

// Widget all of them
type Widget interface{}
