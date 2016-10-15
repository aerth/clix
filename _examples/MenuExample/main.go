package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aerth/clix"
	"github.com/gdamore/tcell"
)

var style1 = tcell.Style(8)
var style2 = tcell.Style(7)
var style3 = tcell.Style(6)
var style4 = tcell.Style(5)
var style5 = tcell.Style(4)
var styleSelected = style2 + 8
var styleDefault = style2
var count int

func main() {
	flag.Parse()
	logfile = "debug.log"
	logger()
	events := clix.NewEventHandler()
	menu := clix.NewMenuBar(nil)

	scroller := clix.NewScrollFrame("Titlescroll")

	menu.NewItem("Red")
	menu.NewItem("Gold")
	menu.AttachScroller(scroller)
	scroller.Buffer.WriteString("Helo scrol world buffer")
	scroller.Buffer.WriteString(clix.Fill())
	//menu.SetMessage("this is a message")
	menu.NewItem("Green")
	menu.AddFunc(tcell.KeyCtrlK, func(interface{}) interface{} {
		count++
		log.Println("CtrlK", count)
		return nil
	})

	data := time.Now()
	menu.AddFunc(tcell.KeyCtrlL, func(i interface{}) interface{} {
		count++
		log.Println("CtrlL", i)
		return nil
	}, data)

	menu2 := clix.NewMenuBar(menu.GetScreen())
	menu2.NewItem("one")
	menu.SetTitle("menunum1")
	menu2.SetTitle("menunum2")
	//	menu2.SetMessage("Welcomez")
	menu2.NewItem("Two")
	menu2.NewItem("threE")
	entry := clix.NewEntry(menu.GetScreen())
	menu.AddEntry("other", entry)

	menu.AddSibling(menu2)

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
			if l == "quit" {
				goto Done
			}
			if l == "fun" {

			}

		}

	}
Done:
	if menu.GetScreen() != nil {
		menu.GetScreen().Fini()
	}
	if mainout != nil {
		fmt.Printf("Got output: %q. ", mainout)
	}

	fmt.Println("Goodbye!")
}
