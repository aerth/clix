package clix

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

// TypeUI to the special parts
func TypeUI(screen tcell.Screen, style tcell.Style, startx, starty int, s string) (endx, endy int) {

	posx, posy := startx, starty
	paragraph := []rune(s)
	// if *verbose {
	// 	log.Println("TypeUI", s, "to", posx, posy)
	// }
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

// Type to a screen, one line no wrap
func Type(scr tcell.Screen, startx int, starty int, style tcell.Style, s string) {
	x, y := startx, starty
	for i := 0; i < len(s); i++ {
		scr.SetCell(x+i, y, style, rune(s[i]))
	}
}

// TypeWriter to a screen, all lines with wrap
func TypeWriter(scr tcell.Screen, startx int, starty int, style tcell.Style, s string) {
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
}

// Eat eats the screen at timing. A good timing is 500, 50, or 1 or 0.
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

// UnEat eats the screen at timing. A good timing is 500, 50, or 1 or 0.
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
