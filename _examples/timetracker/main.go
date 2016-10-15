/*

Alice and Bob always get to work at 9:00am and leave exactly at 5:00pm.
Hours like that makes counting a breeze.
The trouble comes when mallory is scheduled at 3:35pm to 8:20pm and 1:15pm to 7:22pm
Some durations are more difficult to calculate than others.
Allow TimeTracker to do the math for you.


TODO:
  less bugs
  less loops
  less buffer
  less memusage
  less jitter when refreshing screen?
  json input
  csv input
  png graph? :D
  html out
  pdf output
  full zip output (all of outputs in one zip for no reason)
  json output
  csv output
  mouse buttons for repeat and common times (9-5, each hour :00, and :30. maybe 15 too )
  database
  ability to guarantee the time table has not been manipulated? nah
*/

//package timetracker is a time shell to keep track of your workers
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aerth/clix"

	"github.com/gdamore/tcell"
)

var version = "TimeTracker v1"
var timetable []WorkPeriod
var curtime time.Duration

var msg string

// WorkPeriod time
type WorkPeriod struct {
	In  time.Time
	Out time.Time
	Len time.Duration
}

var (
	outflag = flag.String("o", "", "Output file. Alternatively can be an os arg such as: timecalc octoberJohn.txt")
)

func main() {
	flag.Parse()
	logger() // to /dev/null or -log flag
	arg := flag.Args()

	if len(arg) != 1 && *outflag == "" {
		fmt.Println("Need output file. Use -o flag or type " + os.Args[0] + " filename")
		os.Exit(1)
	}

	if *outflag == "" && len(arg) == 1 {
		*outflag = arg[0]
	}

	if *outflag == "" {
		flag.Usage()
		os.Exit(1)
	}

	Loop()

}

// Loop until user types stop or ctrl+c
func Loop() {
	var looped int
	var s tcell.Screen
NewReport:
	looped++

	if *outflag != "/dev/null" {
		_, errer := os.Open(*outflag)
		if errer == nil {
			fmt.Println("File already exists. Not overwriting. Try again.")
			os.Exit(1)
		}
	} else {
		fmt.Println("Discarding output to", os.DevNull)
		time.Sleep(1 * time.Second)
	}

	var ins, o string
	var err error
	var clockout, clockin time.Time
	if looped > 1 {
		s.Clear()
	}

	for i := 0; ; i++ { // Break loop with "stop" command.
		var period WorkPeriod

	ClockIn:
		//		_, ymax := s.Size()
		inentry := clix.NewEntry(nil)
		s = inentry.GetScreen()
		inentry.AddPrompt("[Clock In] Total: " + curtime.String())
		inentry.SetTitle("Timetable " + version)
		inentry.SetSubtitle("Saving to: " + *outflag)
		inentry.AddPrompt("Formatting Help: 3:15pm, 15:15, or 1515")
		inentry.AddPrompt("Insert CLOCK-IN time")
		if len(timetable) != 0 {
			inentry.AddPrompt("When finished, type 'stop'")
		}
		// add msg to prompt if any
		inentry.AddPrompt(msg)
		msg = "" // clear msg after using it

		ins = inentry.Present()

		// Trim up input
		ins = strings.Replace(ins, " ", "", -1)
		ins = strings.ToUpper(ins)
		if ins == "" {
			goto ClockIn
		}

		// Break for loop
		if ins == "STOP" || ins == "DONE" || ins == "SAVE" || ins == "QUIT" || ins == "EXIT" {
			goto Done
		}

		// r repeat command: repeats the last successful entry
		if ins == "R" && len(timetable) > 0 {
			var repeated WorkPeriod
			repeated.In = timetable[len(timetable)-1].In
			repeated.Out = timetable[len(timetable)-1].Out
			repeated.Len = timetable[len(timetable)-1].Len
			msg = fmt.Sprintf("Repeated: %s to %s = %s", repeated.In, repeated.Out, repeated.Len)
			curtime += repeated.Len
			log.Println("Session:", repeated.Len, "Total:", curtime)
			timetable = append(timetable, repeated)
			goto ClockIn
		}

		if len(ins) == 4 && !strings.HasSuffix(ins, "M") {
			ins = twentyfour(ins)
			if ins == "" {
				goto ClockIn
			}
			//clix.Type(s, 0, ymax-4, 3, fmt.Sprintf("%s", ins))
		}
		if len(ins) == 3 && strings.HasSuffix(ins, "M") {
			// byte surgery adding a :00, making 3pm 3:00pm

			ins = string([]byte{ins[0], byte(':'), byte('0'), byte('0'), ins[1], ins[2]})
		}
		if len(ins) == 4 && strings.HasSuffix(ins, "M") {
			// byte surgery adding a :00, making 12pm 12:00pm
			ins = string([]byte{ins[0], ins[1], byte(':'), byte('0'), byte('0'), ins[2], ins[3]})
		}

		if strings.Contains(ins, ":") && strings.Contains(ins, "M") {
			clockin, err = time.Parse(time.Kitchen, ins)
		} else {
			clockin, err = time.Parse("15:04", ins)
		}
		if err != nil {
			log.Println(err)
			msg = err.Error()
			goto ClockIn
		}
		log.Println("Clock In", i, clockin.Format(time.Kitchen))
		period.In = clockin
		//	outerr = nil
	ClockOut:
		//	inerr = nil
		//		_, ymax = s.Size()
		//	s.Clear() // dont need this
		// if outerr != nil {
		// 	clix.Type(s, 0, ymax-3, 2, outerr.Error())
		// 	s.Show()
		// }
		oentry := clix.NewEntry(s)
		oentry.SetTitle(version)
		oentry.SetSubtitle(*outflag)
		oentry.AddPrompt("Insert CLOCK-OUT time")
		oentry.AddPrompt("[Clock Out] Total: " + curtime.String())
		oentry.AddPrompt("Formatting Help: 3:15pm, 15:15, or 1515")
		oentry.AddPrompt("[Clock In: " + period.In.Format(time.Kitchen) + "]")
		oentry.AddPrompt(msg)
		msg = ""
		o = oentry.Present()
		o = strings.Replace(o, " ", "", -1)
		o = strings.ToUpper(o)
		if o == "" {
			goto ClockOut
		}
		if o == "STOP" || o == "DONE" || o == "SAVE" || o == "QUIT" || o == "EXIT" {

			goto Done
		}

		if len(o) == 4 && !strings.HasSuffix(o, "M") {
			o = twentyfour(o) // convert 1200 to 12:00
			msg = o
			if o == "" {
				goto ClockOut
			}
		}
		if len(o) == 3 && strings.HasSuffix(o, "M") {
			// byte surgery adding a :00, making 3pm 3:00pm
			o = string([]byte{o[0], byte(':'), byte('0'), byte('0'), o[1], o[2]})
		}
		if len(o) == 4 && strings.HasSuffix(o, "M") {
			// byte surgery adding a :00, making 12pm 12:00pm
			o = string([]byte{o[0], o[1], byte(':'), byte('0'), byte('0'), o[2], o[3]})
		}

		if strings.Contains(o, ":") && strings.HasSuffix(o, "M") {
			clockout, err = time.Parse(time.Kitchen, o)
		} else {
			clockout, err = time.Parse("15:04", o)
		}

		if err != nil {
			log.Println(err)
			msg = err.Error()
			goto ClockOut
		}

		// Compute
		period.Out = clockout
		durat := clockout.Sub(clockin)
		if int64(durat) < 0 {
			durat = clockout.Add(24 * time.Hour).Sub(clockin)
		}
		period.Len = durat

		// Total current duration
		curtime += durat

		// Send to log
		log.Println("Clock Out", i, clockout.Format(time.Kitchen))
		log.Println("Session:", durat, "Total:", curtime)

		// Add to current timetable
		timetable = append(timetable, period)

		// Display to user
		msg = fmt.Sprintf("'r' to repeat: %s to %s = %s",
			period.In.Format(time.Kitchen),
			period.Out.Format(time.Kitchen),
			period.Len)
	}

Done:
	var stdout string
	if len(timetable) > 0 {
		msg, stdout = GenerateReport()
		if *outflag != os.DevNull {
			cont := clix.NewEntry(s)
			cont.AddPrompt(msg)
			msg = ""
			cont.SetTitle(version)
			cont.AddPrompt("Would you like to create another timetable?")
			resp := cont.Present()
			if resp == "yes" {

				cont = clix.NewEntry(s)
				cont.AddPrompt("Filename")
				cont.SetTitle(version)
				*outflag = cont.Present()

				goto NewReport
			}
		}
	}
	clix.StdOut(s)
	if msg != "" {
		fmt.Println("\n\n" + msg + "\n\n")
	}
	if stdout != "" {
		fmt.Println("\n\n" + stdout + "\n\n")
	}
}

