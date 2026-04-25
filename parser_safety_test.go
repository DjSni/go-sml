package sml

import (
	"strings"
	"testing"
)

func TestNumberParseDetectsTruncatedPayload(t *testing.T) {
	buf := &Buffer{Bytes: []byte{0x63, 0x01}}

	_, err := U16Parse(buf)
	if err == nil || !strings.Contains(err.Error(), "Unexpected end of buffer while parsing number") {
		t.Fatalf("expected truncated-number error, got %v", err)
	}
}

func TestOctetStringParseDetectsTruncatedPayload(t *testing.T) {
	buf := &Buffer{Bytes: []byte{0x03, 0x01}}

	_, err := OctetStringParse(buf)
	if err == nil || !strings.Contains(err.Error(), "Unexpected end of buffer while parsing octet string") {
		t.Fatalf("expected truncated-octet error, got %v", err)
	}
}

func TestBufGetCurrentByteOutOfBoundsReturnsZero(t *testing.T) {
	buf := &Buffer{Bytes: []byte{0x42}, Cursor: 3}

	if got := BufGetCurrentByte(buf); got != 0 {
		t.Fatalf("expected 0 for out-of-bounds cursor, got %x", got)
	}
}

func TestFileParseStopsOnNonListByte(t *testing.T) {
	messages, err := FileParse([]byte{0x60, 0x00, 0x00})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(messages) != 0 {
		t.Fatalf("expected no messages, got %d", len(messages))
	}
}
