package disk

import (
	"encoding/binary"

	"github.com/xephonhq/xephon-k/pkg/util"
)

const (
	// Version is current supported file format version
	Version byte = 1
	// MagicNumber is 'xephon-k', 8 bytes stored in big endian in uint64, used for identify file without relying on extension
	MagicNumber uint64 = 0x786570686F6E2D6B
)

var log = util.Logger.NewEntryWithPkg("k.s.disk")

func MagicBytes() []byte {
	return []byte("xephon-k")
}

func IsMagic(buf []byte) bool {
	// binary.BigEndian.Uint64 don't check the length of the slice
	if len(buf) < 8 {
		return false
	}
	return binary.BigEndian.Uint64(buf) == MagicNumber
}
