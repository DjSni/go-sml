package sml

import (
	"reflect"
	"strings"
	"testing"
)

func TestMessageBodyParseDispatchesKnownTags(t *testing.T) {
	tests := []struct {
		name string
		tag  uint32
		typ  interface{}
	}{
		{name: "OpenRequest", tag: MESSAGEOPENREQUEST, typ: OpenRequest{}},
		{name: "OpenResponse", tag: MESSAGEOPENRESPONSE, typ: OpenResponse{}},
		{name: "CloseRequest", tag: MESSAGECLOSEREQUEST, typ: CloseRequest{}},
		{name: "CloseResponse", tag: MESSAGECLOSERESPONSE, typ: CloseResponse{}},
		{name: "GetProfilePackRequest", tag: MESSAGEGETPROFILEPACKREQUEST, typ: GetProfilePackRequest{}},
		{name: "GetProfilePackResponse", tag: MESSAGEGETPROFILEPACKRESPONSE, typ: GetProfilePackResponse{}},
		{name: "GetProfileListRequest", tag: MESSAGEGETPROFILELISTREQUEST, typ: GetProfileListRequest{}},
		{name: "GetProfileListResponse", tag: MESSAGEGETPROFILELISTRESPONSE, typ: GetProfileListResponse{}},
		{name: "GetProcParameterRequest", tag: MESSAGEGETPROCPARAMETERREQUEST, typ: GetProcParameterRequest{}},
		{name: "GetProcParameterResponse", tag: MESSAGEGETPROCPARAMETERRESPONSE, typ: GetProcParameterResponse{}},
		{name: "SetProcParameterRequest", tag: MESSAGESETPROCPARAMETERREQUEST, typ: SetProcParameterRequest{}},
		{name: "GetListRequest", tag: MESSAGEGETLISTREQUEST, typ: GetListRequest{}},
		{name: "GetListResponse", tag: MESSAGEGETLISTRESPONSE, typ: GetListResponse{}},
		{name: "AttentionResponse", tag: MESSAGEATTENTIONRESPONSE, typ: AttentionResponse{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := buildMessageBody(tc.tag, minimalPayloadForTag(tc.tag))
			buf := &Buffer{Bytes: data}

			msg, err := MessageBodyParse(buf)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}

			if reflect.TypeOf(msg.Data) != reflect.TypeOf(tc.typ) {
				t.Fatalf("unexpected parsed type: got %T want %T", msg.Data, tc.typ)
			}

			if buf.Cursor != len(data) {
				t.Fatalf("parser consumed %d bytes, expected %d", buf.Cursor, len(data))
			}
		})
	}
}

func TestMessageBodyParseRejectsUnsupportedTag(t *testing.T) {
	data := buildMessageBody(MESSAGESETPROCPARAMETERRESPONSE, []byte{0x01})
	buf := &Buffer{Bytes: data}

	_, err := MessageBodyParse(buf)
	if err == nil || !strings.Contains(err.Error(), "Invalid message type") {
		t.Fatalf("expected invalid message type error, got %v", err)
	}
}

func buildMessageBody(tag uint32, payload []byte) []byte {
	data := []byte{0x72, 0x65, byte(tag >> 24), byte(tag >> 16), byte(tag >> 8), byte(tag)}
	return append(data, payload...)
}

func minimalPayloadForTag(tag uint32) []byte {
	switch tag {
	case MESSAGEOPENREQUEST:
		return []byte{0x77, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}
	case MESSAGEOPENRESPONSE:
		return []byte{0x76, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}
	case MESSAGECLOSEREQUEST, MESSAGECLOSERESPONSE:
		return []byte{0x71, 0x01}
	case MESSAGEGETPROFILEPACKREQUEST, MESSAGEGETPROFILELISTREQUEST:
		return []byte{0x79, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}
	case MESSAGEGETPROFILEPACKRESPONSE:
		return []byte{0x78, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}
	case MESSAGEGETPROFILELISTRESPONSE:
		return []byte{0x79, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}
	case MESSAGEGETPROCPARAMETERREQUEST, MESSAGESETPROCPARAMETERREQUEST:
		return []byte{0x75, 0x01, 0x01, 0x01, 0x01, 0x01}
	case MESSAGEGETPROCPARAMETERRESPONSE:
		return []byte{0x73, 0x01, 0x01, 0x01}
	case MESSAGEGETLISTREQUEST:
		return []byte{0x75, 0x01, 0x01, 0x01, 0x01, 0x01}
	case MESSAGEGETLISTRESPONSE:
		return []byte{0x77, 0x01, 0x01, 0x01, 0x01, 0x70, 0x01, 0x01}
	case MESSAGEATTENTIONRESPONSE:
		return []byte{0x74, 0x01, 0x01, 0x01, 0x01}
	default:
		return []byte{0x01}
	}
}
