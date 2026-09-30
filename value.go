package sml

import "fmt"

type Value struct {
	Typ         uint8
	DataBytes   OctetString
	DataBoolean bool
	DataInt     int64
	// DataUnsigned preserves the entire uint64 range. DataInt is signed only.
	DataUnsigned uint64
	DataTime     *Time
}

func ValueParse(buf *Buffer) (Value, error) {
	v := Value{}
	if err := validateRule(buf, valueRule); err != nil {
		return v, err
	}
	if buf.Cursor >= len(buf.Bytes) {
		return v, fmt.Errorf("Unexpected end of buffer while parsing value")
	}
	typ := BufGetNextType(buf)
	v.Typ = typ
	var err error
	switch typ {
	case TYPEOCTETSTRING:
		v.DataBytes, err = OctetStringParse(buf)
	case TYPEBOOLEAN:
		v.DataBoolean, err = BooleanParse(buf)
	case TYPEINTEGER, TYPEUNSIGNED:
		// Inspect the actual TL length, not just its low nibble.
		peek := *buf
		length := BufGetNextLength(&peek)
		if length < 1 || length > 8 {
			return v, fmt.Errorf("Invalid value number length: %d", length)
		}
		size := 1
		for size < length {
			size *= 2
		}
		n, parseErr := NumberParse(buf, typ, size)
		err = parseErr
		v.Typ |= uint8(size)
		if typ == TYPEUNSIGNED {
			v.DataUnsigned = uint64(n)
		} else {
			v.DataInt = n
		}
	case TYPELIST:
		// SML_ListType is a CHOICE whose sole defined alternative is SML_Time.
		if err = Expect(buf, TYPELIST, 2); err != nil {
			return v, err
		}
		var tag uint8
		if tag, err = U8Parse(buf); err != nil {
			return v, err
		}
		if tag != 1 {
			return v, fmt.Errorf("Invalid SML_ListType choice: %02x", tag)
		}
		var t Time
		t, err = TimeParse(buf)
		if !t.Present && err == nil {
			return v, fmt.Errorf("Missing required SML_Time value")
		}
		v.DataTime = &t
	default:
		err = fmt.Errorf("Unexpected value type %02x", typ)
	}
	return v, err
}
