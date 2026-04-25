package sml

type GetProcParameterRequest struct {
	ServerID          OctetString
	Username          OctetString
	Password          OctetString
	ParameterTreePath TreePath
	Attribute         OctetString
}

func GetProcParameterRequestParse(buf *Buffer) (GetProcParameterRequest, error) {
	msg := GetProcParameterRequest{}
	var err error

	if err := Expect(buf, TYPELIST, 5); err != nil {
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

	if msg.ParameterTreePath, err = TreePathParse(buf); err != nil {
		return msg, err
	}

	if msg.Attribute, err = OctetStringParse(buf); err != nil {
		return msg, err
	}

	return msg, nil
}
