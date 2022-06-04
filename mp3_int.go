package mp3

// TODO_TODO_TODO: Add 'unsynchronisation' logic for the non-SyncSafe readers
func readInt(syncSafe bool, s int, p []byte) (uint32, error) {
	var v uint32
	if syncSafe {
		for i := 0; s > i; i++ {
			b := p[i]
			if 0 != (b & 0x80) {
				return 0, Err_InvalidSyncSafe
			}
			v = (v << 7) + uint32(b)
		}
	} else {
		for i := 0; s > i; i++ {
			b := p[i]
			v = (v << 8) + uint32(b)
		}
	}
	return v, nil
}

func (m *MP3File) readInt(syncSafe bool, s int) (uint32, error) {
	var v uint32
	if syncSafe {
		for i := 0; s > i; i++ {
			b, err := m.r.ReadByte()
			if nil != err {
				return 0, err
			}
			if 0 != (b & 0x80) {
				return 0, Err_InvalidSyncSafe
			}
			v = (v << 7) + uint32(b)
		}
	} else {
		for i := 0; s > i; i++ {
			b, err := m.r.ReadByte()
			if nil != err {
				return 0, err
			}
			v = (v << 8) + uint32(b)
		}
	}
	return v, nil
}
