package mp3

import (
	"bufio"
	"os"

	"github.com/jayacarlson/dbg"
	"github.com/jayacarlson/pth"
)

// Info on ID3Tag and ID3Frames:  	https://id3.org/id3v2.3.0
//									https://id3.org/id3v2.4.0-structure

var Dbg = dbg.Dbg{false, 0}

type (
	ID3v2Frame struct {
		unsync bool
		Tag    string
		Size   uint32 // data len size if id3v2.4
		Flags  [2]byte
		Data   []byte
		VerMaj byte
	}

	MP3File struct {
		filename string
		file     *os.File
		r        *bufio.Reader

		verMajor, verMinor byte
		bitFlags           byte
		id3Size            uint32
		extFlags           []byte // extended header flag bytes (2 for v2.3, N for v2.4)
	}
)

const (
	iso8859_1 byte = iota
	ucs2_BOM
	utf16BE
	utf8
)

// Opens a file, may or maynot be an MP3 file...
func OpenMP3File(filename string) (*MP3File, error) {
	path := pth.AsRealPath(filename)
	file, err := os.Open(path)
	if Dbg.ChkErr(err, "Failed to open file: "+path) {
		return nil, err
	}

	m := new(MP3File)
	m.filename = path
	m.file = file
	m.r = bufio.NewReader(file)
	return m, nil
}

func (m *MP3File) Close() {
	m.file.Close()
}

func (m *MP3File) Path() string     { return m.filename }
func (m *MP3File) ID3Area() int     { return int(m.id3Size) }
func (m *MP3File) ID3v2Flags() byte { return m.bitFlags }
func (m *MP3File) ID3v2Version() (int, int) {
	return int(m.verMajor), int(m.verMinor)
}

func (m *MP3File) FindFrame(f string) (ID3v2Frame, error) {
	err := m.ValidateID3v2File()
	if nil != err {
		return ID3v2Frame{}, err
	}
	if !(2 <= m.verMajor && 4 >= m.verMajor) {
		return ID3v2Frame{}, Err_UnknownID3v2Version
	}

	for {
		frm, err := m.ReadID3v2Frame()
		if nil != err {
			return ID3v2Frame{}, err
		}
		if frm.Tag == f {
			return frm, nil
		}
	}

	return ID3v2Frame{}, Err_EOF
}
