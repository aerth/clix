package clix

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

// Type to a screen, one line no wrap.
// If len(s) > ymax-starty, s will go off the screen.
func Type(scr tcell.Screen, startx int, starty int, style tcell.Style, s string) {
	x, y := startx, starty
	for i := 0; i < len(s); i++ {
		scr.SetCell(x+i, y, style, rune(s[i]))
	}
}

// TypeUI is like Type, but writes to the special parts, such as ymax-1
func TypeUI(screen tcell.Screen, style tcell.Style, startx, starty int, s string) (endx, endy int) {
	posx, posy := startx, starty
	paragraph := []rune(s)
	xmax, ymax := screen.Size()
	for _, letter := range paragraph {
		if posx == xmax-1 {
			posy++
			posx = 0
		}
		if posy > ymax { // input bar = 3, output box = everything else
			posy = 0
			posx = 0
		}
		screen.SetCell(posx, posy, style, letter)
		posx++

	}
	return posx, posy
}

// TypeWriter to a screen starting at coordinates: startx,starty.
// All lines are wrapped, so don't assume endy is the same as starty
func TypeWriter(scr tcell.Screen, startx int, starty int, style tcell.Style, s string) (endy int) {
	xmax, ymax := scr.Size()

	x, y := startx, starty
	for i := 0; i < len(s); i++ {
		//		time.Sleep(2 * time.Millisecond)
		x++
		if x >= xmax {
			x = 0
			y++
		}
		if y > ymax {
			y = 0
			x = 0
		}
		scr.SetCell(x, y, style, rune(s[i]))
	}
	scr.Show()
	return y
}

// Eat eats the screen with random runes at (timing). A good timing is 500, 50, or 1 or 0.
func Eat(screen tcell.Screen, timing int) {
	if timing < 0 {
		timing = 500
	}
	xmax, ymax := screen.Size()
	for i := 0; i < ((xmax*ymax)+(xmax*ymax))*2; i++ {
		go func() {
			if timing > 0 {
				time.Sleep(time.Duration(int64(rand.Intn(timing)) * int64(time.Millisecond)))
			}
			screen.SetCell(rand.Intn(xmax), rand.Intn(ymax), tcell.Style(rand.Intn(255)), rune(rand.Intn(800)))
			screen.Show()
		}()
	}
}

// UnEat gradually clears the screen. A good timing is 500, 50, or 1 or 0.
func UnEat(screen tcell.Screen, timing int) {
	if timing < 0 {
		timing = 500
	}
	xmax, ymax := screen.Size()
	for i := 0; i < ((xmax*ymax)+(xmax*ymax))*4; i++ {

		if timing > 0 {
			time.Sleep(time.Duration(int64(rand.Intn(timing)) * int64(time.Millisecond)))
		}
		screen.SetCell(rand.Intn(xmax), rand.Intn(ymax), style2, rune(0))
		screen.Show()

	}
	screen.Clear()
	screen.Show()

}
