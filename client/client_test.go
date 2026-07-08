package client

import (
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
