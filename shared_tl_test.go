package sml

import "testing"

func TestBufGetNextLengthForNonListSingleTL(t *testing.T) {
	buf := &Buffer{Bytes: []byte{0x63}}

	got := BufGetNextLength(buf)
	if got != 2 {
		t.Fatalf("unexpected length: got %d want 2", got)
	}
	if buf.Cursor != 1 {
		t.Fatalf("unexpected cursor: got %d want 1", buf.Cursor)
	}
}

func TestBufGetNextLengthForListSingleTL(t *testing.T) {
	buf := &Buffer{Bytes: []byte{0x72}}

	got := BufGetNextLength(buf)
	if got != 2 {
		t.Fatalf("unexpected length: got %d want 2", got)
	}
	if buf.Cursor != 1 {
		t.Fatalf("unexpected cursor: got %d want 1", buf.Cursor)
	}
}

func TestBufGetNextLengthForNonListMultiTL(t *testing.T) {
	buf := &Buffer{Bytes: []byte{0x81, 0x03}}

	got := BufGetNextLength(buf)
	if got != 17 {
		t.Fatalf("unexpected length: got %d want 17", got)
	}
	if buf.Cursor != 2 {
		t.Fatalf("unexpected cursor: got %d want 2", buf.Cursor)
	}
}

func TestBufGetNextLengthForListMultiTL(t *testing.T) {
	buf := &Buffer{Bytes: []byte{0xF1, 0x03}}

	got := BufGetNextLength(buf)
	if got != 19 {
		t.Fatalf("unexpected length: got %d want 19", got)
	}
	if buf.Cursor != 2 {
		t.Fatalf("unexpected cursor: got %d want 2", buf.Cursor)
	}
}

func TestBufOptionalIsSkippedAdvancesCursor(t *testing.T) {
	buf := &Buffer{Bytes: []byte{OPTIONALSKIPPED, 0x42}}

	if !BufOptionalIsSkipped(buf) {
		t.Fatalf("expected optional field to be skipped")
	}
	if buf.Cursor != 1 {
		t.Fatalf("unexpected cursor: got %d want 1", buf.Cursor)
	}
	if BufOptionalIsSkipped(buf) {
		t.Fatalf("did not expect second byte to be treated as optional marker")
	}
}
