package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aerthlib/clix/clix2"
	"github.com/aerthlib/filer"
)

func main() {
	// Verbose logging to ./debug.log (log.go)
	flag.Parse()
	logger()
	logfile = "debug.log"

	// This presents a new window and makes available two channels.
	// In this example we are only using the Output channel, below.
	win := clix2.NewWindow("Hello")
	mx, my := win.GetScreen().Size()
	/*
		   Debug auto quitter.
		   This exits cleanly after 10 seconds.
			 OK to remove.
	*/

	// DEBUG: Countdown and quit
	go func() {
		var countdown = 60
		c2 := clix2.Coordinates{X: 1, Y: my - 2}
		for {
			win.TypeUI(clix2.Coordinates{X: mx - 2, Y: my - 1}, strconv.Itoa(countdown))
			if countdown == 3 {
				win.TypeUI(c2, "Please insert more credits.")
			}
			if countdown == 1 {
				win.TypeUI(c2, "Demo is over.")
			}
			countdown--
			time.Sleep(1 * time.Second)
			if countdown < 0 {
				win.Output <- "quit"
				return
			}
		}
	}()

	// Lets start somewhere.
	i := clix2.Coordinates{X: 1, Y: 1, Style: 1, Rune: rune(0)}

	/*
		   We can type on the screen.
			 First arg is the coordinates to start.
			 In this case is where the 'H' will be.
	*/
	win.TypeUI(i, "Hello worldz")
	i.Y++ // new line
	win.TypeUI(i, "Welcome!")
	// Present the window, refreshing the screen every input.
	// For example, a key press or mouse movement will refresh the window elements,
	// by doing one iteration of the loop.
	// We want this in a goroutine otherwise we can't get past it.
	go win.Present()

	/*
		   Loop until win.Output gets string "quit", Sending any other output to a file.
			 In a real program this loop would contain more logic,
			 and maybe even listen for receives on the win.Input channel.
	*/

	go func() {
		time.Sleep(1 * time.Second)
		ball := newBall(win.GetScreen())
		ball.Show()
		for {
			time.Sleep(100 * time.Millisecond)
			ball.moveBall(right)
			ball.Show()
			win.GetScreen().Show()
		}
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ball := newBall(win.GetScreen())
		go ball.Show()
		for {
			time.Sleep(100 * time.Millisecond)
			ball.moveBall(left)
			go ball.Show()
		}
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ball := newBall(win.GetScreen())
		ball.Show()
		for {
			time.Sleep(100 * time.Millisecond)
			ball.moveBall(down)
			ball.Show()
			win.GetScreen().Show()
		}
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ball := newBall(win.GetScreen())
		ball.Show()
		for {
			time.Sleep(100 * time.Millisecond)
			ball.moveBall(up)
			ball.Show()
			win.GetScreen().Show()
		}
	}()
Loop:
	for {
		select {
		case output := <-win.Output:
			switch output.(type) {
			case string:
				str := output.(string)
				if str == "quit" {
					//log.Println(`Received "quit" on output channel.`)
					break Loop
				}
				filer.Touch("output.log")
				filer.Append("output.log", []byte(fmt.Sprintf("%q\n", str)))
			}
		case <-time.After(30 * time.Second):
			/*
			 Get API refresh
			*/
			win.GetScreen().Show()
		}
	}

	win.Close()

	fmt.Println("Buffer dump:")
	fmt.Println(win.GetBuffer().String())
}
