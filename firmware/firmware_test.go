package firmware

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestLayoutSegments(t *testing.T) {
	tests := []struct {
		name     string
		segs     []segment
		wantBase uint64
		want     []byte
	}{
		{
			name:     "single segment",
			segs:     []segment{{addr: 0x08000000, data: []byte{1, 2, 3}}},
			wantBase: 0x08000000,
			want:     []byte{1, 2, 3},
		},
		{
			name: "gap is zero-filled and base is the lowest addr",
			segs: []segment{
				{addr: 0x08000000, data: []byte{0xAA, 0xBB}},
				{addr: 0x08000004, data: []byte{0xCC}},
			},
			wantBase: 0x08000000,
			want:     []byte{0xAA, 0xBB, 0x00, 0x00, 0xCC},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, base, err := layoutSegments(tt.segs)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if base != tt.wantBase {
				t.Errorf("base = %#x, want %#x", base, tt.wantBase)
			}
			if !bytes.Equal(out, tt.want) {
				t.Errorf("out = % x, want % x", out, tt.want)
			}
		})
	}
}

// buildELF32 writes a minimal little-endian ARM ELF with one PT_LOAD segment
// at paddr=vaddr=base carrying payload. Keeps the test hermetic (no arm toolchain).
func buildELF32(base uint32, payload []byte) []byte {
	return buildELF32VP(base, base, payload)
}

// buildELF32VP is like buildELF32 but writes distinct p_vaddr and p_paddr, so
// tests can prove placement is by physical (load) address, not virtual address.
func buildELF32VP(vaddr, paddr uint32, payload []byte) []byte {
	const ehsize, phentsize = 52, 32
	dataOff := uint32(ehsize + phentsize)
	buf := make([]byte, dataOff)
	copy(buf, []byte{0x7f, 'E', 'L', 'F', 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	le := binary.LittleEndian
	le.PutUint16(buf[16:], 2)         // e_type ET_EXEC
	le.PutUint16(buf[18:], 40)        // e_machine EM_ARM
	le.PutUint32(buf[20:], 1)         // e_version
	le.PutUint32(buf[24:], paddr)     // e_entry
	le.PutUint32(buf[28:], 52)        // e_phoff
	le.PutUint16(buf[40:], ehsize)    // e_ehsize
	le.PutUint16(buf[42:], phentsize) // e_phentsize
	le.PutUint16(buf[44:], 1)         // e_phnum
	ph := buf[52:84]
	le.PutUint32(ph[0:], 1)                     // p_type PT_LOAD
	le.PutUint32(ph[4:], dataOff)               // p_offset
	le.PutUint32(ph[8:], vaddr)                 // p_vaddr
	le.PutUint32(ph[12:], paddr)                // p_paddr
	le.PutUint32(ph[16:], uint32(len(payload))) // p_filesz
	le.PutUint32(ph[20:], uint32(len(payload))) // p_memsz
	le.PutUint32(ph[24:], 5)                    // p_flags R+X
	le.PutUint32(ph[28:], 4)                    // p_align
	return append(buf, payload...)
}

func TestElfToBin(t *testing.T) {
	payload := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	elf := buildELF32(0x08000000, payload)
	out, base, err := elfToBin(elf)
	if err != nil {
		t.Fatalf("elfToBin: %v", err)
	}
	if base != 0x08000000 {
		t.Errorf("base = %#x, want 0x08000000", base)
	}
	if !bytes.Equal(out, payload) {
		t.Errorf("out = % x, want % x", out, payload)
	}
}

func TestElfToBinNoLoadable(t *testing.T) {
	elf := buildELF32(0x08000000, nil) // p_filesz == 0
	if _, _, err := elfToBin(elf); err == nil {
		t.Fatal("expected error for ELF with no loadable segments")
	}
}

func TestElfToBinUsesPaddr(t *testing.T) {
	payload := []byte{0x11, 0x22}
	elf := buildELF32VP(0x20000000 /*vaddr*/, 0x08000000 /*paddr*/, payload)
	_, base, err := elfToBin(elf)
	if err != nil {
		t.Fatalf("elfToBin: %v", err)
	}
	if base != 0x08000000 {
		t.Errorf("base = %#x, want 0x08000000 (paddr, not vaddr)", base)
	}
}
