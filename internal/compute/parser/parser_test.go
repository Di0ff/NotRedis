package parser

import (
	"NotRedis/internal/myError"
	"go.uber.org/zap"
	"testing"
)

func TestParse(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name    string
		input   string
		want    MySplitRequest
		wantErr error
	}{
		{
			name:    "Valid SET",
			input:   "SET key value",
			want:    MySplitRequest{Type: "SET", Key: "key", Value: "value"},
			wantErr: nil,
		},
		{
			name:    "Valid GET",
			input:   "GET key",
			want:    MySplitRequest{Type: "GET", Key: "key"},
			wantErr: nil,
		},
		{
			name:    "Valid DEL",
			input:   "DEL key",
			want:    MySplitRequest{Type: "DEL", Key: "key"},
			wantErr: nil,
		},
		{
			name:    "Empty request",
			input:   "",
			want:    MySplitRequest{},
			wantErr: myError.EmptyRequest,
		},
		{
			name:    "Only spaces",
			input:   "   ",
			want:    MySplitRequest{},
			wantErr: myError.EmptyRequest,
		},
		{
			name:    "Unknown command",
			input:   "UNKNOWN key",
			want:    MySplitRequest{},
			wantErr: myError.UnknownRequest,
		},
		{
			name:    "SET with wrong args",
			input:   "SET key",
			want:    MySplitRequest{},
			wantErr: myError.SetFail,
		},
		{
			name:    "GET with wrong args",
			input:   "GET key value",
			want:    MySplitRequest{},
			wantErr: myError.GetFail,
		},
		{
			name:    "DEL with wrong args",
			input:   "DEL",
			want:    MySplitRequest{},
			wantErr: myError.DelFail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(logger)
			got, err := p.Parse(tt.input)

			if err != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got.Type != tt.want.Type || got.Key != tt.want.Key || got.Value != tt.want.Value {
					t.Errorf("Parse() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
