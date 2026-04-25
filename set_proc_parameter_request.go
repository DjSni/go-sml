package sml

type SetProcParameterRequest struct {
	ServerID          OctetString
	Username          OctetString
	Password          OctetString
	ParameterTreePath TreePath
	ParameterTree     *Tree
}

func SetProcParameterRequestParse(buf *Buffer) (SetProcParameterRequest, error) {
	msg := SetProcParameterRequest{}
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

	if msg.ParameterTree, err = TreeParse(buf); err != nil {
		return msg, err
	}

	return msg, nil
}
