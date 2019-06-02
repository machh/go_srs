package protocol

import (
	"errors"
	"bytes"
	"log"
	"encoding/binary"
)

type SrsStream struct {
	// current position at bytes.
	p []byte
	// the bytes data for stream to read or write.
	bytes []byte
	// the total number of bytes.
	n_bytes int32
}

func NewSrsStream(data []byte, len int32) *SrsStream {
	return &SrsStream{
		p:       data,
		bytes:   data,
		n_bytes: len,
	}
}

func (s *SrsStream) data() []byte {
	return s.bytes
}

func (s *SrsStream) size() int32 {
	return s.n_bytes
}

func (s *SrsStream) empty() bool {
	return s.bytes == nil || len(s.p) <= 0
}

func (s *SrsStream) require(required_size int32) bool {
	return int(required_size) <= len(s.p)
}

func (s *SrsStream) skip(size int32) {
	s.p = s.p[size:]
}

func (s *SrsStream) read_nbytes(n int32) (b []byte, err error) {
	if !s.require(n) {
		err = errors.New("no enough data")
		return
	}

	b = s.p[0 : n+1]
	s.skip(n)
	return
}

func (s *SrsStream) read_int16() (v int16, err error) {
	b, err := s.read_nbytes(2)
	if err != nil {
		return
	}

	bin_buf := bytes.NewBuffer(b)
	binary.Read(bin_buf, binary.LittleEndian, &v)
	return
}

func (s *SrsStream) read_int32() (v int32, err error) {
	b, err := s.read_nbytes(4)
	if err != nil {
		return
	}

	bin_buf := bytes.NewBuffer(b)
	log.Printf("********************%x %x %x %x", b[0], b[1], b[2], b[3])
	binary.Read(bin_buf, binary.BigEndian, &v)
	return
}

func (s *SrsStream) read_string(len int32) (str string, err error) {
	if !s.require(len) {
		err = errors.New("no enough data")
		return
	}

	str = string(s.p[:len+1])
	err = nil
	return
}

func (s *SrsStream) write_bytes(d []byte) {
	s.p = append(s.p, d...)
}

func (s *SrsStream) write_string(v string) {
	s.p = append(s.p, []byte(v)...)
}