package util

import (
	"errors"
	"strconv"
	"time"
)

var (
	// ErrInvalidSyntax is returned when the input
	// value has an invalid syntax.
	ErrInvalidSyntax = errors.New("invalid syntax")
)

// ParseMemoryStr takes a memory value string s
// (5M)
func ParseMemoryStr(s string) (int64, error) {
	ln := len(s)
	if ln == 0 {
		return 0, nil
	}

	if s[0] < '0' || s[0] > '9' {
		return 0, ErrInvalidSyntax
	}

	if s[ln-1] == 'B' || s[ln-1] == 'b' {
		ln -= 1
	}

	mult := int64(1)
	vStr := s

	if ln >= 2 {
		vStr = s[0 : ln-1]
		switch s[ln-1] {
		case 'k', 'K':
			mult = 1024
		case 'm', 'M':
			mult = 1024 * 1024
		case 'g', 'G':
			mult = 1024 * 1024 * 1024
		case 't', 'T':
			mult = 1024 * 1024 * 1024 * 1024
		default:
			if s[ln-1] < '0' || s[ln-1] > '9' {
				return 0, ErrInvalidSyntax
			}
			vStr = s
		}
	}

	v, err := strconv.ParseInt(vStr, 10, 64)
	return v * mult, err
}

func MeasureTime(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}
