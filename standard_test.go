package sml

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"math"
	"testing"
)

func TestStandardNumbers(t *testing.T) {
	for _, tc := range []struct {
		hex  string
		typ  byte
		size int
		want int64
	}{
		{"52ff", TYPEINTEGER, 8, -1},
		{"538000", TYPEINTEGER, 8, -32768},
		{"52fe", TYPEINTEGER, 4, -2},
		{"6280", TYPEUNSIGNED, 8, 128},
		{"63ffff", TYPEUNSIGNED, 2, 65535},
		{"65ffffffff", TYPEUNSIGNED, 4, 4294967295},
	} {
		t.Run(tc.hex, func(t *testing.T) {
			data, _ := hex.DecodeString(tc.hex)
			got, err := NumberParse(&Buffer{Bytes: data}, tc.typ, tc.size)
			if err != nil || got != tc.want {
				t.Fatalf("got %d, %v; want %d", got, err, tc.want)
			}
		})
	}
	v, err := ValueParse(&Buffer{Bytes: []byte{0x69, 255, 255, 255, 255, 255, 255, 255, 255}})
	if err != nil || v.DataUnsigned != math.MaxUint64 {
		t.Fatalf("unsigned64 lost: %+v, %v", v, err)
	}
	status, err := StatusParse(&Buffer{Bytes: []byte{0x65, 0, 8, 1, 4}})
	if err != nil || status != 0x00080104 {
		t.Fatalf("status lost: %x %v", status, err)
	}
}

func TestStandardTimeChoices(t *testing.T) {
	data := []byte{0x72, 0x62, 3, 0x73, 0x65, 0, 0, 0, 42, 0x53, 0xff, 0xc4, 0x52, 60}
	got, err := TimeParse(&Buffer{Bytes: data})
	if err != nil || got.Tag != 3 || got.Timestamp != 42 || got.LocalOffset != -60 || got.SeasonTimeOffset != 60 {
		t.Fatalf("local timestamp: %+v, %v", got, err)
	}
	v, err := ValueParse(&Buffer{Bytes: append([]byte{0x72, 0x62, 1}, data...)})
	if err != nil || v.DataTime == nil || *v.DataTime != got {
		t.Fatalf("SML_ListType: %+v %v", v, err)
	}
	if _, err := TimeParse(&Buffer{Bytes: []byte{0x72, 0x62, 4, 0x62, 0}}); err == nil {
		t.Fatal("unknown time choice accepted")
	}
}

func TestStandardTLRejectsInvalidHeaders(t *testing.T) {
	for _, data := range [][]byte{{0x80}, {0x81, 0x43}, {0x10}, {0xc2, 0x01}, {0x60}} {
		if _, err := ValueParse(&Buffer{Bytes: data}); err == nil {
			t.Fatalf("invalid TL accepted: % x", data)
		}
	}
	buf := &Buffer{Bytes: []byte{0x81, 0x80, 0x03}}
	if n := BufGetNextLength(buf); n != 256 {
		t.Fatalf("multi-byte TL length wrapped: %d", n)
	}
	for _, b := range []byte{0, 1, 0x80, 255} {
		v, err := ValueParse(&Buffer{Bytes: []byte{0x42, b}})
		if err != nil || v.Typ != TYPEBOOLEAN || v.DataBoolean != (b != 0) {
			t.Fatalf("Boolean %x: %+v %v", b, v, err)
		}
	}
}

func standardCloseMessage() []byte {
	data := []byte{0x76, 0x02, 1, 0x62, 0, 0x62, 0, 0x72, 0x63, 2, 1, 0x71, 1}
	crc := Crc16Calculate(data, len(data))
	return append(data, 0x63, byte(crc>>8), byte(crc), 0)
}

func TestMessageCRCIsMandatoryAndUnmodified(t *testing.T) {
	valid := standardCloseMessage()
	if _, err := MessageParse(&Buffer{Bytes: valid}); err != nil {
		t.Fatal(err)
	}
	for _, validate := range []bool{true, false} {
		broken := append([]byte(nil), valid...)
		broken[len(broken)-2] ^= 0x20
		if _, err := MessageParse(&Buffer{Bytes: broken}, validate); err == nil {
			t.Fatal("CRC disabled")
		}
	}
	if _, err := MessageParse(&Buffer{Bytes: valid[:len(valid)-1]}); err == nil {
		t.Fatal("missing terminator accepted")
	}
	for n := 0; n < len(valid); n++ {
		if _, err := MessageParse(&Buffer{Bytes: valid[:n]}); err == nil {
			t.Fatalf("truncation accepted at %d", n)
		}
	}
	if _, err := U8Parse(&Buffer{Bytes: []byte{0x42, 1}}); err == nil {
		t.Fatal("Boolean treated as unsigned")
	}
	if _, err := MessageBodyParse(&Buffer{Bytes: []byte{0x72, 0x43, 2, 1, 0x71, 1}}); err == nil {
		t.Fatal("Boolean tag repaired")
	}
}

