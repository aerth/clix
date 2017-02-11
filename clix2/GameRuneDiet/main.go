package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aerth/clix"
	"github.com/aerth/clix/clix2"
)

var version = "RuneDiet v1"
var stopfruit = make(chan int)
var enemycount, fruitcount int
var loopnum int
var totalscore int

func main() {

	flag.Parse()
	logger()
	logfile = "debug.log"

Start:
	loopnum++
	win := clix2.NewWindow("Hello")
	win.GetScreen().DisableMouse()
	/*
		Q:	Why doesn't this countdown variable cut games off halfway through level two?
	*/
	go func() {
		var countdown = 60
		c2 := clix2.Coordinates{X: 1, Y: win.Ymax - 2}
		for {
			win.TypeUI(clix2.Coordinates{X: win.Xmax - 2, Y: win.Ymax - 1}, strconv.Itoa(countdown))
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

	header(win)
	go win.Present()

	ball := newBall(win)
	ball.difficulty = loopnum
	//log.Println("Level", ball.difficulty)
	go playball(ball)
	for i := 0; i < ball.difficulty; i++ {
		go fruitgenerator(ball)
	}
Loop:
	for {
		select {
		case output := <-win.Output:
			switch output.(type) {
			case int: // is a score!
				score := output.(int)
				//log.Println("Received Score:", score)
				break Loop
			case string:
				str := output.(string)
				if str == "quit" {
					//log.Println(`Received "quit" on output channel.`)
					break Loop
				} else {
					//log.Println("Output channel:", str)
					s := clean(str)
					if s == "" && !ball.alive {
						break Loop
					}
					if s == "cheat" {
						for i := 0; i < 10; i++ {
							go fruitgenerator(ball)
						}
					}
					win.GetBuffer().Truncate(0)
				}
			}
		case <-time.After(30 * time.Second):
			/*
			 Get API refresh
			*/
			win.GetScreen().Show()
		}
	}

	win.Close()
	if win.GetBuffer().Len() != 0 {
		lines := strings.Split(win.GetBuffer().String(), "\n")
		for _, v := range lines {
			fmt.Println(v)
		}
	}
	if ball.foodEaten != 0 {
		fmt.Println("Your score:", ball.foodEaten)
		b := ball.foodEaten
		totalscore += ball.foodEaten
		var msg string
		switch {
		case b < 10:
			msg = "You did horrible!"
		case b < 100:
			msg = "You did okay!"
		case b < 300:
			msg = "Alright!"
		case b < 500:
			msg = "You are good!"
		default:
			msg = "Great Job!"
		}
		fmt.Println(msg)
		if totalscore != 0 {
			fmt.Println("Total Score:", totalscore)
		}
	}
	if ball.foodEaten < 300*loopnum {
		loopnum = 0 // Will turn to 1 at Start tag
	}
	s := clix.EntryOS(fmt.Sprintf("Play again? Press ENTER for Level %v", loopnum+1))
	if s == "yes" || s == "" {

		goto Start
	}

	fmt.Println("Thank you for playing", version)

}
func header(win *clix2.Window) {

	i := clix2.Coordinates{X: 1, Y: 1, Style: 1, Rune: rune(0)}
	win.TypeUI(i, strings.Repeat("o", 22)) //
	i.Y++
	win.TypeUI(i, "o RuneDiet           o")
	i.Y++
	win.TypeUI(i, "o Don't eat the o's! o")
	i.Y++
	win.TypeUI(i, strings.Repeat("o", 22)) //
}
func clean(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, " ", "", -1)
	return s

}
