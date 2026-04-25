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
	fcsAndFill := []byte{0x12, 0x34, 0x00}

	input := append([]byte{0x00, 0x11, 0x22}, StartSeq...)
	input = append(input, payload...)
	input = append(input, EndSeq...)
	input = append(input, fcsAndFill...)

	got, err := TransportRead(bufio.NewReader(bytes.NewReader(input)))
	if err != nil {
		t.Fatalf("TransportRead returned error: %v", err)
	}

	want := append([]byte{}, StartSeq...)
	want = append(want, payload...)
	want = append(want, EndSeq...)
	want = append(want, fcsAndFill...)

	if !bytes.Equal(got, want) {
		t.Fatalf("unexpected frame bytes: got % x, want % x", got, want)
	}
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
