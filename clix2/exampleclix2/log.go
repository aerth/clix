// Copyright (c) 2016 aerth
// GPLv3
// Easy to use Logger file
// Just set var logfile = "name.txt" and then run logger() after flag.Parse()
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aerthlib/filer"
)

var logfile string

func init() {
	flag.StringVar(&logfile, "log", "debug.log", "Log to file. Default is no log.")
}

//logger switches the log engine to a file, rather than stdout.
func logger() {
	if logfile == "" {
		return
	}
	f, errar := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if errar != nil {
		fmt.Printf("error opening file: %v", errar)
		fmt.Println("Hint: touch " + logfile + ", or chown/chmod it.")
		os.Exit(1)
	}

	log.SetOutput(f)
	log.SetFlags(log.Lshortfile)
	filer.Touch(logfile)
	filer.Append(logfile, []byte("\n\n\n\n"))
	log.Println("Started", time.Now().Format(time.Kitchen))
}
