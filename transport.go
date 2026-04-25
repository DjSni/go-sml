package sml

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

var (
	StartSeq = []byte{0x1b, 0x1b, 0x1b, 0x1b, 0x01, 0x01, 0x01, 0x01}
	EscSeq   = []byte{0x1b, 0x1b, 0x1b, 0x1b}
	EndSeq   = []byte{0x1b, 0x1b, 0x1b, 0x1b, 0x1a}
)

func TransportRead(r *bufio.Reader) ([]byte, error) {
	prefix := make([]byte, 0, len(StartSeq))

	for {
		b, err := r.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, io.EOF
			}
			return nil, err
		}

		if len(prefix) < len(StartSeq) {
			prefix = append(prefix, b)
		} else {
			copy(prefix, prefix[1:])
			prefix[len(prefix)-1] = b
		}

		if len(prefix) == len(StartSeq) && bytes.Equal(prefix, StartSeq) {
			frame := append([]byte(nil), StartSeq...)

			for {
				next, err := r.ReadByte()
				if err != nil {
					if errors.Is(err, io.EOF) {
						return nil, errors.New("premature eof")
					}
					return nil, err
				}
				frame = append(frame, next)

				if len(frame) < len(EscSeq) {
					continue
				}

				if !bytes.Equal(frame[len(frame)-len(EscSeq):], EscSeq) {
					continue
				}

				marker, err := r.ReadByte()
				if err != nil {
					if errors.Is(err, io.EOF) {
						return nil, errors.New("premature eof")
					}
					return nil, err
				}
				frame = append(frame, marker)

				if marker != EndSeq[len(EndSeq)-1] {
					continue
				}

				for i := 0; i < 3; i++ {
					trail, err := r.ReadByte()
					if err != nil {
						if errors.Is(err, io.EOF) {
							return nil, errors.New("premature eof")
						}
						return nil, err
					}
					frame = append(frame, trail)
				}

				return frame, nil
			}
		}
	}
}
