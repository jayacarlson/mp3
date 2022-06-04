package main

import (
	"flag"

	"github.com/jayacarlson/dbg"
	"github.com/jayacarlson/mp3"
)

var (
	frmID string
)

func init() {
	flag.StringVar(&frmID, "id", "", "Output ID")
}

func main() {
	flag.Parse()

	if "" == frmID {
		dbg.Fatal("Must supply -id to output.  Usage: showid3Frame -id XXXX <filepath>")
	}
	if mp3.InvalidID3v2FrameID(frmID) {
		dbg.Fatal("Illegal frame -id.  Usage: showid3Frame -id XXXX <filepath>")
	}
	files := flag.Args()
	if len(files) == 0 {
		dbg.Fatal("Must supply filename to dump.  Usage: showid3Frame -id XXXX <filepath>")
	}

	for _, f := range files {
		m, err := mp3.OpenMP3File(f)
		if nil != err {
			dbg.Error("ERR: %v", err)
			continue
		}
		frm, err := m.FindFrame(frmID)
		if nil != err {
			dbg.Error("ERR: %v", err)
		} else {
			s, err := frm.ToString()
			if nil != err {
				dbg.Error("ERR: %v", err)
			} else {
				if len(files) > 1 {
					dbg.Echo("%s) %s", f, s)
				} else {
					dbg.Echo("%s", s)
				}
			}
		}
		m.Close()
	}
}
