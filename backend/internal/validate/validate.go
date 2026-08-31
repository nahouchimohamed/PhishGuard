package validate

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"
)

// UUID checks a canonical 8-4-4-4-12 hex UUID (any version/variant).
func UUID(s string) error {
	s = strings.TrimSpace(s)
	if len(s) != 36 {
		return fmt.Errorf("invalid uuid")
	}
	for i, c := range s {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return fmt.Errorf("invalid uuid")
			}
		default:
			if !isHex(c) {
				return fmt.Errorf("invalid uuid")
			}
		}
	}
	return nil
}

func isHex(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

// Email performs a lightweight address parse (not full RFC enforcement).
func Email(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("email is required")
	}
	addr, err := mail.ParseAddress(s)
	if err != nil || addr.Address == "" {
		return fmt.Errorf("invalid email")
	}
	if strings.ContainsAny(addr.Address, " <>") {
		return fmt.Errorf("invalid email")
	}
	return nil
}

// Name checks a trimmed display/campaign name length in runes.
func Name(s string, min, max int) error {
	s = strings.TrimSpace(s)
	n := utf8.RuneCountInString(s)
	if n < min || n > max {
		return fmt.Errorf("name length must be between %d and %d", min, max)
	}
	return nil
}
