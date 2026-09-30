package sml

import "fmt"

// Time preserves all three SML_Time choices (BSI IVb, definitions F-H).
// Present distinguishes an absent optional time from timestamp zero.
type Time struct {
	Present          bool
	Tag              uint8
	Timestamp        uint32
	LocalOffset      int16
	SeasonTimeOffset int16
}

func TimeParse(buf *Buffer) (Time, error) {
	t := Time{}
	if BufOptionalIsSkipped(buf) {
		return t, nil
	}
	if err := validateRule(buf, timeRule); err != nil {
		return t, err
	}
	if err := Expect(buf, TYPELIST, 2); err != nil {
		return t, err
	}
	var err error
	if t.Tag, err = U8Parse(buf); err != nil {
		return t, err
	}
	switch t.Tag {
	case 1, 2:
		if BufGetNextType(buf) != TYPEUNSIGNED {
			return t, fmt.Errorf("Invalid time value type %02x", BufGetNextType(buf))
		}
		t.Timestamp, err = U32Parse(buf)
	case 3:
		if err = Expect(buf, TYPELIST, 3); err != nil {
			return t, err
		}
		if t.Timestamp, err = U32Parse(buf); err != nil {
			return t, err
		}
		if t.LocalOffset, err = I16Parse(buf); err != nil {
			return t, err
		}
		t.SeasonTimeOffset, err = I16Parse(buf)
	default:
		return t, fmt.Errorf("Invalid time choice %02x", t.Tag)
	}
	if err != nil {
		return t, err
	}
	t.Present = true
	return t, nil
}
