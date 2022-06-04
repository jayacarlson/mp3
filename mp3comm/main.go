package main

import (
	"flag"

	"github.com/jayacarlson/dbg"
	"github.com/jayacarlson/mp3"
)

// just outputs the COMM frame from a file or list of files

func main() {
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		dbg.Fatal("Must supply filename(s) to dump")
	}

	for _, f := range files {
		m, err := mp3.OpenMP3File(f)
		if nil != err {
			dbg.Error("ERR: %v", err)
			continue
		}
		frm, err := m.FindFrame("COMM")
		if nil != err {
			dbg.Error("ERR: %v", err)
		} else {
			s, err := frm.ToString()
			if nil != err {
				dbg.Error("ERR: %v", err)
			} else {
				if len(files) > 1 {
					dbg.Echo("%s) %s", f, s[6:])
				} else {
					dbg.Echo("%s", s[6:])
				}
			}
		}
		m.Close()
	}
}
