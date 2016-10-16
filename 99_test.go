package clix

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	null, _ := os.OpenFile(os.DevNull, os.O_WRONLY, os.ModeAppend)
	log.SetOutput(null)
}

func TestChoose(t *testing.T) {
	mb := NewMenuBar(nil)
	mb.SetTitle("Select TESTS to run!")
	mb.NewItem("ALL")
	mb.NewItem("MACHINE")
	mb.NewItem("HUMAN")
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
	fmt.Println("Running tests:", str)

	// These tests are placed in a map,
	// So the order they are ran are at random.
	// This helps find issues between widgets.
	tests := map[int]func(*testing.T){}
	switch str {
	case "HUMAN":
		tests[1] = humanTitleMenu
		tests[2] = humanMenuBar
		tests[3] = humanEntry
		tests[4] = humanMenuBarMulti
	case "MACHINE":
	case "ALL":
		tests[1] = humanTitleMenu
		tests[2] = humanMenuBar
		tests[3] = humanEntry
		tests[4] = humanMenuBarMulti
	}
	// Try everything twice
	for try := 1; try <= 2; try++ {
		for _, f := range tests {
			f(t)
		}
	}
}

func humanTitleMenu(t *testing.T) {
	mm := NewTitleMenu()
	lines := `

	Welcome to the clix test suite.

	You have chosed "HUMAN" tests, which require human input.
	`
	mm.AddLines(strings.Split(lines, "\n"))
	mm.AddLine("press any key to continue")
	mm.Present()
	fmt.Println("You did it.")
	time.Sleep(1 * time.Second)
	fmt.Println("Lets try creating a new TitleMenu.")
	time.Sleep(1 * time.Second)
	mm2 := NewTitleMenu()
	mm2.AddLine("press any key to continue")
	mm2.Present()
}

func humanMenuBar(t *testing.T) {
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
		t.FailNow()
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

func humanEntry(t *testing.T) {
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

func humanMenuBarMulti(t *testing.T) {
	mb := NewMenuBar(nil)
	mb.SetTitle("Select GOLD")
	mb.NewItem("One")
	mb.NewItem("Two")
	mb.NewItem("RUST")
	mb2 := NewMenuBar(mb.GetScreen())
	mb2.SetTitle("Select GOLD")
	mb2.NewItem("Яабвгде")
	mb2.NewItem("Testing")
	mb2.NewItem("One Two")
	mb2.NewItem("One Two Three")
	mb2.NewItem("ԊԋԌԍԎԏԐԑԒԓԔԕԖԗ")
	mb2.NewItem("GOLD")
	mb.AddSibling(mb2)

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
		t.FailNow()
	}

}

func machineTest(t *testing.T) {

}

func TestEntry(t *testing.T) {
	e := NewEntry(nil)
	e.AddPrompt("Hello! Enter 'GOLD' to win!")
	str := e.Present()
	fmt.Println("Got:", str)
	assert.Equal(t, "GOLD", str)
}
func TestEvents(t *testing.T) {
	// ev := NewEventHandler()
	// evdone := ev.Launch()
	//
	// go func() {
	// 	time.Sleep(1 * time.Second)
	// 	evdone.Done()
	// }()
	//
	// for {
	//
	// }
}
