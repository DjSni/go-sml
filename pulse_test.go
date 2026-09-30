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

func TestDamagedPulseFrameIsRejectedWithoutRepair(t *testing.T) {
	frame := pulseFrame(t)
	if len(frame) != 232 {
		t.Fatalf("fixture length = %d, want 232", len(frame))
	}
	if _, err := TransportParse(frame); err == nil {
		t.Fatal("damaged Pulse capture was silently repaired")
	}
	if _, err := FileParse(frame[8 : len(frame)-8]); err == nil {
		t.Fatal("message CRC checked repaired bytes instead of original bytes")
	}
}

func TestLivePulseFrameParsesThreeMessagesWithCRC(t *testing.T) {
	frame := pulseFixture(t, "testdata/node_data_live_20260927.hex")
	transport, err := TransportRead(bufio.NewReader(bytes.NewReader(frame)))
	if err != nil {
		t.Fatal(err)
	}
	messages, err := TransportParse(transport)
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) != 3 {
		t.Fatalf("parsed %d messages, want 3", len(messages))
	}
	if messages[1].MessageBody.Tag != MESSAGEGETLISTRESPONSE {
		t.Fatal("missing GetListResponse")
	}
	response := messages[1].MessageBody.Data.(GetListResponse)
	if len(response.ValList) != 5 || response.ValList[2].Value.DataUnsigned != 0x04008350 {
		t.Fatalf("unexpected consumption values: %+v", response.ValList)
	}
}

func pulseFixture(t *testing.T, path string) []byte {
	t.Helper()
	hexData, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data, err := hex.DecodeString(strings.TrimSpace(string(hexData)))
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 232 {
		t.Fatalf("fixture %s length = %d, want 232", path, len(data))
	}
	return data
}

func TestTimeParseRejectsPulse45Variant(t *testing.T) {
	buf := &Buffer{Bytes: []byte{0x72, 0x62, 0x01, 0x45, 0x04, 0xA2, 0x9B, 0x2A}}
	if _, err := TimeParse(buf); err == nil {
		t.Fatal("nonstandard 0x45 time accepted")
	}
}

func TestTimeParseUnsigned65Variant(t *testing.T) {
	buf := &Buffer{Bytes: []byte{0x72, 0x62, 0x01, 0x65, 0x04, 0xA2, 0x9B, 0x2A}}
	got, err := TimeParse(buf)
	if err != nil || got.Timestamp != 0x04A29B2A || got.Tag != 1 {
		t.Fatalf("TimeParse(0x65) = %+v, %v", got, err)
	}
}

func TestValueParseRejectsLengthlessBoolean(t *testing.T) {
	if _, err := ValueParse(&Buffer{Bytes: []byte{0x40}}); err == nil {
		t.Fatal("lengthless Boolean accepted")
	}
}

func TestBooleanValueRemainsBoolean(t *testing.T) {
	value, err := ValueParse(&Buffer{Bytes: []byte{0x42, 0x01}})
	if err != nil {
		t.Fatal(err)
	}
	if value.Typ != TYPEBOOLEAN || !value.DataBoolean {
		t.Fatalf("Boolean value = type %#x value %t", value.Typ, value.DataBoolean)
	}
}

func TestPulseFrameCRCErrorIsRejected(t *testing.T) {
	frame := pulseFixture(t, "testdata/node_data_live_20260927.hex")
	payload, err := TransportPayload(frame)
	if err != nil {
		t.Fatal(err)
	}
	payload[20] ^= 0x01
	_, err = FileParse(payload)
	if err == nil {
		t.Fatal("CRC-corrupted Pulse frame was accepted")
	}
}
