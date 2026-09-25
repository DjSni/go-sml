package sml

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

func pulseFrame(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(`S:\node_data.bin`)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestPulseFrameParsesThreeMessagesWithCRC(t *testing.T) {
	frame := pulseFrame(t)
	if len(frame) != 232 {
		t.Fatalf("fixture length = %d, want 232", len(frame))
	}
	transport, err := TransportRead(bufio.NewReader(bytes.NewReader(frame)))
	if err != nil {
		t.Fatal(err)
	}
	payload := transport[8 : len(transport)-8]
	buf := &Buffer{Bytes: payload}
	messages := make([]Message, 0, 3)
	for buf.Cursor < len(buf.Bytes) {
		message, parseErr := MessageParse(buf, true)
		if parseErr != nil {
			t.Fatalf("message %d at offset %d: %v", len(messages)+1, buf.Cursor, parseErr)
		}
		messages = append(messages, message)
	}
	if len(messages) != 3 {
		t.Fatalf("parsed %d messages, want 3", len(messages))
	}
	if messages[1].MessageBody.Tag != MESSAGEGETLISTRESPONSE {
		t.Fatalf("message 2 tag = %#x, want MESSAGEGETLISTRESPONSE", messages[1].MessageBody.Tag)
	}
}

func TestTimeParsePulse45Variant(t *testing.T) {
	buf := &Buffer{Bytes: []byte{0x72, 0x62, 0x01, 0x45, 0x04, 0xA2, 0x9B, 0x2A}}
	got, err := TimeParse(buf)
	if err != nil || got != 0x04A29B2A {
		t.Fatalf("TimeParse(0x45) = %#x, %v", got, err)
	}
}

func TestTimeParseUnsigned65Variant(t *testing.T) {
	buf := &Buffer{Bytes: []byte{0x72, 0x62, 0x01, 0x65, 0x04, 0xA2, 0x9B, 0x2A}}
	got, err := TimeParse(buf)
	if err != nil || got != 0x04A29B2A {
		t.Fatalf("TimeParse(0x65) = %#x, %v", got, err)
	}
}

func TestPulseFrameCRCErrorIsRejected(t *testing.T) {
	frame := pulseFrame(t)
	frame[20] ^= 0x01
	_, err := FileParse(frame[8 : len(frame)-8])
	if err == nil {
		t.Fatal("CRC-corrupted Pulse frame was accepted")
	}
}
