package sml

type ObjReqEntry = OctetString

type GetProfilePackRequest struct {
	ServerID          OctetString
	Username          OctetString
	Password          OctetString
	WithRawdata       bool
	BeginTime         Time
	EndTime           Time
	ParameterTreePath TreePath
	ObjectList        []ObjReqEntry
	DasDetails        *Tree
}

func GetProfilePackRequestParse(buf *Buffer) (GetProfilePackRequest, error) {
	msg := GetProfilePackRequest{}
	var err error

	if err := Expect(buf, TYPELIST, 9); err != nil {
		return msg, err
	}

	if msg.ServerID, err = OctetStringParse(buf); err != nil {
		return msg, err
	}

	if msg.Username, err = OctetStringParse(buf); err != nil {
		return msg, err
	}

	if msg.Password, err = OctetStringParse(buf); err != nil {
		return msg, err
	}

	if msg.WithRawdata, err = BooleanParse(buf); err != nil {
		return msg, err
	}

	if msg.BeginTime, err = TimeParse(buf); err != nil {
		return msg, err
	}

	if msg.EndTime, err = TimeParse(buf); err != nil {
		return msg, err
	}

	if msg.ParameterTreePath, err = TreePathParse(buf); err != nil {
		return msg, err
	}

	if !BufOptionalIsSkipped(buf) {
		if err := ExpectType(buf, TYPELIST); err != nil {
			return msg, err
		}

		elems := BufGetNextLength(buf)
		msg.ObjectList = make([]ObjReqEntry, 0, elems)
		for i := 0; i < elems; i++ {
			obj, err := OctetStringParse(buf)
			if err != nil {
				return msg, err
			}
			if obj != nil {
				msg.ObjectList = append(msg.ObjectList, obj)
			}
		}
	}

	if msg.DasDetails, err = TreeParse(buf); err != nil {
		return msg, err
	}

	return msg, nil
}
