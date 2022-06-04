package main

import (
	"flag"
	"io/ioutil"
	"path"

	"github.com/jayacarlson/dbg"
	"github.com/jayacarlson/mp3"
	"github.com/jayacarlson/pth"
)

var (
	dprog, fprog bool
)

func init() {
	flag.BoolVar(&dprog, "d", false, "Show progress by outputing directory")
	flag.BoolVar(&fprog, "f", false, "Show progress by outputing file")
}

func processDir(srcPath string, recurse bool, dirHandler func(string) error, filHandler func(string) error) error {
	if !pth.Exists(srcPath) {
		return pth.Err_NotExist
	}

	if dprog {
		dbg.Info("processing: %s", srcPath)
	}

	entries, err := ioutil.ReadDir(srcPath)
	err = pth.ChkFileErr(err)
	if nil != err {
		return err
	}
	for _, entry := range entries {
		if !entry.Mode().IsDir() && (nil != filHandler) {
			err = filHandler(path.Join(srcPath, entry.Name()))
			err = pth.ChkFileErr(err)
			if nil != err {
				return err
			}
			continue
		}
		if entry.Mode().IsDir() {
			if nil != dirHandler {
				err = dirHandler(path.Join(srcPath, entry.Name()))
				if nil != err {
					return err
				}
			}
			if recurse {
				err = processDir(path.Join(srcPath, entry.Name()), recurse, dirHandler, filHandler)
				if nil != err {
					return err
				}
			}
		}
	}
	return nil
}

func fileHandler(path string) error {
	_, _, ext := pth.Split(path)
	if fprog {
		dbg.Echo("checking: %s", path)
	}
	if ".mp3" != ext {
		return nil
	}
	m, err := mp3.OpenMP3File(path)
	if nil != err {
		dbg.Error("ERR: %v for file `%s`", err, path)
		return nil
	}
	defer m.Close()

	err = m.ValidateID3v2File()
	if nil != err {
		dbg.Error("ERR: %v for file `%s`", err, path)
		return nil
	}
	verMajor, _ := m.ID3v2Version()
	if !(2 <= verMajor && 4 >= verMajor) {
		dbg.Error("ERR: %v for file `%s`", mp3.Err_UnknownID3v2Version, path)
		return nil
	}
	if 2 == verMajor {
		dbg.Error("ERR: %v for file `%s`", mp3.Err_ObsoleteID3v2Version, path)
		return nil
	}

	for {
		frm, err := m.ReadID3v2Frame()
		if mp3.Err_EOF == err {
			return nil
		}
		if nil != err {
			dbg.Error("ERR: %v", err)
			return nil
		}
		_, err = frm.ToString()
		if nil != err {
			dbg.Error("ERR: %v", err)
			return nil
		}
	}

	return nil
}

func main() {
	flag.Parse()

	dirs := flag.Args()
	if len(dirs) == 0 {
		dbg.Fatal("Must supply dir to validate")
	}

	for _, dir := range dirs {
		err := processDir(pth.AsRealPath(dir), true, nil, fileHandler)
		if nil != err {
			return
		}
	}
}
