package tcpheader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// ErrUnknownHeader ...
var ErrUnknownHeader error = fmt.Errorf("unknown header")

// HeaderType ...
type HeaderType int

const (
	// Unknown ...
	Unknown = iota
	// TwoByteUnsigned ...
	TwoByteUnsigned
	// FourByteUnsigned ...
	FourByteUnsigned
)

func (h HeaderType) String() string {
	names := [...]string{
		"Unknown type",
		"twoByteUnsigned",
		"fourByteUnsigned",
	}

	if h < Unknown || h > FourByteUnsigned {
		return "Unknown type"
	}

	return names[h]
}

// ReadLen ...
func ReadLen(r io.Reader, h HeaderType) (int, error) {
	switch h {
	case TwoByteUnsigned:
		var header uint16

		err := binary.Read(r, binary.BigEndian, &header)
		if nil != err {
			return 0, err
		}

		return int(header), nil
	case FourByteUnsigned:
		var header uint32

		err := binary.Read(r, binary.BigEndian, &header)
		if nil != err {
			return 0, err
		}

		return int(header), nil
	}

	return 0, ErrUnknownHeader
}

// ReadMessage ...
func ReadMessage(r io.Reader, p []byte) error {
	for len(p) > 0 {
		n, err := r.Read(p)
		p = p[n:]

		if err != nil {
			return err
		}
	}

	return nil
}

// WriteMessage ...
func WriteMessage(r io.Writer, p []byte, h HeaderType) error {
	var err error

	//Build a temp buffer with full message
	buf := new(bytes.Buffer)

	//Add BinEndian length indicator to send buffer
	switch h {
	case TwoByteUnsigned:
		err = binary.Write(buf, binary.BigEndian, uint16(len(p)))
	case FourByteUnsigned:
		err = binary.Write(buf, binary.BigEndian, uint32(len(p)))
	}

	if nil != err {
		return err
	}

	//Add message to buffer
	if _, err = buf.Write(p); nil != err {
		return err
	}

	//Write buffer
	if _, err = r.Write(buf.Bytes()); nil != err {
		return err
	}

	return nil
}

// ReadC return a channel unbuffered channel
// that returns bytes read from reader
func ReadC(r io.Reader, header HeaderType) chan []byte {
	output := make(chan []byte)
	go func() {
		defer close(output)
		for {
			l, err := ReadLen(r, header)
			if err != nil {
				return
			}

			buf := make([]byte, l)
			if err = ReadMessage(r, buf); err != nil {
				return
			}
			output <- buf
		}
	}()

	return output
}
