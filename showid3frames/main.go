package main

import (
	"flag"

	"github.com/jayacarlson/dbg"
	"github.com/jayacarlson/mp3"
)

func DumpID3v2Info(m *mp3.MP3File) error {
	err := m.ValidateID3v2File()
	if nil != err {
		return err
	}
	path := m.Path()
	verMajor, verMinor := m.ID3v2Version()
	flags := m.ID3v2Flags()

	if !(2 <= verMajor && 4 >= verMajor) {
		return mp3.Err_UnknownID3v2Version
	}

	dbg.Info("File `%s` uses ID3 v2.%d.%d  (x%02x)  %d x%x", path, verMajor, verMinor, flags, m.ID3Area(), m.ID3Area())
	if 2 == verMajor {
		return mp3.Err_ObsoleteID3v2Version
	}

	for {
		frm, err := m.ReadID3v2Frame()
		if nil != err {
			if mp3.Err_EOF == err {
				break
			}
			return err
		}
		dbg.Echo("  %s: (%02x,%02x)  %d", frm.Tag, frm.Flags[0], frm.Flags[1], frm.Size)

		s, err := frm.ToString()
		if nil != err {
			return err
		}
		dbg.Message("     `%s`", s)
	}
	dbg.Info("   %d x%x bytes remaining", m.ID3Area(), m.ID3Area())

	return nil
}

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
		err = DumpID3v2Info(m)
		if nil != err {
			dbg.Error("ERR: %v", err)
		}
		m.Close()
	}
}
