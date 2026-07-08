package firmware

import (
	"fmt"
	"strings"
)

type Format int

const (
	FormatRaw Format = iota // .bin / .uf2 — passed through unchanged
	FormatELF
	FormatIHEX
)

func (f Format) String() string {
	switch f {
	case FormatELF:
		return "ELF"
	case FormatIHEX:
		return "Intel HEX"
	default:
		return "raw binary"
	}
}

// Info describes what Normalize did.
type Info struct {
	Format   Format
	BaseAddr uint64 // resolved load base (ELF/HEX only)
	Size     int    // output byte count
}

const stm32FlashBase = 0x08000000

// Normalize converts firmware into the flashable image the pod expects for
// boardType. ELF and Intel HEX are converted to raw bin for STM32 boards; every
// other input passes through unchanged.
func Normalize(data []byte, boardType string) ([]byte, Info, error) {
	format := detect(data)
	if format == FormatRaw {
		return data, Info{Format: FormatRaw, Size: len(data)}, nil
	}

	if !strings.HasPrefix(boardType, "stm32-") {
		return nil, Info{}, fmt.Errorf(
			"%s input is only supported for STM32 boards; for %s provide %s",
			format, boardType, nativeHint(boardType))
	}

	var (
		bin  []byte
		base uint64
		err  error
	)
	if format == FormatELF {
		bin, base, err = elfToBin(data)
	} else {
		bin, base, err = ihexToBin(data)
	}
	if err != nil {
		return nil, Info{}, err
	}
	if base != stm32FlashBase {
		return nil, Info{}, fmt.Errorf(
			"firmware is linked at %#x, but %s flashes at %#x — link at the flash base",
			base, boardType, stm32FlashBase)
	}
	return bin, Info{Format: format, BaseAddr: base, Size: len(bin)}, nil
}

func detect(data []byte) Format {
	if len(data) >= 4 && data[0] == 0x7f && data[1] == 'E' && data[2] == 'L' && data[3] == 'F' {
		return FormatELF
	}
	for _, b := range data {
		if b == ' ' || b == '\t' || b == '\r' || b == '\n' {
			continue
		}
		if b == ':' {
			return FormatIHEX
		}
		break
	}
	return FormatRaw
}

func nativeHint(boardType string) string {
	if boardType == "rp2350" {
		return "a .uf2"
	}
	return "a .bin"
}
