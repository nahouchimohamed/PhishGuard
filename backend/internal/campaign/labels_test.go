package campaign

import (
	"testing"
	"time"
)

func TestRecipientLabel(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		o, c, r *time.Time
		want    string
	}{
		{"delivered only", nil, nil, nil, "Delivered"},
		{"opened", &now, nil, nil, "Opened"},
		{"clicked", &now, &now, nil, "Clicked"},
		{"reported wins over click", &now, &now, &now, "Reported"},
		{"reported without click", &now, nil, &now, "Reported"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RecipientLabel(tt.o, tt.c, tt.r)
			if got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}
