package disk

import (
	"encoding/binary"

	"github.com/xephonhq/xephon-k/pkg/util"
)

// NOTE: The following struct are generated by protobuf, see disk.proto and disk.pb.go for detail
// 	IndexEntries
// 	IndexEntry

const (
	// Version is current supported file format version
	Version byte = 1
	// MagicNumber is 'xephon-k', 8 bytes stored in big endian in uint64, used for identify file without relying on extension
	MagicNumber uint64 = 0x786570686F6E2D6B
)

var log = util.Logger.NewEntryWithPkg("k.storage.disk")

func MagicBytes() []byte {
	return []byte("xephon-k")
}

func IsMagic(buf []byte) bool {
	// binary.BigEndian.Uint64 don't check the indexLength of the slice
	if len(buf) < 8 {
		return false
	}
	return binary.BigEndian.Uint64(buf) == MagicNumber
}

func IsValidFormat(buf []byte) bool {
	if len(buf) < 9 {
		return false
	}
	return buf[0] == Version && IsMagic(buf[1:])
}
