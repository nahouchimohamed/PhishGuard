package validate

import "testing"

func TestUUID(t *testing.T) {
	tests := []struct {
		in      string
		wantErr bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", false},
		{"550E8400-E29B-41D4-A716-446655440000", false},
		{"", true},
		{"not-a-uuid", true},
		{"550e8400e29b41d4a716446655440000", true},
		{"550e8400-e29b-41d4-a716-44665544000g", true},
		{"550e8400-e29b-41d4-a716-44665544000", true},
		{"550e8400-e29b-41d4-a716-4466554400000", true},
		{"550e8400-e29b-61d4-a716-446655440000", false},
		{"  550e8400-e29b-41d4-a716-446655440000  ", false},
		{"550e8400-e29b-41d4-g716-446655440000", true},
	}
	for _, tt := range tests {
		err := UUID(tt.in)
		if (err != nil) != tt.wantErr {
			t.Fatalf("UUID(%q)=%v wantErr=%v", tt.in, err, tt.wantErr)
		}
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		in      string
		wantErr bool
	}{
		{"a@example.com", false},
		{"  a@example.com ", false},
		{"", true},
		{"not-an-email", true},
		{"@example.com", true},
	}
	for _, tt := range tests {
		err := Email(tt.in)
		if (err != nil) != tt.wantErr {
			t.Fatalf("Email(%q)=%v wantErr=%v", tt.in, err, tt.wantErr)
		}
	}
}

func TestName(t *testing.T) {
	tests := []struct {
		in      string
		min     int
		max     int
		wantErr bool
	}{
		{"Campaign", 1, 120, false},
		{"", 1, 120, true},
		{"   ", 1, 120, true},
		{stringsRepeat("x", 121), 1, 120, true},
		{stringsRepeat("x", 120), 1, 120, false},
	}
	for _, tt := range tests {
		err := Name(tt.in, tt.min, tt.max)
		if (err != nil) != tt.wantErr {
			t.Fatalf("Name(%q)=%v wantErr=%v", tt.in, err, tt.wantErr)
		}
	}
}

func stringsRepeat(s string, n int) string {
	b := make([]byte, 0, n*len(s))
	for i := 0; i < n; i++ {
		b = append(b, s...)
	}
	return string(b)
}
