package sml

type GetProfilePackResponse struct {
	ServerID          OctetString
	ActTime           Time
	RegPeriod         uint32
	ParameterTreePath TreePath
	HeaderList        []*ProfObjHeaderEntry
	PeriodList        []*ProfObjPeriodEntry
	Rawdata           OctetString
	ProfileSignature  OctetString
}

type ProfObjHeaderEntry struct {
	ObjName OctetString
	Unit    uint8
	Scaler  int8
}

type ProfObjPeriodEntry struct {
	ValTime         Time
	Status          uint64
	ValueList       []*ValueEntry
	PeriodSignature OctetString
}

type ValueEntry struct {
	Value          Value
	ValueSignature OctetString
}

func GetProfilePackResponseParse(buf *Buffer) (GetProfilePackResponse, error) {
	msg := GetProfilePackResponse{}
	var err error

	if err := Expect(buf, TYPELIST, 8); err != nil {
		return msg, err
	}

	if msg.ServerID, err = OctetStringParse(buf); err != nil {
		return msg, err
	}

	if msg.ActTime, err = TimeParse(buf); err != nil {
		return msg, err
	}

	if msg.RegPeriod, err = U32Parse(buf); err != nil {
		return msg, err
	}

	if msg.ParameterTreePath, err = TreePathParse(buf); err != nil {
		return msg, err
	}

	if msg.HeaderList, err = parseProfObjHeaderEntrySequence(buf); err != nil {
		return msg, err
	}

	if msg.PeriodList, err = parseProfObjPeriodEntrySequence(buf); err != nil {
		return msg, err
	}

	if msg.Rawdata, err = OctetStringParse(buf); err != nil {
		return msg, err
	}

	if msg.ProfileSignature, err = OctetStringParse(buf); err != nil {
		return msg, err
	}

	return msg, nil
}

func ProfObjHeaderEntryParse(buf *Buffer) (*ProfObjHeaderEntry, error) {
	entry := &ProfObjHeaderEntry{}
	var err error

	if err := Expect(buf, TYPELIST, 3); err != nil {
		return nil, err
	}

	if entry.ObjName, err = OctetStringParse(buf); err != nil {
		return nil, err
	}
	if entry.Unit, err = U8Parse(buf); err != nil {
		return nil, err
	}
	if entry.Scaler, err = I8Parse(buf); err != nil {
		return nil, err
	}

	return entry, nil
}

func ProfObjPeriodEntryParse(buf *Buffer) (*ProfObjPeriodEntry, error) {
	entry := &ProfObjPeriodEntry{}
	var err error

	if err := Expect(buf, TYPELIST, 4); err != nil {
		return nil, err
	}

	if entry.ValTime, err = TimeParse(buf); err != nil {
		return nil, err
	}
	if entry.Status, err = U64Parse(buf); err != nil {
		return nil, err
	}
	if entry.ValueList, err = parseValueEntrySequence(buf); err != nil {
		return nil, err
	}
	if entry.PeriodSignature, err = OctetStringParse(buf); err != nil {
		return nil, err
	}

	return entry, nil
}

func ValueEntryParse(buf *Buffer) (*ValueEntry, error) {
	entry := &ValueEntry{}
	var err error

	if err := Expect(buf, TYPELIST, 2); err != nil {
		return nil, err
	}

	if entry.Value, err = ValueParse(buf); err != nil {
		return nil, err
	}

	if entry.ValueSignature, err = OctetStringParse(buf); err != nil {
		return nil, err
	}

	return entry, nil
}

func parseProfObjHeaderEntrySequence(buf *Buffer) ([]*ProfObjHeaderEntry, error) {
	if BufOptionalIsSkipped(buf) {
		return nil, nil
	}

	if err := ExpectType(buf, TYPELIST); err != nil {
		return nil, err
	}

	elems := BufGetNextLength(buf)
	out := make([]*ProfObjHeaderEntry, 0, elems)
	for i := 0; i < elems; i++ {
		entry, err := ProfObjHeaderEntryParse(buf)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			out = append(out, entry)
		}
	}

	return out, nil
}

func parseProfObjPeriodEntrySequence(buf *Buffer) ([]*ProfObjPeriodEntry, error) {
	if BufOptionalIsSkipped(buf) {
		return nil, nil
	}

	if err := ExpectType(buf, TYPELIST); err != nil {
		return nil, err
	}

	elems := BufGetNextLength(buf)
	out := make([]*ProfObjPeriodEntry, 0, elems)
	for i := 0; i < elems; i++ {
		entry, err := ProfObjPeriodEntryParse(buf)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			out = append(out, entry)
		}
	}

	return out, nil
}

func parseValueEntrySequence(buf *Buffer) ([]*ValueEntry, error) {
	if BufOptionalIsSkipped(buf) {
		return nil, nil
	}

	if err := ExpectType(buf, TYPELIST); err != nil {
		return nil, err
	}

	elems := BufGetNextLength(buf)
	out := make([]*ValueEntry, 0, elems)
	for i := 0; i < elems; i++ {
		entry, err := ValueEntryParse(buf)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			out = append(out, entry)
		}
	}

	return out, nil
}
