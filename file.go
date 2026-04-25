package sml

func FileParse(bytes []byte) ([]Message, error) {
	buf := &Buffer{}
	buf.Bytes = append([]byte(nil), bytes...)

	messages := make([]Message, 0)

	for buf.Cursor < len(buf.Bytes) {
		if BufGetCurrentByte(buf) == MESSAGEEND {
			// reading trailing zeroed bytes
			BufUpdateBytesRead(buf, 1)
			continue
		}

		if BufGetCurrentByte(buf)&TYPEFIELD != TYPELIST {
			break
		}

		msg, err := MessageParse(buf, true)
		if err != nil {
			return messages, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}
