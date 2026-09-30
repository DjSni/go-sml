package sml

import "fmt"

// These validators implement the binary ASN.1 shapes in BSI IVb 5.1.
// They distinguish a missing OPTIONAL field from a missing required field;
// low-level primitive decoders alone cannot make that distinction.
type wireValue struct {
	typ      byte
	data     []byte
	children []wireValue
}

func readWire(buf *Buffer, depth int) (wireValue, error) {
	n := wireValue{}
	if depth > 128 {
		return n, fmt.Errorf("SML nesting exceeds 128 levels")
	}
	if buf.Cursor < 0 || buf.Cursor >= len(buf.Bytes) {
		return n, fmt.Errorf("Unexpected end of SML value")
	}
	n.typ = BufGetNextType(buf)
	length := BufGetNextLength(buf)
	if buf.err != nil {
		return n, buf.err
	}
	if length < 0 {
		return n, fmt.Errorf("Invalid TL length")
	}
	if n.typ == TYPELIST {
		if length > len(buf.Bytes)-buf.Cursor {
			return n, fmt.Errorf("Truncated SML list")
		}
		for i := 0; i < length; i++ {
			child, err := readWire(buf, depth+1)
			if err != nil {
				return n, err
			}
			n.children = append(n.children, child)
		}
	} else {
		if length > len(buf.Bytes)-buf.Cursor {
			return n, fmt.Errorf("Truncated SML value")
		}
		n.data = buf.Bytes[buf.Cursor : buf.Cursor+length]
		buf.Cursor += length
	}
	return n, nil
}

type rule func(wireValue) error

func optional(r rule) rule {
	return func(n wireValue) error {
		if n.typ == TYPEOCTETSTRING && len(n.data) == 0 {
			return nil
		}
		return r(n)
	}
}

func primitive(typ byte, max int) rule {
	return func(n wireValue) error {
		if n.typ != typ {
			return fmt.Errorf("Unexpected type %02x (expected %02x)", n.typ, typ)
		}
		if typ == TYPEOCTETSTRING {
			return nil
		}
		if len(n.data) < 1 || len(n.data) > max {
			return fmt.Errorf("Invalid primitive width %d (maximum %d)", len(n.data), max)
		}
		return nil
	}
}

func sequence(fields ...rule) rule {
	return func(n wireValue) error {
		if n.typ != TYPELIST || len(n.children) != len(fields) {
			return fmt.Errorf("Expected sequence of %d fields", len(fields))
		}
		for i, r := range fields {
			if err := r(n.children[i]); err != nil {
				return fmt.Errorf("field %d: %w", i+1, err)
			}
		}
		return nil
	}
}

func sequenceOf(r rule, min int) rule {
	return func(n wireValue) error {
		if n.typ != TYPELIST || len(n.children) < min {
			return fmt.Errorf("Expected list with at least %d elements", min)
		}
		for i, c := range n.children {
			if err := r(c); err != nil {
				return fmt.Errorf("element %d: %w", i+1, err)
			}
		}
		return nil
	}
}

var (
	octets  = primitive(TYPEOCTETSTRING, 0)
	u8      = primitive(TYPEUNSIGNED, 1)
	u16     = primitive(TYPEUNSIGNED, 2)
	u32     = primitive(TYPEUNSIGNED, 4)
	u64     = primitive(TYPEUNSIGNED, 8)
	i8      = primitive(TYPEINTEGER, 1)
	i16     = primitive(TYPEINTEGER, 2)
	i64     = primitive(TYPEINTEGER, 8)
	boolean = primitive(TYPEBOOLEAN, 1)
)

func wireUnsigned(n wireValue) uint64 {
	var v uint64
	for _, b := range n.data {
		v = v<<8 | uint64(b)
	}
	return v
}

func timeRule(n wireValue) error {
	if n.typ != TYPELIST || len(n.children) != 2 {
		return fmt.Errorf("Expected SML_Time CHOICE")
	}
	if err := u8(n.children[0]); err != nil {
		return err
	}
	switch wireUnsigned(n.children[0]) {
	case 1, 2:
		return u32(n.children[1])
	case 3:
		return sequence(u32, i16, i16)(n.children[1])
	default:
		return fmt.Errorf("Unknown SML_Time choice")
	}
}

func valueRule(n wireValue) error {
	switch n.typ {
	case TYPEOCTETSTRING:
		return octets(n)
	case TYPEBOOLEAN:
		return boolean(n)
	case TYPEINTEGER:
		return i64(n)
	case TYPEUNSIGNED:
		return u64(n)
	case TYPELIST:
		if err := sequence(u8, timeRule)(n); err != nil {
			return err
		}
		if wireUnsigned(n.children[0]) != 1 {
			return fmt.Errorf("Unknown SML_ListType choice")
		}
		return nil
	default:
		return fmt.Errorf("Reserved SML_Value type")
	}
}

