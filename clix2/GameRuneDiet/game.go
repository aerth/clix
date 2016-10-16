package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/aerth/clix/clix2"
	"github.com/gdamore/tcell"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Ball to play
type Ball struct {
	difficulty int
	Coords     *clix2.Coordinates
	Runes      []*clix2.Coordinates
	screen     tcell.Screen
	Length     int
	alive      bool
	window     *clix2.Window
	foodEaten  int
}

func newBall(win *clix2.Window) *Ball {
	coords := new(clix2.Coordinates)
	coords.X = 10
	coords.Y = 10
	coords.Rune = '@'

	ball := new(Ball)
	ball.Coords = coords
	ball.screen = win.GetScreen()
	ball.window = win
	ball.Runes = append(ball.Runes, coords)
	ball.alive = true
	ball.difficulty = 1
	return ball
}

const (
	up    = 0
	down  = 1
	left  = 2
	right = 3
)

func makefruit(screen tcell.Screen) rune {
	var coords clix2.Coordinates
	maxx, maxy := screen.Size()
	coords.Screen = screen
	coords.X = rand.Intn(maxx)
	coords.Y = rand.Intn(maxy)
	coords.Style = 4
	runes := []rune("Helloo WorldoHellooWorldoo!1234567890z")
	runes = append(runes, []rune{tcell.RuneDiamond, tcell.RuneLantern, tcell.RunePi}...)
	coords.Rune = rune(runes[rand.Intn(len(runes))])
	coords.Draw()
	screen.Show()
	return coords.Rune
}

func playball(ball *Ball) {
	for {
		select {
		case input := <-ball.window.Input:
			switch input.(type) {
			default:
			case string:
				//log.Println("Input chan:", input)
				switch input.(string) {
				default:
				case "Enter":
					if !ball.alive {
						ball.window.Output <- ball.foodEaten
						stopfruit <- 1
					}
					return
				case "Up":
					ball.moveBall(up)

				case "Down":
					ball.moveBall(down)

				case "Left":
					ball.moveBall(left)

				case "Right":

					ball.moveBall(right)

				}
			}
		}
	}
}
func (b *Ball) relocate(coordinates clix2.Coordinates) {
	b.Hide()
	b.Coords.X = coordinates.X
	b.Coords.Y = coordinates.Y
	b.Show()
}
func (b *Ball) moveBall(direction int) {
	if !b.alive {
		return
	}
	b.Hide()
	maxx, maxy := b.screen.Size()
	maxy--
	maxx--
	switch direction {
	case up:
		if b.Coords.Y > 0 {
			b.Coords.Y--
		}
	case down:
		if b.Coords.Y < maxy {
			b.Coords.Y++
		}
	case left:
		if b.Coords.X > 0 {
			b.Coords.X--
		}
	case right:
		if b.Coords.X < maxx {
			b.Coords.X++
		}
	default:
		panic("cant move ball")
	}
	b.Show()
}

// Show a ball
func (b *Ball) Show() {
	if !b.alive {
		return
	}
	for _, r := range b.Runes {
		l, _, _, _ := b.screen.GetContent(r.X, r.Y)
		if l == rune('o') {
			b.Kill()
			return
		}

		if l == rune(32) || l == rune(0) { // Blanks
			r.Screen = b.screen
			r.Draw()
			b.screen.Show()
			return
		}

		if l == rune('z') {
			// teleport
			coords := clix2.Coordinates{X: rand.Intn(b.window.Xmax), Y: rand.Intn(b.window.Ymax)}
			b.relocate(coords)
			return
		}
		// Food eaten!

		str := fmt.Sprintf(strings.Repeat("o", 22)+"\n"+
			"ooo %q (Rune %v) ooo\n"+
			strings.Repeat("o", 22)+"\n", string(l), l)
		b.window.TypeUI(clix2.Coordinates{X: 1, Y: b.window.Ymax - 3}, str)

		b.foodEaten += int(l)

		r.Screen = b.screen
		r.Draw()
	}
	b.window.TypeUI(clix2.Coordinates{X: 1, Y: 3}, fmt.Sprintf("o Points: %+03v / %+04v o", b.foodEaten, (loopnum)*300))
	b.screen.Show()
}

// Hide a ball
func (b *Ball) Hide() {
	for _, r := range b.Runes {
		var co clix2.Coordinates
		co.Screen = b.screen
		co.Rune = rune(0)
		co.X = r.X
		co.Y = r.Y
		co.Draw()
	}
	b.screen.Show()

}

// Kill a ball
func (b *Ball) Kill() {
	for _, r := range b.Runes {
		var co clix2.Coordinates
		co.Screen = b.screen
		co.Rune = rune(0)
		co.X = r.X
		co.Y = r.Y
		co.Draw()
	}
	b.alive = false
	b.window.TypeUI(clix2.Coordinates{X: 1, Y: 3},
		fmt.Sprintf(strings.Repeat("o", 22)+"\n"+
			"Game Over.\nTotal: %+02v\n"+
			"Press ENTER"+
			strings.Repeat("o", 22)+"\n",
			b.foodEaten))
	b.screen.Show()
	b.window.Output <- b.foodEaten
	stopfruit <- 1

}

// Generate the fruit to eat (and to avoid)
func fruitgenerator(ball *Ball) {
	for {
		select {
		case <-stopfruit:
			return
		case <-time.After(900 * time.Millisecond):
			if ball.foodEaten >= (300*loopnum) && enemycount > 0 {
				ball.window.TypeUI(clix2.Coordinates{X: 1, Y: 6},
					fmt.Sprintf("Level %v completed! Score: %v",
						loopnum+1, ball.foodEaten))
				ball.screen.Show()
				ball.alive = false
				//log.Println("Finished level", loopnum+1)
				//stopfruit <- 1
				//log.Println("Finished level!", loopnum+1)
				return
			}
			i := makefruit(ball.window.GetScreen())
			if i != 'o' {
				fruitcount++
			} else {
				enemycount++
			}
		}
	}
}
