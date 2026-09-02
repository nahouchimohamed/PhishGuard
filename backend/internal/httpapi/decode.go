package httpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const maxJSONBody = 1 << 20 // 1 MiB

// decodeJSON reads a JSON object body with a size cap, known fields only,
// and an application/json content type (charset optional).
func decodeJSON(r *http.Request, dst any) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		media := strings.TrimSpace(strings.Split(ct, ";")[0])
		if !strings.EqualFold(media, "application/json") {
			return fmt.Errorf("content-type must be application/json")
		}
	}
	dec := json.NewDecoder(io.LimitReader(r.Body, maxJSONBody+1))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	// Reject trailing junk / a second value.
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("trailing data")
	}
	return nil
}
