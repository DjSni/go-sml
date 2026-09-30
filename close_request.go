package sml

type CloseRequest struct {
	GlobalSignature OctetString
}

func CloseRequestParse(buf *Buffer) (CloseRequest, error) {
	msg := CloseRequest{}
	if err := validateBody(buf, MESSAGECLOSEREQUEST); err != nil {
		return msg, err
	}
	var err error

	if err := Expect(buf, TYPELIST, 1); err != nil {
		return msg, err
	}

	if msg.GlobalSignature, err = OctetStringParse(buf); err != nil {
		return msg, err
	}

	return msg, nil
}
