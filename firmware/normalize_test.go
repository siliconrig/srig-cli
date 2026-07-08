package firmware

import (
	"bytes"
	"strings"
	"testing"
)

func TestNormalizeRawPassthrough(t *testing.T) {
	raw := []byte{0x01, 0x02, 0x03, 0x04}
	out, info, err := Normalize(raw, "esp32-s3")
	if err != nil {
		t.Fatalf("Normalize: %v", err)
	}
	if info.Format != FormatRaw {
		t.Errorf("format = %v, want FormatRaw", info.Format)
	}
	if !bytes.Equal(out, raw) {
		t.Error("raw input must pass through byte-identical")
	}
}

func TestNormalizeStm32ELF(t *testing.T) {
	elf := buildELF32(0x08000000, []byte{0xAA, 0xBB, 0xCC})
	out, info, err := Normalize(elf, "stm32-h753")
	if err != nil {
		t.Fatalf("Normalize: %v", err)
	}
	if info.Format != FormatELF || info.BaseAddr != 0x08000000 {
		t.Errorf("info = %+v", info)
	}
	if !bytes.Equal(out, []byte{0xAA, 0xBB, 0xCC}) {
		t.Errorf("out = % x", out)
	}
}

func TestNormalizeStm32WrongBase(t *testing.T) {
	elf := buildELF32(0x08020000, []byte{0x01})
	_, _, err := Normalize(elf, "stm32-f446")
	if err == nil || !strings.Contains(err.Error(), "0x8000000") {
		t.Fatalf("expected flash-base error, got %v", err)
	}
}

func TestNormalizeELFOnNonStm32(t *testing.T) {
	elf := buildELF32(0x08000000, []byte{0x01})
	_, _, err := Normalize(elf, "esp32-s3")
	if err == nil || !strings.Contains(err.Error(), "STM32") {
		t.Fatalf("expected STM32-only error, got %v", err)
	}
}