// GenerateReport prints report to stdout and writes to *filename
func GenerateReport() (string, string) {

	var b bytes.Buffer
	fmt.Fprintln(&b, "Timetable Report")
	fmt.Fprintln(&b, *outflag+"\n")
	fmt.Fprintf(&b, "Created %s\n\n", time.Now().String())
	fmt.Fprintf(&b, "%s \t%+2s \t%+3s\n", "clock in", "clock out", "duration")
	fmt.Fprintf(&b, "%s \t %+2s \t %+3s\n", "--------", "---------", "--------")

	for _, v := range timetable {
		fmt.Fprintf(&b, "%+0s \t\t%+0s\t\t%-2s\n",
			v.In.Format(time.Kitchen),
			v.Out.Format(time.Kitchen), v.Len)
	}
	fmt.Fprintln(&b, "Total:", curtime)

	fmt.Fprintln(&b, "Number of sessions:", len(timetable))
	if *outflag == "/dev/null" {
		*outflag = os.DevNull
	}
	err := ioutil.WriteFile(*outflag, b.Bytes(), 0700)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("Wrote %v bytes to %q\n", b.Len(), *outflag)
	}
	sout := fmt.Sprintf("Wrote %v bytes to %q\n", b.Len(), *outflag)
	// Clear info
	timetable = nil
	curtime = 0
	return sout, b.String()
}

func twentyfour(s string) string {
	if len(s) != 4 {
		return s
	}

	num1, err := strconv.Atoi(string([]byte{s[0], s[1]}))
	if err != nil {
		return s
	}
	num2, err := strconv.Atoi(string([]byte{s[2], s[3]}))
	if err != nil {
		return s
	}
	return fmt.Sprintf("%+02v:%+02v", num1, num2)
}

func puts(s tcell.Screen, str string) {
	xmax, ymax := s.Size()
	for j := 0; j < xmax; j++ {

		s.SetCell(j, ymax-3, tcell.Style(1), rune(' '))
	}
	clix.Type(s, 0, ymax-3, 2, str)
	s.Show()
}
