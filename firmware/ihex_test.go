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