func TestTransportCRCAndEscaping(t *testing.T) {
	payload := []byte{0x10, 0x1b, 0x1b, 0x1b, 0x1b, 0x1a, 0x11}
	frame := buildTransportFrame(payload)
	raw, err := TransportRead(bufio.NewReader(bytes.NewReader(frame)))
	if err != nil {
		t.Fatal(err)
	}
	got, err := TransportPayload(raw)
	if err != nil || !bytes.Equal(got, payload) {
		t.Fatalf("escape decoding: % x %v", got, err)
	}
	frame[8] ^= 0x20
	if _, err := TransportPayload(frame); err == nil {
		t.Fatal("damaged transport accepted")
	}
	frame = buildTransportFrame([]byte{1})
	frame[len(frame)-3] = 4
	crc := Crc16Calculate(frame, len(frame)-2)
	frame[len(frame)-2], frame[len(frame)-1] = byte(crc>>8), byte(crc)
	if _, err := TransportPayload(frame); err == nil {
		t.Fatal("invalid fill accepted")
	}
}

func TestRequiredSequenceFields(t *testing.T) {
	if _, err := GetProfilePackRequestParse(&Buffer{Bytes: []byte{0x79, 1, 1, 1, 1, 1, 1, 1, 1, 1}}); err == nil {
		t.Fatal("missing tree path accepted")
	}
	if _, err := ListEntryParse(&Buffer{Bytes: []byte{0x77, 1, 1, 1, 1, 1, 0x40, 1}}); err == nil {
		t.Fatal("malformed list value repaired")
	}
	if _, err := FileParse(nil); err == nil {
		t.Fatal("empty file accepted")
	}
}

func standardMessage(tag uint32, group byte) []byte {
	data := []byte{0x76, 0x02, 1, 0x62, group, 0x62, 0}
	data = append(data, buildMessageBody(tag, minimalPayloadForTag(tag))...)
	crc := crc16Reference(data)
	return append(data, 0x63, byte(crc>>8), byte(crc), 0)
}

func TestStandardFileEnvelopeAndAtomicFailure(t *testing.T) {
	open := standardMessage(MESSAGEOPENRESPONSE, 0)
	close := standardMessage(MESSAGECLOSERESPONSE, 0)
	valid := append(append([]byte{}, open...), close...)
	if got, err := FileParse(valid); err != nil || len(got) != 2 {
		t.Fatalf("valid file: %v", err)
	}
	duplicate := append(append(append([]byte{}, open...), open...), close...)
	if _, err := FileParse(duplicate); err == nil {
		t.Fatal("duplicate Open accepted")
	}
	broken := append([]byte{}, valid...)
	broken[len(broken)-2] ^= 1
	if got, err := FileParse(broken); err == nil || got != nil {
		t.Fatal("partial messages returned on file failure")
	}
	if _, err := FileParse(append(valid, 0)); err == nil {
		t.Fatal("trailing padding accepted as file content")
	}
	groups := append(append(append([]byte{}, open...), standardMessage(MESSAGEGETLISTRESPONSE, 1)...), close...)
	if _, err := FileParse(groups); err == nil {
		t.Fatal("interleaved group accepted")
	}
}

func TestStandardProcParListEntry(t *testing.T) {
	data := []byte{0x72, 0x62, 5, 0x77, 0x07, 1, 0, 1, 8, 0, 255, 1, 1, 1, 1, 0x62, 42, 1}
	got, err := ProcParValueParse(&Buffer{Bytes: data})
	if err != nil || got.ListEntry == nil || got.ListEntry.Value.DataUnsigned != 42 {
		t.Fatalf("list-entry CHOICE: %+v %v", got, err)
	}
}

func FuzzStrictMessageParser(f *testing.F) {
	f.Add(standardCloseMessage())
	f.Add([]byte{0x76, 0x80})
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 65536 {
			t.Skip()
		}
		buf := &Buffer{Bytes: data}
		_, _ = MessageParse(buf)
		if buf.Cursor < 0 || buf.Cursor > len(data) {
			t.Fatalf("cursor out of bounds: %d/%d", buf.Cursor, len(data))
		}
	})
}
