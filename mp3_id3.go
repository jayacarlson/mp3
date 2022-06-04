package mp3

import (
	"github.com/jayacarlson/dbg"
)

/*
	ID3v2.3 Header:
		[3]byte		"ID3"
		byte, byte	majorVersion (3), minorVersion
		byte		flagBits
		[4]byte		size -- non-syncsafe!

	  Extended Header: (identified by flagBit)
		[4]byte		size -- non-syncsafe!
		[2]byte		flagBytes
		[4]byte		padding -- non-syncsafe!
		[ extra data ]

	ID3v2.4 Header:
		[3]byte		"ID3"
		byte, byte	majorVersion (4), minorVersion
		byte		flagBits
		[4]byte		size -- syncsafe bytes (4* %0nnnnnnn)

	  Extended Header: (identified by flagBit)
		[4]byte		size -- syncsafe bytes (4* %0nnnnnnn)
		byte		numFlagBytes
		[]byte		flagBytes
		[ extra data ]
*/

const (
	v2flagBit_Unsync         = 0x80
	v2flagBit_ExtendedHeader = 0x40
	v2flagBit_Experimental   = 0x20
	v2flagBit_Footer         = 0x10 // ID3v2.4
)

// Validates a file has ID3v2 tag at start, captures info & updates bufio buffer only if ID3v2 tag found
func (m *MP3File) ValidateID3v2File() error {
	if 0 != m.verMajor { // already parse the header
		return nil
	}

	skipSize := 10

	p, err := m.r.Peek(32)
	if nil != err {
		return err
	}

	if p[0] != 'I' || p[1] != 'D' || p[2] != '3' {
		return Err_NoID3v2Tag
	}

	m.verMajor = p[3]
	m.verMinor = p[4]
	if 0xFF == m.verMajor || 0xFF == m.verMinor {
		return Err_IllegalID3v2Header
	}
	m.bitFlags = p[5]

	sz, err := readInt(4 == m.verMajor, 4, p[6:])
	if nil != err {
		return err
	}
	m.id3Size = sz

	// any extended header?
	if 0 != (m.bitFlags & v2flagBit_ExtendedHeader) {
		sz, err = readInt(4 == m.verMajor, 4, p[10:])
		if nil != err {
			return err
		}
		skipSize += int(sz)
		m.id3Size -= sz // remove extended header size

		if 3 == m.verMajor {
			m.extFlags = make([]byte, 2)

			m.extFlags[0] = p[14]
			m.extFlags[1] = p[15]
		} else if 4 == m.verMajor {
			if 6 > sz {
				return Err_IllegalID3v2ExtHeader
			}

			if 1 != p[14] { // currently restricted to 1
				return Err_IllegalID3v2ExtHeader
			}

			// keep our options open...
			n := int(p[14])
			m.extFlags = make([]byte, n)
			for i := 0; i < n; i++ {
				m.extFlags[i] = p[15+i]
			}
		} else {
			dbg.Warning("%s", Err_UnknownID3v2Version)
		}
	}
	m.r.Discard(skipSize)

	return nil
}