func listEntryRule(n wireValue) error {
	return sequence(octets, optional(u64), optional(timeRule), optional(u8), optional(i8), valueRule, optional(octets))(n)
}

func periodEntryRule(n wireValue) error {
	return sequence(octets, u8, i8, valueRule, optional(octets))(n)
}

func tupelRule(n wireValue) error {
	return sequence(octets, timeRule, u64, u8, i8, i64, u8, i8, i64, u8, i8, i64, octets, u8, i8, i64, u8, i8, i64, u8, i8, i64, octets)(n)
}

func treeRule(n wireValue) error {
	return sequence(octets, optional(procParRule), optional(sequenceOf(treeRule, 1)))(n)
}

func procParRule(n wireValue) error {
	if n.typ != TYPELIST || len(n.children) != 2 {
		return fmt.Errorf("Expected SML_ProcParValue CHOICE")
	}
	if err := u8(n.children[0]); err != nil {
		return err
	}
	switch wireUnsigned(n.children[0]) {
	case 1:
		return valueRule(n.children[1])
	case 2:
		return periodEntryRule(n.children[1])
	case 3:
		return tupelRule(n.children[1])
	case 4:
		return timeRule(n.children[1])
	case 5:
		return listEntryRule(n.children[1])
	default:
		return fmt.Errorf("Unknown SML_ProcParValue choice")
	}
}

func bodyRule(tag uint32) rule {
	path := sequenceOf(octets, 1)
	profileRequest := sequence(optional(octets), optional(octets), optional(octets), optional(boolean), optional(timeRule), optional(timeRule), path, optional(sequenceOf(octets, 0)), optional(treeRule))
	switch tag {
	case MESSAGEOPENREQUEST:
		return sequence(optional(octets), octets, octets, optional(octets), optional(octets), optional(octets), optional(u8))
	case MESSAGEOPENRESPONSE:
		return sequence(optional(octets), optional(octets), octets, octets, optional(timeRule), optional(u8))
	case MESSAGECLOSEREQUEST, MESSAGECLOSERESPONSE:
		return sequence(optional(octets))
	case MESSAGEGETLISTREQUEST:
		return sequence(octets, optional(octets), optional(octets), optional(octets), optional(octets))
	case MESSAGEGETLISTRESPONSE:
		return sequence(optional(octets), octets, optional(octets), optional(timeRule), sequenceOf(listEntryRule, 0), optional(octets), optional(timeRule))
	case MESSAGEGETPROFILEPACKREQUEST, MESSAGEGETPROFILELISTREQUEST:
		return profileRequest
	case MESSAGEGETPROFILEPACKRESPONSE:
		return func(n wireValue) error {
			header := sequenceOf(sequence(octets, u8, i8), 0)
			period := sequence(timeRule, u64, sequenceOf(sequence(valueRule, optional(octets)), 0), optional(octets))
			if err := sequence(octets, timeRule, u32, path, header, sequenceOf(period, 0), optional(octets), optional(octets))(n); err != nil {
				return err
			}
			for _, p := range n.children[5].children {
				if len(p.children[2].children) != len(n.children[4].children) {
					return fmt.Errorf("Profile header/value count mismatch")
				}
			}
			return nil
		}
	case MESSAGEGETPROFILELISTRESPONSE:
		return sequence(octets, timeRule, u32, path, timeRule, u64, sequenceOf(periodEntryRule, 0), optional(octets), optional(octets))
	case MESSAGEGETPROCPARAMETERREQUEST:
		return sequence(optional(octets), optional(octets), optional(octets), path, optional(octets))
	case MESSAGEGETPROCPARAMETERRESPONSE:
		return sequence(octets, path, treeRule)
	case MESSAGESETPROCPARAMETERREQUEST:
		return sequence(optional(octets), optional(octets), optional(octets), path, treeRule)
	case MESSAGEATTENTIONRESPONSE:
		return sequence(octets, octets, optional(octets), optional(treeRule))
	default:
		return func(wireValue) error { return fmt.Errorf("Invalid message type: % x", tag) }
	}
}

func validateBody(buf *Buffer, tag uint32) error {
	return validateRule(buf, bodyRule(tag))
}

func validateRule(buf *Buffer, r rule) error {
	peek := *buf
	n, err := readWire(&peek, 0)
	if err != nil {
		return err
	}
	return r(n)
}
