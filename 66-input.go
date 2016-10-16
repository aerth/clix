package clix

import
//"log"

"github.com/gdamore/tcell"

var promptText = "[cmd]"
var numchar int
var color1 = tcell.StyleDefault
var inserting bool

var selchar int
var style1 = tcell.StyleDefault.Background(2).Foreground(0)
var style2 = tcell.StyleDefault.Background(0).Foreground(2)
var style3 = tcell.StyleDefault.Background(0).Foreground(4)

func getText(screen tcell.Screen) string {
	_, ymax := screen.Size()
	var input []rune
	for posx := 2; posx < numchar+1; posx++ {
		rune1, _, _, _ := screen.GetContent(posx, ymax-1)
		input = append(input, rune1)
	}
	return string(input)
}

// Peeks at text from input bar, returns as a string.
func peekInput(screen tcell.Screen) string {
	_, ymax := screen.Size()
	var input []rune
	for posx := 2; posx <= numchar; posx++ {
		rune1, _, _, _ := screen.GetContent(posx, ymax-1)
		input = append(input, rune1)
	}
	return string(input)
}

func getTextMin(screen tcell.Screen) string {
	xmax, ymax := screen.Size()
	var input []rune
	////log.Println("Getting", numchar)
	for posx := xmax - 20 + 1; posx < numchar+xmax-20+1; posx++ {

		rune1, _, _, _ := screen.GetContent(posx, ymax-1)
		screen.SetCell(posx, ymax-1, tcell.StyleDefault, rune(0))
		////log.Println(rune1, posx)
		input = append(input, rune1)
	}

	screen.HideCursor()
	////log.Printf("got text: %q\n", string(input))
	return string(input)
}

// Peeks at text from input bar, returns as a string.
func peekInputMin(screen tcell.Screen) string {
	xmax, ymax := screen.Size()
	var input []rune
	////log.Println("Getting", numchar)
	for posx := xmax - 20; posx <= numchar+xmax-20+1; posx++ {
		rune1, _, _, _ := screen.GetContent(posx, ymax)
		input = append(input, rune1)
	}
	return string(input)
}

// Draw the actual '>'
func drawPromptRune(screen tcell.Screen) {
	_, ymax := screen.Size()
	screen.SetContent(0, ymax-1, '>', nil, color1)
	// bool flags... Way to simplify this?
	numflags := 0
	if inserting {
		screen.SetContent(len(promptText)+numflags, ymax-2, tcell.RuneDegree, nil, color1)
		numflags++
	}
}

// Clear and Draw the Inputbar
func drawInputBar(screen tcell.Screen, code int) {
	_, ymax := screen.Size()
	if code == 1 {
		numchar = 1 // no char
		selchar = 1 // no selection
		clearline(screen, ymax-1)
		//clearline(screen, ymax-2)
	}
	posx := 1

	posy := ymax - 1
	xmax := 5
	drawPromptRune(screen)
	for {
		if posy > ymax {
			screen.ShowCursor(numchar+1, ymax-1)
			screen.Show()
			return
		}
		if posx > xmax {
			posx = 0
			posy++
		}
		if code == 1 {
			screen.SetCell(posx, posy, color1, 0)
		}
		posx++
	}
}

func clearSpace(screen tcell.Screen, i int) {
	_, ymax := screen.Size()
	for i := 0; i < ymax-3; i++ {
		clearline(screen, i)
	}
}
func clearline(screen tcell.Screen, y int) {
	xmax, ymax := screen.Size()
	if y < ymax {
		for i := 0; i < xmax; i++ {
			screen.SetCell(i, y, tcell.StyleDefault, rune(0))
		}
	}
}
func clearchar(screen tcell.Screen, x, y, num int) {
	xmax, ymax := screen.Size()
	if num > (xmax * ymax) {
		num = xmax * ymax
	}
	if y < ymax {
		for i := 0; i < num; i++ {
			x++
			if x > xmax {
				y++
				x = 0
			}
			if y > ymax {
				return
			}
			screen.SetCell(x, y, tcell.StyleDefault, rune(0))
		}
	}
}
