package mp3

import (
	"github.com/jayacarlson/dbg"
)

// Info on MP3 file header:  	https://en.wikipedia.org/wiki/MP3

func isMP3Header(b []byte) bool {
	return 0xFF == b[0] && 0xF0 <= b[1]
}

func (m *MP3File) showMP3Header() (uint32, error) {
	p, err := m.r.Peek(4)
	if nil != err {
		return 0, err
	}
	dbg.Warning("  FFF %x(v%d l%d e%d) %x(br%d) %x(f%d p%d X) %x(m%d i%d s%d) %x(c%d o%d E%d)",
		(p[1] & 0xF),   // ver/layer/ep
		(p[1]&0x8)>>3,  //   version
		(p[1]&0x6)>>1,  //   layer
		(p[1] & 0x1),   //   error protection
		(p[2]&0xF0)>>4, // bitrate (hex)
		(p[2]&0xF0)>>4, //   bitrate (int)
		(p[2] & 0xF),   // freq/padb
		(p[2]&0xC)>>2,  //   frequency
		(p[2]&0x2)>>1,  //   padded
		//(p[2] & 0x1), //   reserved
		(p[3]&0xF0)>>4, // mode/int/ms
		(p[3]&0xC0)>>6, //   mode
		(p[3]&0x20)>>5, //   intensity stereo
		(p[3]&0x10)>>4, //   MS stereo
		(p[3] & 0xF),   // cpywrt/orig/emph
		(p[3]&0x8)>>3,  //   copywrite
		(p[3]&0x4)>>2,  //   original
		(p[3] & 0x3))   //   emphesis

	return 0, nil
}
