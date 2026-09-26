package sml

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"os"
	"strings"
	"testing"
)

func TestNamedPulseCaptures(t *testing.T) {
	for _, name := range []string{"01", "05", "06", "07", "09", "10"} {
		t.Run(name, func(t *testing.T) {
			hexData, err := os.ReadFile("testdata/node_data_10x_20260927_010951_" + name + ".hex")
			if err != nil {
				t.Fatal(err)
			}
			data, err := hex.DecodeString(strings.TrimSpace(string(hexData)))
			if err != nil {
				t.Fatal(err)
			}
			frame, err := TransportRead(bufio.NewReader(bytes.NewReader(data)))
			if err != nil {
				t.Fatal(err)
			}
			messages, err := FileParse(frame[8 : len(frame)-8])
			if err != nil {
				t.Fatalf("%d messages: %v", len(messages), err)
			}
			if len(messages) != 3 {
				t.Fatalf("parsed %d messages", len(messages))
			}
		})
	}
}
