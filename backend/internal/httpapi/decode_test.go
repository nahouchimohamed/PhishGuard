package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSONOK(t *testing.T) {
	var got struct {
		Email string `json:"email"`
	}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"email":"a@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	if err := decodeJSON(req, &got); err != nil {
		t.Fatal(err)
	}
	if got.Email != "a@example.com" {
		t.Fatalf("got=%q", got.Email)
	}
}

func TestDecodeJSONRejectsUnknownField(t *testing.T) {
	var got struct {
		Email string `json:"email"`
	}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"email":"a@example.com","extra":1}`))
	req.Header.Set("Content-Type", "application/json")
	if err := decodeJSON(req, &got); err == nil {
		t.Fatal("expected error")
	}
}

func TestDecodeJSONRejectsBadContentType(t *testing.T) {
	var got map[string]any
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "text/plain")
	if err := decodeJSON(req, &got); err == nil {
		t.Fatal("expected error")
	}
}

func TestDecodeJSONCapsBody(t *testing.T) {
	var got map[string]any
	body := `{"x":"` + strings.Repeat("a", maxJSONBody) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	if err := decodeJSON(req, &got); err == nil {
		t.Fatal("expected error for oversized body")
	}
}

func TestDecodeJSONRejectsTrailingJunk(t *testing.T) {
	var got map[string]any
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"ok":true}{"extra":1}`))
	req.Header.Set("Content-Type", "application/json")
	if err := decodeJSON(req, &got); err == nil {
		t.Fatal("expected error for trailing JSON")
	}
}

func TestDecodeJSONAllowsCharset(t *testing.T) {
	var got struct {
		Email string `json:"email"`
	}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"email":"a@example.com"}`))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err := decodeJSON(req, &got); err != nil {
		t.Fatal(err)
	}
	if got.Email != "a@example.com" {
		t.Fatalf("got=%q", got.Email)
	}
}
