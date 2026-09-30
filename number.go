package sml

import "fmt"

const (
	TYPENUMBER_8  = 1
	TYPENUMBER_16 = 2
	TYPENUMBER_32 = 4
	TYPENUMBER_64 = 8
)

func U8Parse(buf *Buffer) (uint8, error) {
	num, err := NumberParse(buf, TYPEUNSIGNED, TYPENUMBER_8)
	return uint8(num), err
}

func U16Parse(buf *Buffer) (uint16, error) {
	num, err := NumberParse(buf, TYPEUNSIGNED, TYPENUMBER_16)
	return uint16(num), err
}

func U32Parse(buf *Buffer) (uint32, error) {
	num, err := NumberParse(buf, TYPEUNSIGNED, TYPENUMBER_32)
	return uint32(num), err
}

func U64Parse(buf *Buffer) (uint64, error) {
	num, err := NumberParse(buf, TYPEUNSIGNED, TYPENUMBER_64)
	return uint64(num), err
}

func I8Parse(buf *Buffer) (int8, error) {
	num, err := NumberParse(buf, TYPEINTEGER, TYPENUMBER_8)
	return int8(num), err
}

func I16Parse(buf *Buffer) (int16, error) {
	num, err := NumberParse(buf, TYPEINTEGER, TYPENUMBER_16)
	return int16(num), err
}

func I32Parse(buf *Buffer) (int32, error) {
	num, err := NumberParse(buf, TYPEINTEGER, TYPENUMBER_32)
	return int32(num), err
}

func I64Parse(buf *Buffer) (int64, error) {
	num, err := NumberParse(buf, TYPEINTEGER, TYPENUMBER_64)
	return int64(num), err
}

func NumberParse(buf *Buffer, numtype uint8, maxSize int) (int64, error) {
	if skip := BufOptionalIsSkipped(buf); skip {
		return 0, nil
	}

	Debug(buf, "NumberParse")

	typefield := BufGetNextType(buf)
	if typefield != numtype {
		return 0, fmt.Errorf("Unexpected type %02x (expected %02x)", typefield, numtype)
	}

	length := BufGetNextLength(buf)
	if length < 1 || length > maxSize {
		return 0, fmt.Errorf("Invalid length: %d", length)
	}

	if length > len(buf.Bytes)-buf.Cursor {
		return 0, fmt.Errorf("Unexpected end of buffer while parsing number")
	}

	if maxSize != 1 && maxSize != 2 && maxSize != 4 && maxSize != 8 {
		return 0, fmt.Errorf("Invalid number type size %02x", maxSize)
	}
	var bits uint64
	if typefield == TYPEINTEGER && buf.Bytes[buf.Cursor]&0x80 != 0 {
		bits = ^uint64(0)
	}
	for _, b := range buf.Bytes[buf.Cursor : buf.Cursor+length] {
		bits = bits<<8 | uint64(b)
	}
	num := int64(bits)

	BufUpdateBytesRead(buf, length)
	// fmt.Printf("num: %d\n", num)

	return num, nil
}
