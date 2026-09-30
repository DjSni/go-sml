package sml

import "fmt"

func FileParse(bytes []byte) ([]Message, error) {
	buf := &Buffer{}
	buf.Bytes = append([]byte(nil), bytes...)

	messages := make([]Message, 0)

	for buf.Cursor < len(buf.Bytes) {

		msg, err := MessageParse(buf, true)
		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}
	if len(messages) == 0 {
		return nil, fmt.Errorf("Empty SML file")
	}
	if err := validateFileEnvelope(messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func validateFileEnvelope(messages []Message) error {
	var open uint32
	groups := map[uint8]bool{}
	var previous uint8
	for i, msg := range messages {
		tag := msg.MessageBody.Tag
		if open == 0 {
			if tag != MESSAGEOPENREQUEST && tag != MESSAGEOPENRESPONSE && tag != MESSAGEATTENTIONRESPONSE {
				return fmt.Errorf("SML file must start with Open or Attention (message %d)", i+1)
			}
			open = tag
			groups = map[uint8]bool{}
			previous = msg.GroupID
		} else if tag == MESSAGEOPENREQUEST || tag == MESSAGEOPENRESPONSE {
			return fmt.Errorf("Duplicate Open before Close (message %d)", i+1)
		}
		if msg.GroupID != previous {
			groups[previous] = true
			if groups[msg.GroupID] {
				return fmt.Errorf("Noncontiguous SML message group %d", msg.GroupID)
			}
			previous = msg.GroupID
		}
		if tag == MESSAGECLOSEREQUEST || tag == MESSAGECLOSERESPONSE {
			if (open == MESSAGEOPENREQUEST) != (tag == MESSAGECLOSEREQUEST) {
				return fmt.Errorf("Request/response Close mismatch")
			}
			open = 0
		} else if tag == MESSAGEATTENTIONRESPONSE && (i == len(messages)-1 || messages[i+1].MessageBody.Tag == MESSAGEOPENREQUEST || messages[i+1].MessageBody.Tag == MESSAGEOPENRESPONSE) {
			if open == MESSAGEOPENREQUEST {
				return fmt.Errorf("Request file cannot end with Attention")
			}
			open = 0
		}
	}
	if open != 0 {
		return fmt.Errorf("SML file must end with Close or Attention")
	}
	return nil
}
