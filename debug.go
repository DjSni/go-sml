package sml

import (
	"fmt"
)

var DebugEnable bool

func Debug(buf *Buffer, function string) {
	if DebugEnable {
		end := min(len(buf.Bytes), buf.Cursor+30)
		if buf.Cursor < 0 || buf.Cursor > end {
			return
		}
		fmt.Printf("%-22s % x\n", function, buf.Bytes[buf.Cursor:end])
	}
}
