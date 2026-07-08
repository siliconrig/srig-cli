package firmware

import (
	"bytes"
	"testing"
)

func TestIhexToBinSimple(t *testing.T) {
	// 4 data bytes DE AD BE EF at 0x0000, then EOF.
	hex := ":04000000DEADBEEFC4\n:00000001FF\n"
	out, base, err := ihexToBin([]byte(hex))
	if err != nil {
		t.Fatalf("ihexToBin: %v", err)
	}
	if base != 0 {
		t.Errorf("base = %#x, want 0", base)
	}
	if !bytes.Equal(out, []byte{0xDE, 0xAD, 0xBE, 0xEF}) {
		t.Errorf("out = % x", out)
	}
}

func TestIhexToBinExtendedLinear(t *testing.T) {
	// Extended Linear Address 0x0800 -> base 0x08000000, then 2 data bytes.
	hex := ":020000040800F2\n:02000000AABB99\n:00000001FF\n"
	out, base, err := ihexToBin([]byte(hex))
	if err != nil {
		t.Fatalf("ihexToBin: %v", err)
	}
	if base != 0x08000000 {
		t.Errorf("base = %#x, want 0x08000000", base)
	}
	if !bytes.Equal(out, []byte{0xAA, 0xBB}) {
		t.Errorf("out = % x", out)
	}
}

func TestIhexToBinBadChecksum(t *testing.T) {
	hex := ":04000000DEADBEEF00\n:00000001FF\n" // wrong checksum
	if _, _, err := ihexToBin([]byte(hex)); err == nil {
		t.Fatal("expected checksum error")
	}
}

func TestIhexToBinExtendedSegment(t *testing.T) {
	// Extended Segment Address 0x1000 -> base 0x10000 (<<4), then 1 data byte.
	hex := ":020000021000EC\n:0100000042BD\n:00000001FF\n"
	out, base, err := ihexToBin([]byte(hex))
	if err != nil {
		t.Fatalf("ihexToBin: %v", err)
	}
	if base != 0x10000 {
		t.Errorf("base = %#x, want 0x10000", base)
	}
	if !bytes.Equal(out, []byte{0x42}) {
		t.Errorf("out = % x", out)
	}
}

func TestIhexToBinMissingEOF(t *testing.T) {
	// A valid data record with no EOF (:00000001FF) line must be rejected.
	hex := ":04000000DEADBEEFC4\n"
	if _, _, err := ihexToBin([]byte(hex)); err == nil {
		t.Fatal("expected missing-EOF error")
	}
}

func TestIhexToBinMalformedExtAddr(t *testing.T) {
	// Type 0x04 record with only 1 data byte but a valid checksum must error,
	// not panic on the BigEndian.Uint16 read.
	hex := ":0100000408F3\n"
	if _, _, err := ihexToBin([]byte(hex)); err == nil {
		t.Fatal("expected malformed ext-addr error")
	}
}
