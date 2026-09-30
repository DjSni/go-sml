package sml

import (
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Optional external capture audit; the portable regression fixtures are in
// testdata. No drive letter or machine-local path is embedded in the tests.
func TestStrictPulseLongStream(t *testing.T) {
	root := os.Getenv("PULSE_CAPTURE_DIR")
	if root == "" {
		t.Skip("set PULSE_CAPTURE_DIR to audit the 60-capture stream")
	}
	paths, err := filepath.Glob(filepath.Join(root, "node_data_long60x_20260927_015854_*.hex"))
	if err != nil || len(paths) != 60 {
		t.Fatalf("expected 60 captures, got %d (%v)", len(paths), err)
	}
	accepted, rejected := 0, 0
	for _, path := range paths {
		h, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		frame, err := hex.DecodeString(strings.Join(strings.Fields(string(h)), ""))
		if err != nil {
			t.Fatal(err)
		}
		bin, err := os.ReadFile(strings.TrimSuffix(path, ".hex") + ".bin")
		if err != nil {
			t.Fatal(err)
		}
		if len(frame) != 232 || !bytes.Equal(frame, bin) {
			t.Fatalf("inconsistent capture: %s", path)
		}
		original := append([]byte(nil), frame...)
		// Independent bitwise CRC reference: rejection is an expectation,
		// not permission for the audit to swallow an arbitrary parser failure.
		wantRejected := crc16Reference(frame[:len(frame)-2]) != uint16(frame[len(frame)-2])<<8|uint16(frame[len(frame)-1])
		messages, parseErr := TransportParse(frame)
		if !bytes.Equal(frame, original) {
			t.Fatal("parser modified capture")
		}
		if parseErr == nil {
			if wantRejected {
				t.Fatalf("CRC-damaged frame accepted: %s", path)
			}
			if len(messages) != 3 || messages[1].MessageBody.Tag != MESSAGEGETLISTRESPONSE {
				t.Fatalf("incomplete accepted capture: %s", path)
			}
			accepted++
		} else {
			if !wantRejected {
				t.Fatalf("CRC-valid frame rejected: %s: %v", path, parseErr)
			}
			if messages != nil {
				t.Fatal("error returned partial readings")
			}
			rejected++
			t.Logf("%s: rejected: %v", filepath.Base(path), parseErr)
		}
	}
	if accepted != 27 || rejected != 33 {
		t.Fatalf("changed stream classification: %d/%d", accepted, rejected)
	}
	t.Logf("strict stream: %d accepted, %d rejected (no repairs)", accepted, rejected)
}
