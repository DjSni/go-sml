package sml

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	StartSeq = []byte{0x1b, 0x1b, 0x1b, 0x1b, 1, 1, 1, 1}
	EscSeq   = []byte{0x1b, 0x1b, 0x1b, 0x1b}
	EndSeq   = []byte{0x1b, 0x1b, 0x1b, 0x1b, 0x1a}
)

// TransportRead finds and validates one version-1 frame and returns its
// original wire bytes. CRC, escape sequences, alignment and fill are checked.
func TransportRead(r *bufio.Reader) ([]byte, error) {
	prefix := make([]byte, 0, 8)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		prefix = append(prefix, b)
		if len(prefix) > 8 {
			prefix = prefix[1:]
		}
		if len(prefix) == 8 && bytes.Equal(prefix[:4], EscSeq) && prefix[4] == 2 {
			return nil, errors.New("SML transport version 2 is not supported")
		}
		if bytes.Equal(prefix, StartSeq) {
			break
		}
	}
	frame := append([]byte(nil), StartSeq...)
	run := 0
	for {
		b, err := r.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, errors.New("premature eof")
			}
			return nil, err
		}
		frame = append(frame, b)
		if b == 0x1b {
			run++
		} else {
			run = 0
		}
		if run != 4 {
			continue
		}
		var marker [4]byte
		if _, err = io.ReadFull(r, marker[:]); err != nil {
			return nil, errors.New("premature eof")
		}
		frame = append(frame, marker[:]...)
		if bytes.Equal(marker[:], EscSeq) {
			run = 0
			continue
		}
		if marker[0] != 0x1a {
			return nil, fmt.Errorf("Unsupported transport escape: % x", marker)
		}
		if _, err = TransportPayload(frame); err != nil {
			return nil, err
		}
		return frame, nil
	}
}

// TransportPayload validates a complete version-1 frame, removes its framing,
// unescapes the payload and removes precisely the declared fill bytes.
func TransportPayload(frame []byte) ([]byte, error) {
	if len(frame) < 16 || !bytes.Equal(frame[:8], StartSeq) || !bytes.Equal(frame[len(frame)-8:len(frame)-3], EndSeq) {
		return nil, errors.New("Invalid SML transport framing")
	}
	crcEnd := len(frame) - 2
	if Crc16Calculate(frame[:crcEnd], crcEnd) != binary.BigEndian.Uint16(frame[crcEnd:]) {
		return nil, errors.New("Transport CRC error")
	}
	wire := frame[8 : len(frame)-8]
	if len(wire)%4 != 0 {
		return nil, errors.New("Invalid transport payload alignment")
	}
	fill := int(frame[len(frame)-3])
	if fill > 3 || fill > len(wire) {
		return nil, errors.New("Invalid transport fill count")
	}
	for _, b := range wire[len(wire)-fill:] {
		if b != 0 {
			return nil, errors.New("Nonzero transport fill byte")
		}
	}
	wire = wire[:len(wire)-fill]
	payload := make([]byte, 0, len(wire))
	for i := 0; i < len(wire); {
		if len(wire)-i >= 4 && bytes.Equal(wire[i:i+4], EscSeq) {
			if len(wire)-i < 8 || !bytes.Equal(wire[i+4:i+8], EscSeq) {
				return nil, errors.New("Unescaped transport escape sequence")
			}
			payload = append(payload, EscSeq...)
			i += 8
		} else {
			payload = append(payload, wire[i])
			i++
		}
	}
	return payload, nil
}

// TransportParse validates both the transport CRC and every message CRC.
func TransportParse(frame []byte) ([]Message, error) {
	payload, err := TransportPayload(frame)
	if err != nil {
		return nil, err
	}
	return FileParse(payload)
}
