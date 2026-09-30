package sml

import "fmt"

const (
	MESSAGEEND = 0x00

	TYPEFIELD   = 0x70
	LENGTHFIELD = 0x0F
	ANOTHERTL   = 0x80

	TYPEOCTETSTRING = 0x00
	TYPEBOOLEAN     = 0x40
	TYPEINTEGER     = 0x50
	TYPEUNSIGNED    = 0x60
	TYPELIST        = 0x70

	OPTIONALSKIPPED = 0x01
)

type Buffer struct {
	Bytes  []byte
	Cursor int

	err error
}

func BufGetCurrentByte(buf *Buffer) byte {
	if buf.Cursor < 0 || buf.Cursor >= len(buf.Bytes) {
		return 0
	}
	return buf.Bytes[buf.Cursor]
}

func BufUpdateBytesRead(buf *Buffer, delta int) {
	buf.Cursor += delta
}

func Expect(buf *Buffer, expectedType uint8, expectedLength int) error {
	if err := ExpectType(buf, expectedType); err != nil {
		return err
	}

	if length := BufGetNextLength(buf); length != expectedLength {
		return fmt.Errorf("Invalid length: %d (expected %d)", length, expectedLength)
	}

	return nil
}

func ExpectType(buf *Buffer, expectedType uint8) error {
	if buf.err != nil {
		return buf.err
	}
	if buf.Cursor < 0 || buf.Cursor >= len(buf.Bytes) {
		return fmt.Errorf("Unexpected end of buffer at offset %d", buf.Cursor)
	}
	if typefield := BufGetNextType(buf); typefield != expectedType {
		return fmt.Errorf("Unexpected type %02x (expected %02x)", typefield, expectedType)
	}

	return nil
}

func BufGetNextType(buf *Buffer) uint8 {
	return BufGetCurrentByte(buf) & TYPEFIELD
}

func BufGetNextLength(buf *Buffer) int {
	start := buf.Cursor
	if buf.err != nil || start < 0 || start >= len(buf.Bytes) {
		buf.err = fmt.Errorf("Unexpected end of TL field at offset %d", start)
		return -1
	}
	first := buf.Bytes[start]
	typ := first & TYPEFIELD
	if typ != TYPELIST && typ != TYPEOCTETSTRING && typ != TYPEBOOLEAN && typ != TYPEINTEGER && typ != TYPEUNSIGNED {
		buf.err = fmt.Errorf("Reserved TL type at offset %d", start)
		return -1
	}
	if typ == TYPEBOOLEAN && first&ANOTHERTL != 0 {
		buf.err = fmt.Errorf("Extended Boolean TL field at offset %d", start)
		return -1
	}
	length := 0
	for {
		if buf.Cursor >= len(buf.Bytes) {
			buf.err = fmt.Errorf("Unexpected end of TL field at offset %d", start)
			return -1
		}
		b := buf.Bytes[buf.Cursor]
		if buf.Cursor != start && b&TYPEFIELD != 0 {
			buf.err = fmt.Errorf("Reserved continuation TL bits at offset %d", buf.Cursor)
			return -1
		}
		if length > (int(^uint(0)>>1)-int(b&LENGTHFIELD))/16 {
			buf.err = fmt.Errorf("TL length overflow at offset %d", start)
			return -1
		}
		length = length*16 + int(b&LENGTHFIELD)
		buf.Cursor++
		if b&ANOTHERTL == 0 {
			break
		}
	}
	if typ != TYPELIST {
		length -= buf.Cursor - start
	}
	if length < 0 {
		buf.err = fmt.Errorf("Invalid TL length at offset %d", start)
		return -1
	}
	return length
}

func BufOptionalIsSkipped(buf *Buffer) bool {
	if BufGetCurrentByte(buf) == OPTIONALSKIPPED {
		BufUpdateBytesRead(buf, 1)
		return true
	}

	return false
}
