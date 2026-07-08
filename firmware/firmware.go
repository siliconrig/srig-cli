// Package firmware converts user-supplied firmware artifacts (ELF, Intel HEX,
// raw binary) into the flashable bytes the pod expects for a given board.
package firmware

import (
	"bytes"
	"debug/elf"
	"errors"
	"fmt"
	"io"
)

type segment struct {
	addr uint64
	data []byte
}

// layoutSegments lays load segments out into one contiguous image starting at
// the lowest address, zero-filling gaps (matching `objcopy -O binary`).
func layoutSegments(segs []segment) ([]byte, uint64, error) {
	if len(segs) == 0 {
		return nil, 0, errors.New("no loadable segments")
	}
	min := segs[0].addr
	var max uint64
	for _, s := range segs {
		if s.addr < min {
			min = s.addr
		}
		if end := s.addr + uint64(len(s.data)); end > max {
			max = end
		}
	}
	out := make([]byte, max-min)
	for _, s := range segs {
		copy(out[s.addr-min:], s.data)
	}
	return out, min, nil
}

// elfToBin extracts PT_LOAD segments (by physical/load address) into a raw image.
func elfToBin(data []byte) ([]byte, uint64, error) {
	f, err := elf.NewFile(bytes.NewReader(data))
	if err != nil {
		return nil, 0, fmt.Errorf("parse ELF: %w", err)
	}
	defer f.Close()
	var segs []segment
	for _, p := range f.Progs {
		if p.Type != elf.PT_LOAD || p.Filesz == 0 {
			continue
		}
		b := make([]byte, p.Filesz)
		if _, err := io.ReadFull(p.Open(), b); err != nil {
			return nil, 0, fmt.Errorf("read ELF segment: %w", err)
		}
		segs = append(segs, segment{addr: p.Paddr, data: b})
	}
	if len(segs) == 0 {
		return nil, 0, errors.New("ELF has no loadable segments")
	}
	return layoutSegments(segs)
}
