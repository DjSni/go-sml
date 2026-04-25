package sml

import "testing"

func TestPortedMessageParsersAcceptOptionalFields(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		fn   func(*Buffer) error
	}{
		{
			name: "AttentionResponse",
			data: []byte{0x74, 0x01, 0x01, 0x01, 0x01},
			fn: func(buf *Buffer) error {
				_, err := AttentionResponseParse(buf)
				return err
			},
		},
		{
			name: "GetProcParameterRequest",
			data: []byte{0x75, 0x01, 0x01, 0x01, 0x01, 0x01},
			fn: func(buf *Buffer) error {
				_, err := GetProcParameterRequestParse(buf)
				return err
			},
		},
		{
			name: "GetProcParameterResponse",
			data: []byte{0x73, 0x01, 0x01, 0x01},
			fn: func(buf *Buffer) error {
				_, err := GetProcParameterResponseParse(buf)
				return err
			},
		},
		{
			name: "SetProcParameterRequest",
			data: []byte{0x75, 0x01, 0x01, 0x01, 0x01, 0x01},
			fn: func(buf *Buffer) error {
				_, err := SetProcParameterRequestParse(buf)
				return err
			},
		},
		{
			name: "GetProfilePackRequest",
			data: []byte{0x79, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			fn: func(buf *Buffer) error {
				_, err := GetProfilePackRequestParse(buf)
				return err
			},
		},
		{
			name: "GetProfileListRequest",
			data: []byte{0x79, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			fn: func(buf *Buffer) error {
				_, err := GetProfileListRequestParse(buf)
				return err
			},
		},
		{
			name: "GetProfileListResponse",
			data: []byte{0x79, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			fn: func(buf *Buffer) error {
				_, err := GetProfileListResponseParse(buf)
				return err
			},
		},
		{
			name: "GetProfilePackResponse",
			data: []byte{0x78, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			fn: func(buf *Buffer) error {
				_, err := GetProfilePackResponseParse(buf)
				return err
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := &Buffer{Bytes: tc.data}
			if err := tc.fn(buf); err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			if buf.Cursor != len(tc.data) {
				t.Fatalf("parser consumed %d bytes, expected %d", buf.Cursor, len(tc.data))
			}
		})
	}
}

func TestMessageBodyParsePortedTypes(t *testing.T) {
	tests := []struct {
		name string
		tag  []byte
		body []byte
	}{
		{name: "GetProfilePackRequest", tag: []byte{0x65, 0x00, 0x00, 0x03, 0x00}, body: []byte{0x79, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}},
		{name: "GetProfilePackResponse", tag: []byte{0x65, 0x00, 0x00, 0x03, 0x01}, body: []byte{0x78, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}},
		{name: "GetProfileListRequest", tag: []byte{0x65, 0x00, 0x00, 0x04, 0x00}, body: []byte{0x79, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}},
		{name: "GetProfileListResponse", tag: []byte{0x65, 0x00, 0x00, 0x04, 0x01}, body: []byte{0x79, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}},
		{name: "GetProcParameterRequest", tag: []byte{0x65, 0x00, 0x00, 0x05, 0x00}, body: []byte{0x75, 0x01, 0x01, 0x01, 0x01, 0x01}},
		{name: "GetProcParameterResponse", tag: []byte{0x65, 0x00, 0x00, 0x05, 0x01}, body: []byte{0x73, 0x01, 0x01, 0x01}},
		{name: "SetProcParameterRequest", tag: []byte{0x65, 0x00, 0x00, 0x06, 0x00}, body: []byte{0x75, 0x01, 0x01, 0x01, 0x01, 0x01}},
		{name: "AttentionResponse", tag: []byte{0x65, 0x00, 0x00, 0xFF, 0x01}, body: []byte{0x74, 0x01, 0x01, 0x01, 0x01}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := append([]byte{0x72}, tc.tag...)
			data = append(data, tc.body...)
			buf := &Buffer{Bytes: data}

			msg, err := MessageBodyParse(buf)
			if err != nil {
				t.Fatalf("message body parse failed: %v", err)
			}
			if msg.Data == nil {
				t.Fatalf("message data is nil")
			}
			if buf.Cursor != len(data) {
				t.Fatalf("message parser consumed %d bytes, expected %d", buf.Cursor, len(data))
			}
		})
	}
}
