package sml

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"os"
	"strings"
	"testing"
)

func TestNamedDamagedPulseCapturesAreRejected(t *testing.T) {
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
			if _, err := TransportRead(bufio.NewReader(bytes.NewReader(data))); err == nil {
				t.Fatal("damaged transport frame accepted")
			}
			if _, err := FileParse(data[8 : len(data)-8]); err == nil {
				t.Fatal("damaged SML messages silently repaired")
			}
		})
	}
}
