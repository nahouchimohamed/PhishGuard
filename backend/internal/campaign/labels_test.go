package campaign

import (
	"testing"
	"time"
)

func TestRecipientLabel(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		d, o, c, r *time.Time
		want       string
	}{
		{"pending", nil, nil, nil, nil, "Pending"},
		{"delivered only", &now, nil, nil, nil, "Delivered"},
		{"opened", &now, &now, nil, nil, "Opened"},
		{"clicked", &now, &now, &now, nil, "Clicked"},
		{"reported wins over click", &now, &now, &now, &now, "Reported"},
		{"reported without click", &now, &now, nil, &now, "Reported"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RecipientLabel(tt.d, tt.o, tt.c, tt.r)
			if got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}
