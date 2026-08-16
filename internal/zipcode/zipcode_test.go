package zipcode

import (
	"testing"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{name: "valid", raw: `"29902555"`, want: "29902555"},
		{name: "short", raw: `"2990255"`, wantErr: true},
		{name: "letters", raw: `"29902abc"`, wantErr: true},
		{name: "number", raw: `29902555`, wantErr: true},
		{name: "missing", raw: ``, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Validate(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				if !IsValidationError(err) {
					t.Fatalf("expected validation error, got %T", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
