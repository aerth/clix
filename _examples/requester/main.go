package main

import (
	"github.com/aerth/clix"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	s := clix.Load()
	s.EnableMouse()

	s.SetStyle(4)

	mm := clix.NewMainMenu()
	mm.SetJustifyRight()
	mm.SetTitle("REQUESTOR v1")
	mm.AddLine("Welcome to Requestor")
	mm.AddLine("You can generate a nice request with this tool.")

	var art = `
                           _
  ___ ___  ___  __ _  ___ | |
 / __/ _ \/ __|/ _  |/ _ \| |
| (_| (_) \__ \ (_| | (_) |_|
 \___\___/|___/\__, |\___/(_)
               |___/
             http://aerth.xyz
                            `

	mm.AddLines(strings.Split(art, "\n"))
	mm.Present(s)

	//s = clix.Load()
	m := clix.NewMenuBar()

	m.NewItem("GET FAST")
	m.Present(s, true)
	method := m.MainLoop(s)
	if method == "GET FAST" {
		//s = clix.Load()
		ustr := clix.Entry(s, "URL")
		url, err := url.Parse(ustr)
		if err != nil {
			fmt.Println(err)
			panic(err)
			goto After
		}
		resp, err := http.Get(url.String())
		if err != nil {
			panic(err)
			goto After
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
			goto After
		}
		clix.TypeWriter(s, 0, 0, 0, string(b))
		s.Show()

	}

After:
	s.Show()
	time.Sleep(3 * time.Second)
	s.Fini()
}
