package sml

type AttentionResponse struct {
	ServerID         OctetString
	AttentionNumber  OctetString
	AttentionMessage OctetString
	AttentionDetails *Tree
}

func AttentionResponseParse(buf *Buffer) (AttentionResponse, error) {
	msg := AttentionResponse{}
	var err error

	if err := Expect(buf, TYPELIST, 4); err != nil {
		return msg, err
	}

	if msg.ServerID, err = OctetStringParse(buf); err != nil {
		return msg, err
	}

	if msg.AttentionNumber, err = OctetStringParse(buf); err != nil {
		return msg, err
	}

	if msg.AttentionMessage, err = OctetStringParse(buf); err != nil {
		return msg, err
	}

	if msg.AttentionDetails, err = TreeParse(buf); err != nil {
		return msg, err
	}

	return msg, nil
}
