package client

import (
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFlashFirmwareBytes(t *testing.T) {
	var gotField, gotName string
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, params, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		mr := multipart.NewReader(r.Body, params["boundary"])
		part, err := mr.NextPart()
		if err != nil {
			t.Errorf("no part: %v", err)
			return
		}
		gotField = part.FormName()
		gotName = part.FileName()
		gotBody, _ = io.ReadAll(part)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	c := New(srv.URL, "sk_test")
	if err := c.FlashFirmwareBytes("sess_1", "firmware.bin", []byte{1, 2, 3, 4}); err != nil {
		t.Fatalf("FlashFirmwareBytes: %v", err)
	}
	if gotField != "firmware" {
		t.Errorf("form field = %q, want firmware", gotField)
	}
	if gotName != "firmware.bin" {
		t.Errorf("filename = %q, want firmware.bin", gotName)
	}
	if !strings.EqualFold(string(gotBody), string([]byte{1, 2, 3, 4})) {
		t.Errorf("body = % x, want 01 02 03 04", gotBody)
	}
}

func TestIsTransient(t *testing.T) {
	cases := []struct {
		err  error
		want bool
	}{
		{&APIError{StatusCode: 503, Message: "no available board"}, true},
		{&APIError{StatusCode: 500, Message: "internal error"}, true},
		{&APIError{StatusCode: 402, Message: "insufficient credits"}, false},
		{&APIError{StatusCode: 409, Message: "concurrent session limit reached"}, false},
		{&APIError{StatusCode: 401, Message: "unauthorized"}, false},
		{&APIError{StatusCode: 400, Message: "board_type is required"}, false},
		{errors.New("dial tcp: connection refused"), true},
	}
	for _, tc := range cases {
		if got := IsTransient(tc.err); got != tc.want {
			t.Errorf("IsTransient(%v) = %v, want %v", tc.err, got, tc.want)
		}
	}
}
