package mp3

import "errors"

var (
	Err_NoID3v2Tag             = errors.New("No ID3v2 Tag found")
	Err_EOF                    = errors.New("End of frames")
	Err_InvalidSyncSafe        = errors.New("Non-SyncSafe byte in integer")
	Err_IllegalID3v2Header     = errors.New("Illegal ID3v2 header data")
	Err_IllegalID3v2ExtHeader  = errors.New("Illegal ID3v2 extended header data")
	Err_UnknownID3v2Version    = errors.New("Unknown ID3v2 version")
	Err_ObsoleteID3v2Version   = errors.New("Obsolete ID3v2 version")
	Err_IllegalFrame           = errors.New("Illegal frame header")
	Err_IllegalFrameData       = errors.New("Illegal/Unknown frame data")
	Err_IllegalFrameString     = errors.New("Illegal frame string")
	Err_IllegalFrameStringType = errors.New("Illegal frame string type")
)
