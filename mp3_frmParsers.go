package mp3

import (
	"fmt"
	"strconv"
	"time"
)

// Parsers based off of:	https://id3.org/id3v2.3.0#ID3v2_header

/*
	<Header for 'Unique file identifier', ID: "UFID">
	Owner identifier    <text string> $00
	Identifier          <up to 64 bytes binary data>
*/
func (frm ID3v2Frame) ufidParser() (string, error) {
	o, b, err := stringFromBytesSlice(frm.Data)
	if nil != err {
		return "", err
	}
	if 64 < len(b) {
		return "", Err_IllegalFrameData
	}
	return fmt.Sprintf("%s [%v]", o, b), nil
}

/*
	<Header for 'Comment', ID: "COMM">
	Text encoding           $xx
	Language                $xx xx xx
	Short content descrip.  <text string according to encoding> $00 (00)
	The actual text         <full text string according to encoding>
*/
func (frm ID3v2Frame) commParser() (string, error) {
	b := frm.Data
	s, err := readID3StringsByType(b[0], b[4:])
	if nil != err {
		return "", err
	}
	if 2 != len(s) {
		return "", Err_IllegalFrameData
	}
	lang := string(b[1:4])
	if "" != s[0] {
		return fmt.Sprintf("(%s) [%s] %s", lang, s[0], s[1]), nil
	}
	return fmt.Sprintf("(%s) %s", lang, s[1]), nil
}

//	The 'Length' frame contains the length of the audiofile in milliseconds, represented as a numeric string.
func (frm ID3v2Frame) tlenParser() (string, error) {
	b := frm.Data
	s, err := readID3StringByType(b[0], b[1:])
	if nil != err {
		return "", err
	}
	// convert from milliseconds to HH:MM:SS.s
	millis, err := strconv.Atoi(s)
	if nil != err {
		return s, err
	}
	tm := time.Duration(millis) * time.Millisecond
	hr := tm.Truncate(time.Hour) / time.Hour
	tm -= hr * time.Hour
	mn := tm.Truncate(time.Minute) / time.Minute
	tm -= mn * time.Minute
	sc := tm.Truncate(time.Second) / time.Second
	tm -= sc * time.Second
	tm /= (time.Millisecond * 100)

	if hr > 0 {
		return fmt.Sprintf("(%s) %02d:%02d:%02d", s, hr, mn, sc), nil
	}
	return fmt.Sprintf("(%s) %02d:%02d.%d", s, mn, sc, tm), nil
}

/*
	<Header for 'User defined text information frame', ID: "TXXX">
	Text encoding    $xx
	Description      <text string according to encoding> $00 (00)
	Value            <text string according to encoding>
*/
func (frm ID3v2Frame) txxxParser() (string, error) {
	b := frm.Data
	s, err := readID3StringsByType(b[0], b[1:])
	if nil != err {
		return "", err
	}
	if 2 != len(s) {
		return "", Err_IllegalFrameData
	}
	if "" != s[0] {
		return fmt.Sprintf("[%s] %s", s[0], s[1]), nil
	}
	return fmt.Sprintf("%s", s[1]), nil
}

/*
	<Header for 'Text information frame', ID: "T000" - "TZZZ", excluding "TXXX" described in 4.2.2.>
	<Header for 'Involved people list', ID: "IPLS">
	Text encoding    $xx
	Information      <text string according to encoding>
*/
func (frm ID3v2Frame) t___Parser() (string, error) {
	b := frm.Data
	return readID3StringByType(b[0], b[1:])
}

/*
	<Header for 'URL link frame', ID: "W000" - "WZZZ", excluding "WXXX" described in 4.3.2.>
	URL <text string>
*/
func (frm ID3v2Frame) w___Parser() (string, error) {
	return readID3StringByType(iso8859_1, frm.Data)
}

/*
	<Header for 'User defined URL link frame', ID: "WXXX">
	Text encoding    $xx
	Description    <text string according to encoding> $00 (00)
	URL    <text string>
*/
func (frm ID3v2Frame) wxxxParser() (string, error) {
	b := frm.Data
	s, err := readID3StringsByType(b[0], b[1:])
	if nil != err {
		return "", err
	}
	if 2 != len(s) {
		return "", Err_IllegalFrameData
	}
	if "" != s[0] {
		return fmt.Sprintf("[%s] %s", s[0], s[1]), nil
	}
	return fmt.Sprintf("[---] %s", s[1]), nil
}

func (frm ID3v2Frame) ToString() (string, error) {
	if (4 == frm.VerMaj && (0 != frm.Flags[1]&(frm4FmtFlg_Compressed|frm4FmtFlg_Encrypted))) ||
		(3 == frm.VerMaj && (0 != frm.Flags[1]&(frm3FmtFlg_Compressed|frm3FmtFlg_Encrypted))) {
		return "<Encrypted / Compompressed data>", nil
	}

	if 'T' == frm.Tag[0] {
		if "TXXX" == frm.Tag {
			return frm.txxxParser()
		} else if "TLEN" == frm.Tag {
			return frm.tlenParser()
		}
		return frm.t___Parser()
	}
	if 'W' == frm.Tag[0] {
		if "WXXX" == frm.Tag {
			return frm.wxxxParser()
		}
		return frm.w___Parser()
	}

	switch frm.Tag {
	case "COMM":
		return frm.commParser()
		// other one off parsers...
	case "IPLS":
		return frm.t___Parser()
	}
	return fmt.Sprintf("Unrecognized frame <<%s>>", frm.Tag), nil
}
