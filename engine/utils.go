package engine

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/c2h5oh/datasize"
	"golang.org/x/time/rate"
)

func rateLimiter(rstr string) (*rate.Limiter, error) {
	var rateSize int
	rstr = strings.ToLower(strings.TrimSpace(rstr))
	switch rstr {
	case "low":
		// ~50k/s
		rateSize = 50000
	case "medium":
		// ~500k/s
		rateSize = 500000
	case "high":
		// ~1500k/s
		rateSize = 1500000
	case "unlimited", "0", "":
		// unlimited
		return rate.NewLimiter(rate.Inf, 0), nil
	default:
		var v datasize.ByteSize
		err := v.UnmarshalText([]byte(rstr))
		if err != nil {
			return nil, err
		}
		if v > 2147483647 {
			// max of int, unlimited
			return nil, errors.New("excceed int val")
		}

		rateSize = int(v)
	}
	return rate.NewLimiter(rate.Limit(rateSize), rateSize*3), nil
}

func genEnv(dir, path, hash, ttype, api string, size int64, ts int64) []string {
	env := append(os.Environ(),
		fmt.Sprintf("CLD_DIR=%s", dir),
		fmt.Sprintf("CLD_PATH=%s", path),
		fmt.Sprintf("CLD_HASH=%s", hash),
		fmt.Sprintf("CLD_TYPE=%s", ttype),
		fmt.Sprintf("CLD_RESTAPI=%s", api),
		fmt.Sprintf("CLD_SIZE=%d", size),
		fmt.Sprintf("CLD_STARTTS=%d", ts),
	)
	return env
}
