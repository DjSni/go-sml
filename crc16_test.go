package sml

import "testing"

func TestCrc16MatchesReferenceImplementation(t *testing.T) {
	vectors := [][]byte{
		{},
		{0x00},
		{0x12, 0x34, 0x56, 0x78},
		[]byte("123456789"),
	}

	for _, vector := range vectors {
		got := Crc16Calculate(vector, len(vector))
		want := crc16Reference(vector)
		if got != want {
			t.Fatalf("crc mismatch for % x: got 0x%04x want 0x%04x", vector, got, want)
		}
	}
}

func TestCrc16KnownVector123456789(t *testing.T) {
	got := Crc16Calculate([]byte("123456789"), 9)
	if got != 0x6e90 {
		t.Fatalf("unexpected crc for 123456789: got 0x%04x want 0x6e90", got)
	}
}

func crc16Reference(data []byte) uint16 {
	var fcs uint16 = 0xffff

	for _, b := range data {
		fcs ^= uint16(b)
		for i := 0; i < 8; i++ {
			if fcs&0x0001 != 0 {
				fcs = (fcs >> 1) ^ 0x8408
			} else {
				fcs >>= 1
			}
		}
	}

	fcs ^= 0xffff
	return ((fcs & 0x00ff) << 8) | ((fcs & 0xff00) >> 8)
}
