package sml

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The upstream suite checks all .bin/.hex pairs. This runs that same check
// without requiring Python, then audits every complete frame against this fork.
// Invalid real-world dumps are reported, not repaired or declared conformant.
func TestExternalLibSMLFixtures(t *testing.T) {
	root := os.Getenv("LIBSML_TESTING_DIR")
	if root == "" {
		t.Skip("set LIBSML_TESTING_DIR to run the complete external corpus")
	}
	files, err := filepath.Glob(filepath.Join(root, "*.bin"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no external fixtures: %v", err)
	}
	hexFiles, err := filepath.Glob(filepath.Join(root, "*.hex"))
	if err != nil || len(hexFiles) != len(files) {
		t.Fatal("missing bin/hex counterparts")
	}
	if len(files) != len(strictCorpusCounts) {
		t.Fatal("external corpus changed; review new fixtures against the specification")
	}
	totalAccepted, totalRejected := 0, 0
	for _, path := range files {
		t.Run(filepath.Base(path), func(t *testing.T) {
			bin, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			h, err := os.ReadFile(strings.TrimSuffix(path, ".bin") + ".hex")
			if err != nil {
				t.Fatal(err)
			}
			h = []byte(strings.Join(strings.Fields(string(h)), ""))
			decoded, err := hex.DecodeString(string(h))
			if err != nil || !bytes.Equal(bin, decoded) {
				t.Fatal("upstream bin/hex inconsistency")
			}
			r := bufio.NewReader(bytes.NewReader(bin))
			accepted, rejected := 0, 0
			for {
				frame, err := TransportRead(r)
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					rejected++
					t.Logf("transport rejected: %v", err)
					if r.Buffered() == 0 || err.Error() == "premature eof" {
						break
					}
					continue
				}
				if _, err = TransportParse(frame); err != nil {
					rejected++
					t.Logf("message rejected: %v", err)
				} else {
					accepted++
				}
			}
			t.Logf("audit: %d accepted, %d rejected", accepted, rejected)
			want, ok := strictCorpusCounts[filepath.Base(path)]
			if !ok || want != [2]int{accepted, rejected} {
				t.Fatalf("unexpected strict-corpus result: got %v want %v", [2]int{accepted, rejected}, want)
			}
			totalAccepted += accepted
			totalRejected += rejected
		})
	}
	t.Logf("complete corpus: %d pairs, %d accepted frames, %d rejected frames/tails", len(files), totalAccepted, totalRejected)
}

// Corpus revision a3c7869. Counts include incomplete serial-dump tails.
// Holley ZDBA omits the SML_Time CHOICE required by IVb F-H.
// Corrupt CRCs and the named error dump stay rejected.
var strictCorpusCounts = map[string][2]int{
	"DZG_DVS-7412.2_jmberg.bin":                 {1, 0},
	"DZG_DVS-7420.2V.G2_mtr0.bin":               {1, 1},
	"DZG_DVS-7420.2V.G2_mtr1.bin":               {4, 1},
	"DZG_DVS-7420.2V.G2_mtr1_error.bin":         {0, 4},
	"DZG_DVS-7420.2V.G2_mtr2.bin":               {3, 0},
	"DZG_DVS-7420.2V.G2_mtr2_neg.bin":           {3, 1},
	"DrNeuhaus_SMARTY_ix-130.bin":               {12, 1},
	"EMH-ED300L_consumption.bin":                {1, 1},
	"EMH-ED300L_delivery.bin":                   {2, 1},
	"EMH_eHZ-GW8E2A500AK2.bin":                  {16, 1},
	"EMH_eHZ-HW8E2A5L0EK2P.bin":                 {12, 1},
	"EMH_eHZ-HW8E2A5L0EK2P_1.bin":               {12, 1},
	"EMH_eHZ-HW8E2A5L0EK2P_2.bin":               {1, 0},
	"EMH_eHZ-HW8E2AWL0EK2P.bin":                 {13, 1},
	"EMH_eHZ-IW8E2A5L0EK2P_with_error.bin":      {11, 1},
	"EMH_eHZ-IW8E2AWL0EK2P.bin":                 {12, 1},
	"EMH_eHZ361L5R.bin":                         {1, 0},
	"EMH_eHZ361L5R_1.bin":                       {1, 0},
	"EMH_mME40-AE6AKF0K0.bin":                   {12, 1},
	"EasyMeter_Q3A_A1064V1009.bin":              {4, 4},
	"HOLLEY_DTZ541-BDBA_with_PIN.bin":           {2, 0},
	"HOLLEY_DTZ541-BDBA_without_PIN.bin":        {1, 0},
	"HOLLEY_DTZ541-ZDBA.bin":                    {0, 8},
	"ISKRA_MT175_D1A52-V22-K0t.bin":             {8, 1},
	"ISKRA_MT175_eHZ.bin":                       {10, 1},
	"ISKRA_MT631-D1A52-K0z-H01_with_PIN.bin":    {5, 0},
	"ISKRA_MT631-D1A52-K0z-H01_without_PIN.bin": {5, 0},
	"ISKRA_MT631-D2A51-V22-K0z_with_PIN.bin":    {4, 0},
	"ISKRA_MT631-D2A51-V22-K0z_without_PIN.bin": {2, 0},
	"ISKRA_MT691_eHZ-MS2020.bin":                {18, 1},
	"ITRON_OpenWay-3.HZ.bin":                    {1, 0},
	"ITRON_OpenWay-3.HZ_with_PIN.bin":           {5, 0},
	"ITRON_OpenWay-3.HZ_without_PIN.bin":        {2, 0},
	"dzg_dwsb20_2th_2byte.bin":                  {15, 1},
	"dzg_dwsb20_2th_3byte.bin":                  {14, 3},
	"eBZ_DD3_DD32R06DTA-SMZ1.bin":               {4, 0},
	"eBZ_DD3_DD3BZ06DTA-SMZ1_without_PIN.bin":   {2, 0},
}
