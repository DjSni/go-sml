package sml

type GetProfileListResponse struct {
	ServerID          OctetString
	ActTime           Time
	RegPeriod         uint32
	ParameterTreePath TreePath
	ValTime           Time
	Status            uint64
	PeriodList        []*PeriodEntry
	Rawdata           OctetString
	PeriodSignature   OctetString
}

func GetProfileListResponseParse(buf *Buffer) (GetProfileListResponse, error) {
	msg := GetProfileListResponse{}
	var err error

	if err := Expect(buf, TYPELIST, 9); err != nil {
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

	if msg.ValTime, err = TimeParse(buf); err != nil {
		return msg, err
	}

	if msg.Status, err = U64Parse(buf); err != nil {
		return msg, err
	}

	if msg.PeriodList, err = parsePeriodEntrySequence(buf); err != nil {
		return msg, err
	}

	if msg.Rawdata, err = OctetStringParse(buf); err != nil {
		return msg, err
	}

	if msg.PeriodSignature, err = OctetStringParse(buf); err != nil {
		return msg, err
	}

	return msg, nil
}

func parsePeriodEntrySequence(buf *Buffer) ([]*PeriodEntry, error) {
	if BufOptionalIsSkipped(buf) {
		return nil, nil
	}

	if err := ExpectType(buf, TYPELIST); err != nil {
		return nil, err
	}

	elems := BufGetNextLength(buf)
	out := make([]*PeriodEntry, 0, elems)
	for i := 0; i < elems; i++ {
		entry, err := PeriodEntryParse(buf)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			out = append(out, entry)
		}
	}

	return out, nil
}
