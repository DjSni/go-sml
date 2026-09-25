package sml

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"os"
	"strings"
	"testing"
)

func pulseFrame(t *testing.T) []byte {
	t.Helper()
	hexData, err := os.ReadFile("testdata/node_data.hex")
	if err != nil {
		t.Fatal(err)
	}
	data, err := hex.DecodeString(strings.TrimSpace(string(hexData)))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

const nodeDataHex = "1B1B1B1B0101010176050DE7EB4C62006200726301017601010504A2A3C40B0A0149534B000516465C7262016504A2A20862016335990076050DE7EB4D620062007263070177010B0A0149534B000516465C070100620AFFFF7262016504A2A20875770701004032010101010101010449534B0177070100600100FF010101010B0A0149534B000516465C0177070100010800FF650008010401621E52FF6503FEC26A0177070100020800FF0101621E52FF62000177070100100700FF0101621B520053017F01010163BB0E0076050DE7EB4E6200620072630201710163F1A9001B1B1B1B1A00F61F"

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
