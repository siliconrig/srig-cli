package firmware

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// ihexToBin parses Intel HEX (record types 00,01,02,04; 03/05 ignored) into a
// raw image via layoutSegments.
func ihexToBin(data []byte) ([]byte, uint64, error) {
	var segs []segment
	var upper uint64
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	ln := 0
	for sc.Scan() {
		ln++
		rec := strings.TrimSpace(sc.Text())
		if rec == "" {
			continue
		}
		if rec[0] != ':' {
			return nil, 0, fmt.Errorf("line %d: record must start with ':'", ln)
		}
		b, err := hex.DecodeString(rec[1:])
		if err != nil || len(b) < 5 {
			return nil, 0, fmt.Errorf("line %d: malformed record", ln)
		}
		n := int(b[0])
		if len(b) != n+5 {
			return nil, 0, fmt.Errorf("line %d: byte-count mismatch", ln)
		}
		var sum byte
		for _, x := range b {
			sum += x
		}
		if sum != 0 {
			return nil, 0, fmt.Errorf("line %d: bad checksum", ln)
		}
		off := uint64(b[1])<<8 | uint64(b[2])
		payload := b[4 : 4+n]
		switch b[3] {
		case 0x00: // data
			segs = append(segs, segment{addr: upper + off, data: append([]byte(nil), payload...)})
		case 0x01: // EOF
			return layoutSegments(segs)
		case 0x02: // extended segment address
			if len(payload) != 2 {
				return nil, 0, fmt.Errorf("line %d: type 0x02 record needs 2 data bytes, got %d", ln, len(payload))
			}
			upper = uint64(binary.BigEndian.Uint16(payload)) << 4
		case 0x04: // extended linear address
			if len(payload) != 2 {
				return nil, 0, fmt.Errorf("line %d: type 0x04 record needs 2 data bytes, got %d", ln, len(payload))
			}
			upper = uint64(binary.BigEndian.Uint16(payload)) << 16
		case 0x03, 0x05: // start address — ignored
		default:
			return nil, 0, fmt.Errorf("line %d: unsupported record type %#02x", ln, b[3])
		}
	}
	if err := sc.Err(); err != nil {
		return nil, 0, err
	}
	return nil, 0, errors.New("HEX file has no EOF record (:00000001FF)")
}
