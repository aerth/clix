package clix

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

/*
Test going in and out of the screen and back to STDOUT.
Test we dont mess up the terminal
Test we dont break API by changing function names etc
*/
func init() {
	null, _ := os.OpenFile(os.DevNull, os.O_WRONLY, os.ModeAppend)
	log.SetOutput(null)
}
func TestTitleMenu(t *testing.T) {
	mm := NewTitleMenu()
		line1 := "Prepare to calibrate your joystick."
		line2 := "It may feel like its a human test, but really..."
	mm.AddLine(line1)
		mm.AddLine(line2)
			mm.AddLine("press any key to continue")
	mm.Present()
	fmt.Println("You did it.")
	time.Sleep(1 * time.Second)
	fmt.Println("Lets try again.")
	time.Sleep(1 * time.Second)
	mm2 := NewTitleMenu()
	mm2.AddLine("press ENTER one more time!")
	mm2.Present()
}

func TestMenuBar(t *testing.T) {
	mb := NewMenuBar(nil)
	mb.SetTitle("Select GOLD")
	mb.NewItem("One")
	mb.NewItem("Two")
	mb.NewItem("GOLD")
	events := NewEventHandler()
	events.AddMenuBar(mb)
	events.Launch()
	var str string
	for {
		select {
		case l := <-events.Input:
			fmt.Println(l)
			str = l.(string)
			break
		case l := <-events.Output:
			str = l.(string)
			fmt.Println(l)
			break
		}
		mb.GetScreen().Fini()
		break
	}
	fmt.Println("Got:", str)
	if str != "GOLD" {
		fmt.Println("Round two Got:", str)
		t.Fail()
	}
	mb2 := NewMenuBar(nil)
	mb2.SetTitle("Select GOLD")
	mb2.NewItem("Two One")
	mb2.NewItem("GOLD")
	mb2.NewItem("Two Three")
	events2 := NewEventHandler()
	events2.AddMenuBar(mb2)
	events2.Launch()
	var str2 string
	for {
		select {
		case l := <-events2.Input:
			fmt.Println(l)
			str2 = l.(string)
			break
		case l := <-events2.Output:
			str2 = l.(string)
			fmt.Println(l)
			break
		}
		mb2.GetScreen().Fini()
		break
	}
	if str2 != "GOLD" {
		fmt.Println("Round two Got:", str2)
		t.FailNow()
	}
}

func TestEntry(t *testing.T) {
	e := NewEntry(nil)
	e.AddPrompt("Greetings, tester...")
	e.AddPrompt("Are you testing? Type: 1")
	one := e.Present()
	if one != "1" {
		t.FailNow()
	}
	e = NewEntry(nil)
	e.AddPrompt("Round two: Type: 2")
	two := e.Present()
	if two != "2" {
		t.FailNow()
	}
	fmt.Println("Now level 3")
	time.Sleep(1 * time.Second)
	e = NewEntry(nil)
	e.AddPrompt("Round 3: Type: 3")
	tree := e.Present()
	if tree != "3" {
		t.FailNow()
	}

}
