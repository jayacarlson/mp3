package mp3

import (
	"bytes"
	"unicode/utf16"
)

func stringFromBytesSlice(b []byte) (string, []byte, error) {
	if 0 == len(b) {
		return "", b, nil
	}

	n := bytes.IndexByte(b, 0x00)
	if -1 == n {
		return string(b), []byte{}, nil
	}
	return string(b[:n]), b[n+1:], nil
}

func stringFromByteSliceOfUCS2BOM(b []byte) (string, []byte, error) {
	var i, o int
	var hi, lo int

	l := len(b)
	if Dbg.ChkTru(0 == l&1, "Byte len not even for Uint16 parsing") {
		return "", []byte{}, Err_IllegalFrameString
	}

	// check for NULL terminated string (may or maynot have BOM)
	if 0 == b[0] && 0 == b[1] {
		return "", b[2:], nil
	}
	// chk BOM
	if Dbg.ChkTru((l > 1) && ((b[0] == 0xFF && b[1] == 0xFE) || (b[1] == 0xFF && b[0] == 0xFE)), "Illegal BOM") {
		return "", []byte{}, Err_IllegalFrameString
	}
	if b[0] == 0xFF {
		hi, lo = 1, 0
	} else {
		hi, lo = 0, 1
	}

	l -= 2
	b = b[2:]

	r := make([]uint16, l/2)

	for i = 0; i < l; {
		u := (uint16(b[i+hi]) << 8) | (uint16(b[i+lo]))
		i += 2
		if 0 == u {
			r = r[:o]
			break
		}
		r[o] = u
		o++
	}
	return string(utf16.Decode(r)), b[i:], nil
}

func stringFromByteSliceOfUint16BE(b []byte) (string, []byte, error) {
	l := len(b)
	if Dbg.ChkTru(0 == l&1, "Byte len not even for Uint16 parsing") {
		return "", []byte{}, Err_IllegalFrameString
	}

	r := make([]uint16, l/2)
	var i, o int

	// data is in big-endian form
	for i = 0; i < l; {
		u := (uint16(b[i]) << 8) | (uint16(b[i+1]))
		i += 2
		if 0 == u {
			r = r[:o]
			break
		}
		r[o] = u
		o++
	}
	return string(utf16.Decode(r)), b[i:], nil
}

func readID3StringByType(e byte, b []byte) (string, error) {
	switch e {
	case iso8859_1, utf8:
		s, _, err := stringFromBytesSlice(b)
		return s, err
	case ucs2_BOM:
		s, _, err := stringFromByteSliceOfUCS2BOM(b)
		return s, err
	case utf16BE:
		s, _, err := stringFromByteSliceOfUint16BE(b)
		return s, err
	}
	return "", Err_IllegalFrameStringType
}

func readID3StringsByType(e byte, b []byte) ([]string, error) {
	r := []string{}
	switch e {
	case iso8859_1, utf8, ucs2_BOM, utf16BE:
	default:
		return r, Err_IllegalFrameStringType
	}
	var s string
	var err error

	for 0 != len(b) {
		switch e {
		case iso8859_1, utf8:
			s, b, err = stringFromBytesSlice(b)
		case ucs2_BOM:
			s, b, err = stringFromByteSliceOfUCS2BOM(b)
		case utf16BE:
			s, b, err = stringFromByteSliceOfUint16BE(b)
		}
		if nil != err {
			break
		}
		r = append(r, s)
	}
	return r, err
}
