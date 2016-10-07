// Skeleton clix CUI library example app
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aerthlib/clix"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func initialize() {
	flag.Parse()
	logger()
}

func main() {
	initialize()
	m := clix.NewMenuBar()
	m.SetTitle("[Disaster]")
	m.NewItem("Earthquake")
	m.NewItem("Tornado")
	m.NewItem("Wildfire")
	m.NewItem("Blizzard")
	m.NewItem("Flood")
	m.NewItem("Hurricane")
	m.NewItem("Mudslide")
	m.NewItem("Other")

	m2 := clix.NewMenuBar()
	m2.SetTitle("[Number]")
	for i := 0; i < 7; i++ {
		m2.NewItem(strconv.Itoa(i))
	}
	m3 := clix.NewMenuBar()
	m3.SetTitle("[Region]")
	for _, v := range []string{"North", "South", "East", "West", "Within", "Without", "Up", "Down"} {
		m3.NewItem(v)
	}
	m.AddSibling(m2)
	m.AddSibling(m3)
Begin: // All the above is not to be repeated every loop
	s := clix.Load()
	s.HideCursor()
	m.Present(s, true)
	userSelection := m.MainLoop(s)
	s.Clear()
	var o string
	var outs = map[string]string{}
	log.Println(userSelection)
Name: // Cut right to a specific menu item

	o = clix.Entry(s, "How many fingers am I holding up?")
	log.Println(o)
	o = clix.Entry(s, "Does \""+userSelection+"\" sound familiar?")
	log.Println(o)
	if o == "yes" {
		clix.StdOut(s)
		fmt.Println(userSelection)
		os.Exit(0)
	}

	outs["name"] = o
	o = clix.Entry(s, "What is your name?")
	outs["name"] = o
	o = clix.Entry(s, "What is your favorite software library?")
	outs["lib"] = o
	o = clix.Entry(s, "Mouse is cool!")
	outs["lib"] = o
	o = clix.Entry(s, "What is your name again?")
	if outs["name"] != o {
		// User error
		clix.Eat(s, 0)
		clix.UnEat(s, 0)
		goto Name
	}
	clix.StdOut(s) // Back to fmt.Println-able (and ctrl+c able) area

	for _, v := range fmt.Sprintf("%q chooses %q\n", strings.ToUpper(o), strings.ToLower(userSelection)) { // Output selection
		fmt.Print(string(v))
		time.Sleep(20 * time.Millisecond)
	}
	for _, v := range fmt.Sprintf("%s\n", outs) { // Output selection
		fmt.Print(string(v))
		time.Sleep(20 * time.Millisecond)
	}

	str := clix.EntryOS("Type YES to start over!")
	if str == "YES" {
		goto Begin
	}
}
