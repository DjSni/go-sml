package sml

import (
	"strings"
	"testing"
)

func TestGetListResponseParseWithEmptyValueList(t *testing.T) {
	buf := &Buffer{Bytes: []byte{0x77, 0x01, 0x01, 0x01, 0x01, 0x70, 0x01, 0x01}}

	msg, err := GetListResponseParse(buf)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if msg.ValList == nil {
		t.Fatalf("expected empty value list, got nil")
	}
	if len(msg.ValList) != 0 {
		t.Fatalf("expected empty value list, got %d entries", len(msg.ValList))
	}
	if buf.Cursor != len(buf.Bytes) {
		t.Fatalf("parser consumed %d bytes, expected %d", buf.Cursor, len(buf.Bytes))
	}
}

func TestGetListResponseParseWithSingleOptionalEntry(t *testing.T) {
	listEntry := []byte{0x77, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}
	data := []byte{0x77, 0x01, 0x01, 0x01, 0x01, 0x71}
	data = append(data, listEntry...)
	data = append(data, 0x01, 0x01)
	buf := &Buffer{Bytes: data}

	msg, err := GetListResponseParse(buf)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(msg.ValList) != 1 {
		t.Fatalf("expected one list entry, got %d", len(msg.ValList))
	}
	if msg.ValList[0].ObjName != nil {
		t.Fatalf("expected optional object name to be nil")
	}
	if buf.Cursor != len(data) {
		t.Fatalf("parser consumed %d bytes, expected %d", buf.Cursor, len(data))
	}
}

func TestGetListResponseParseRejectsInvalidValueListType(t *testing.T) {
	buf := &Buffer{Bytes: []byte{0x77, 0x01, 0x01, 0x01, 0x01, 0x63, 0x01, 0x01, 0x01}}

	_, err := GetListResponseParse(buf)
	if err == nil || !strings.Contains(err.Error(), "Unexpected type") {
		t.Fatalf("expected unexpected-type error, got %v", err)
	}
}
