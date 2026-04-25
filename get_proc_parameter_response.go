package sml

type GetProcParameterResponse struct {
	ServerID          OctetString
	ParameterTreePath TreePath
	ParameterTree     *Tree
}

func GetProcParameterResponseParse(buf *Buffer) (GetProcParameterResponse, error) {
	msg := GetProcParameterResponse{}
	var err error

	if err := Expect(buf, TYPELIST, 3); err != nil {
		return msg, err
	}

	if msg.ServerID, err = OctetStringParse(buf); err != nil {
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
