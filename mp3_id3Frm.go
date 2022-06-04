package mp3

const (
	frm4StatFlg_TagPres   = 0x40
	frm4StatFlg_FilePres  = 0x20
	frm4StatFlg_ReadOnly  = 0x10
	frm4FmtFlg_Grouping   = 0x40
	frm4FmtFlg_Compressed = 0x08
	frm4FmtFlg_Encrypted  = 0x04
	frm4FmtFlg_Unsync     = 0x02
	frm4FmtFlg_DataLen    = 0x01

	frm3StatFlg_TagPres   = 0x80
	frm3StatFlg_FilePres  = 0x40
	frm3StatFlg_ReadOnly  = 0x20
	frm3FmtFlg_Compressed = 0x80
	frm3FmtFlg_Encrypted  = 0x40
	frm3FmtFlg_Grouping   = 0x20
)

func invalidFrmChar(b byte) bool {
	if b >= 'A' && b <= 'Z' {
		return false
	}
	if b >= '0' && b <= '9' {
		return false
	}
	return true
}

func InvalidID3v2FrameID(s string) bool {
	return (4 != len(s)) || invalidFrmChar(s[0]) || invalidFrmChar(s[1]) || invalidFrmChar(s[2]) || invalidFrmChar(s[3])
}

func (m *MP3File) ReadID3v2Frame() (ID3v2Frame, error) {
	frm := ID3v2Frame{}
	if 10 > m.id3Size { // large enough to read a frame header?
		return frm, Err_EOF
	}
	p, err := m.r.Peek(10)
	if nil != err {
		return frm, err
	}

	if invalidFrmChar(p[0]) || invalidFrmChar(p[1]) || invalidFrmChar(p[2]) || invalidFrmChar(p[3]) {
		return frm, Err_EOF
	}

	sz, err := readInt(4 == m.verMajor, 4, p[4:])
	if nil != err {
		return frm, err
	}
	if 0 == sz { // frame must be at least 1 byte
		return frm, Err_IllegalFrame
	}
	// id3 tag not large enough for data
	if sz+10 > m.id3Size {
		return frm, Err_IllegalFrame
	}

	frm.Tag = string(p[0:4])
	frm.Size = sz
	frm.Flags[0] = p[8]
	frm.Flags[1] = p[9]
	frm.VerMaj = m.verMajor

	m.id3Size -= (sz + 10)

	m.r.Discard(10)

	if (4 == m.verMajor) && (0 != frm.Flags[1]&frm4FmtFlg_DataLen) {
		if 4 > frm.Size { // must have 4 bytes for DataLen field
			return ID3v2Frame{}, Err_IllegalFrameData
		}
		sz, err := m.readInt(true, 4)
		if nil != err {
			return ID3v2Frame{}, err
		}
		frm.Data = make([]byte, frm.Size-4)
		n, err := m.r.Read(frm.Data)
		if nil != err {
			return ID3v2Frame{}, err
		}
		if n != int(frm.Size-4) {
			return ID3v2Frame{}, Err_EOF // maybe a better error?  an OS error?
		}
		frm.Size = sz
	} else {
		frm.Data = make([]byte, frm.Size)
		n, err := m.r.Read(frm.Data)
		if nil != err {
			return ID3v2Frame{}, err
		}
		if n != int(frm.Size) {
			return ID3v2Frame{}, Err_EOF // maybe a better error?  an OS error?
		}
	}

	return frm, nil
}
