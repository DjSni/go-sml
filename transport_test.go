package sml

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestTransportReadFindsFrameAfterNoise(t *testing.T) {
	payload := []byte{0xaa, 0xbb, 0xcc}
	want := buildTransportFrame(payload)
	input := append([]byte{0x00, 0x11, 0x22}, want...)

	got, err := TransportRead(bufio.NewReader(bytes.NewReader(input)))
	if err != nil {
		t.Fatalf("TransportRead returned error: %v", err)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("unexpected frame bytes: got % x, want % x", got, want)
	}
}

func buildTransportFrame(payload []byte) []byte {
	wire := bytes.ReplaceAll(payload, EscSeq, append(append([]byte{}, EscSeq...), EscSeq...))
	fill := (4 - len(wire)%4) % 4
	frame := append([]byte{}, StartSeq...)
	frame = append(frame, wire...)
	frame = append(frame, make([]byte, fill)...)
	frame = append(frame, EndSeq...)
	frame = append(frame, byte(fill))
	crc := Crc16Calculate(frame, len(frame))
	return append(frame, byte(crc>>8), byte(crc))
}

func TestTransportReadReturnsEOFWithoutStartSequence(t *testing.T) {
	_, err := TransportRead(bufio.NewReader(bytes.NewReader([]byte{0x01, 0x02, 0x03})))
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected io.EOF, got %v", err)
	}
}

func TestTransportReadReturnsPrematureEOFForTruncatedFrame(t *testing.T) {
	input := append([]byte{}, StartSeq...)
	input = append(input, 0xaa, 0xbb)
	input = append(input, EscSeq...)
	input = append(input, 0x1a)
	input = append(input, 0x01)

	_, err := TransportRead(bufio.NewReader(bytes.NewReader(input)))
	if err == nil || err.Error() != "premature eof" {
		t.Fatalf("expected premature eof, got %v", err)
	}
}
