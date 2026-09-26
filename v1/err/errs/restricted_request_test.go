package errs

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorRestrictedRequest(t *testing.T) {
	tests := []struct {
		name      string
		userID    fields.Key
		source    fields.Key
		key       fields.Key
		command   fields.Command
		requestID fields.RequestID
		id        fields.TracebackID
		want      []byte
		wantErr   bool
	}{
		{"nil keys", nil, nil, nil, fields.Handshake, 0, tracebackIDOne, encoded(17, tracebackIDOne), false},
		{"empty id", fields.Key{}, nil, nil, fields.Answer, 0, tracebackIDOne, encoded(17, tracebackIDOne), false},
		{"empty source", nil, fields.Key{}, nil, fields.Read, 0, tracebackIDOne, encoded(17, tracebackIDOne), false},
		{"empty key", nil, nil, fields.Key{}, fields.Write, 0, tracebackIDOne, encoded(17, tracebackIDOne), false},
		{"full request", fields.Key("user"), fields.Key("192.168.1.1"), fields.Key("key"), fields.Delete, 1, tracebackIDOne, encoded(17, tracebackIDOne), false},
		{"max request id", fields.Key("user"), fields.Key("192.168.1.1"), fields.Key("key"), fields.Delete, math.MaxUint32, tracebackIDMax, encoded(17, tracebackIDMax), false},
		{"zero value", nil, nil, nil, fields.Handshake, 0, fields.TracebackID{}, encoded(17, fields.TracebackID{}), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorRestrictedRequestWithTracebackID(tt.userID, tt.source, tt.key, tt.command, tt.requestID, tt.id)
			assertProtocolError(t, e, fields.RestrictedRequest, tt.id, tt.want, tt.wantErr)
		})
	}
}
