package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aerth/clix"
	"github.com/gdamore/tcell"
)

// var style1 = tcell.Style(8)
var style2 = tcell.Style(7)

// var style3 = tcell.Style(6)
// var style4 = tcell.Style(5)
// var style5 = tcell.Style(4)
var styleSelected = style2 + 8
var styleDefault = style2
var count int

func main() {
	flag.Parse()
	logfile = "debug.log"
	logger()
	doTitleMenu()
	doMenuBar()
}

var msg = "this is a message"

func doMenuBar() {

	events := clix.NewEventHandler()
	menu := clix.NewMenuBar(nil)

	scroller := clix.NewScrollFrame("Titlescroll")

	menu.NewItem("Red")
	menu.NewItem("Gold")
	menu.AttachScroller(scroller)
	scroller.Buffer.WriteString("Helo scrol world buffer")
	scroller.Buffer.WriteString(clix.Fill())

	menu.SetMessage(msg)
	menu.NewItem("Green")
	menu.AddFunc(tcell.KeyCtrlK, func(interface{}) interface{} {
		count++
		log.Println("CtrlK", count)
		return nil
	})

	menu.AddFunc(tcell.KeyCtrlL, func(i interface{}) interface{} {
		count++
		log.Println("CtrlL", i)
		return nil
	}, time.Now)

	menu2 := clix.NewMenuBar(menu.GetScreen())
	menu3 := clix.NewMenuBar(menu.GetScreen())
	//	menu4 := clix.NewMenuBar(menu.GetScreen())

	menu2.NewItem("Launch Entry Demo")
	menu3.NewItem("one")
	menu3.NewItem("two2")
	menu3.NewItem("three")

	menu.SetTitle("menunum1")
	menu2.SetTitle("menunum2")
	//	menu2.SetMessage("Welcomez")
	menu2.NewItem("Two")
	menu2.NewItem("threE")
	entry := clix.NewEntry(menu.GetScreen())
	menu.AddEntry("other", entry)

	menu.AddSibling(menu2)
	menu.AddSibling(menu3)

	var mainout interface{}

	// Timeout // Close menu from outside
	// go func() {
	// 	time.Sleep(10 * time.Second)
	// 	log.Println("trying timeout:")
	// 	menu.GetWidgetController().Input <- "end"
	// 	//menu.GetWidgetController().Input <- "minimize"
	// }()
	// Timeout // Close menu from outside
	go func() {
		time.Sleep(400 * time.Second)
		log.Println("trying test:")
		menu.GetWidgetController().Input <- "test"
		//menu.GetWidgetController().Input <- "minimize"
	}()

	//Start:
	/*
		events.Launch begins a couple goroutines to monitor input.
		To close it, send "end" string to events.Input channel
		To utilize it, use AddMenuBar or AddEntry** or AddTitleMenu**
	*/
	events.AddMenuBar(menu)
	events.Launch()
	select {

	// case l := <-menu.GetWidgetController().Output:
	// 	log.Println("MC Out:", l)
	case l := <-events.Output:
		log.Println("Output:", l)
		mainout = l

		switch l.(type) {
		case string:
			switch l.(string) {
			case "quit":

				goto Done

			case "fun":

			case "Launch Entry Demo":

				s := entrydemo(menu.GetScreen())
				mainout = s
				msg = s
				doMenuBar()

			case "TitleMenu Demo":

				doTitleMenu()

			}
		}

	}

Done:

	if menu.GetScreen() != nil {
		if menu.GetScreen().PollEvent() != nil {
			menu.GetScreen().Fini()
		}
	}

	if mainout != nil {
		fmt.Printf("\n\nGot output: %q.\n\n", mainout)
	}

	fmt.Println("Goodbye!")
}

func doTitleMenu() {
	t := clix.NewTitleMenu()
	q := `

	Greetings!


	Press ENTER to see the demo.
	`
	t.AddLines(strings.Split(q, "\n"))
	t.Present()
}
func entrydemo(s tcell.Screen) string {
	e := clix.NewEntry(s)

	e.AddPrompt("What is name")
	str := e.Present()
	log.Println(str)
	return str
}
