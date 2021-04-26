package util

import (
	"strconv"
	"time"
)

func ParseMemoryStr(s string) (int64, error) {
	ln := len(s)
	if ln == 0 {
		return 0, nil
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
			vStr = s
		}
	}

	v, err := strconv.ParseInt(vStr, 10, 64)
	return v * mult, err
}

func MeasureExecTime(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}
