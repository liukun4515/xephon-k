package disk

import (
	"os"

	"syscall"

	"encoding/binary"

	"github.com/pkg/errors"
	"github.com/xephonhq/xephon-k/pkg/common"
)

var _ DataFileReader = (*LocalDataFileReader)(nil)

type DataFileReader interface {
	ReadIndexOfIndexes() error
	SeriesCount() int
	Close() error
}

type IndexOfIndex struct {
	offset uint32
	length uint32
}

type LocalDataFileReader struct {
	f                  *os.File
	fStat              os.FileInfo
	b                  []byte
	size               int
	indexOffset        uint64
	indexOfIndexOffset uint32
	indexLength        uint32
	indexOfIndexes     map[common.SeriesID]IndexOfIndex
}

func NewLocalDataFileReader(f *os.File) (*LocalDataFileReader, error) {
	name := f.Name()
	// we check version after we mmap the file, because normally, it should be the file we need
	stat, err := f.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "can't get file stat")
	}
	// mmap the file
	// https://github.com/golang/exp/blob/master/mmap/mmap_unix.go
	size := stat.Size()
	if size == 0 {
		return nil, errors.Errorf("mmap: file %s is empty", name)
	}
	if size < 0 {
		return nil, errors.Errorf("mmap: file %s has negative size %d", name, size)
	}
	if size != int64(int(size)) {
		return nil, errors.Errorf("mmap: file %s is too large, it's likely you have a file larger than 4GB on a 32 bit OS", name)
	}
	if size < FooterLength {
		return nil, errors.Errorf("file is too short, file is %d bytes, footer is %d bytes", size, FooterLength)
	}
	b, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, errors.Errorf("mmap: file %s can't be mmaped", name)
	}

	r := &LocalDataFileReader{
		f:     f,
		fStat: stat,
		b:     b,
		size:  int(size),
	}

	footer := b[size-FooterLength:]

	// read index position
	indexOffset := binary.BigEndian.Uint64(footer[:8])
	indexOfIndexOffset := binary.BigEndian.Uint32(footer[8:12])
	indexLength := binary.BigEndian.Uint32(footer[12:16])
	if uint64(indexLength) != (uint64(size) - indexOffset - uint64(FooterLength)) {
		// unmap and close the file
		if err := r.Close(); err != nil {
			return nil, errors.Wrap(err, "can't close reader after invalid index length is detected")
		}
		return nil, errors.Errorf("indexLength %d != (size %d - indexOffset %d - FooterLength %d)",
			indexLength, size, indexOffset, FooterLength)
	}
	if indexOfIndexOffset > indexLength {
		// unmap and close the file
		if err := r.Close(); err != nil {
			return nil, errors.Wrap(err, "can't close reader after invalid index of index offset is detected")
		}
		return nil, errors.Errorf("index of index offset %d is larger than total index length %d", indexOfIndexOffset, indexLength)
	}
	r.indexOffset = indexOffset
	r.indexOfIndexOffset = indexOfIndexOffset
	r.indexLength = indexLength

	// check version
	if !IsValidFormat(footer[16:]) {
		// unmap and close the file
		if err := r.Close(); err != nil {
			return nil, errors.Wrap(err, "can't close reader after invalid format is detected")
		}
		return nil, errors.Errorf("version and/or magic does not match, expected %v %d but got %v %d", Version, MagicNumber, b[size-9], b[size-8:])
	}

	return r, nil
}

func (reader *LocalDataFileReader) ReadIndexOfIndexes() error {
	if reader.indexOfIndexes != nil {
		// TODO: return error if called multiple times? currently we just silently return
		return nil
	}

	seriesCount := int((reader.indexLength - reader.indexOfIndexOffset) / (IndexOfIndexUnitLength))
	reader.indexOfIndexes = make(map[common.SeriesID]IndexOfIndex, seriesCount)
	log.Infof("size %d idx offset %d idx of idx offset %d length %d series count %d",
		reader.size, reader.indexOffset, reader.indexOfIndexOffset, reader.size, seriesCount)
	// load all the needed bytes
	start := reader.indexOffset + uint64(reader.indexOfIndexOffset)
	b := reader.b[start : start+uint64(seriesCount*IndexOfIndexUnitLength)]

	var (
		id     uint64
		offset uint32
		length uint32
	)
	for i := 0; i < seriesCount; i++ {
		id = binary.BigEndian.Uint64(b[i*IndexOfIndexUnitLength : i*IndexOfIndexUnitLength+8])
		offset = binary.BigEndian.Uint32(b[i*IndexOfIndexUnitLength+8 : i*IndexOfIndexUnitLength+12])
		length = binary.BigEndian.Uint32(b[i*IndexOfIndexUnitLength+12 : i*IndexOfIndexUnitLength+16])
		reader.indexOfIndexes[common.SeriesID(id)] = IndexOfIndex{
			offset: offset,
			length: length,
		}
	}
	return nil
}

func (reader *LocalDataFileReader) SeriesCount() int {
	if reader.indexOfIndexes == nil {
		return 0
	} else {
		return len(reader.indexOfIndexes)
	}
}

func (reader *LocalDataFileReader) Close() error {
	// the reader is not initialized or already closed
	if reader.b == nil {
		return nil
	}
	if err := syscall.Munmap(reader.b); err != nil {
		return errors.Wrapf(err, "mmap: can't unmap file %s", reader.f.Name())
	}
	reader.b = nil
	if err := reader.f.Close(); err != nil {
		return errors.Wrapf(err, "can't close file %s after unmap", reader.f.Name())
	}
	return nil
}
