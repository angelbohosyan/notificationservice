package providers

import (
	"math"
	"time"
)

func Retry(calls int, resolution time.Duration, retry func() error) error {
	for i := 0; i < calls; i++ {
		if err := retry(); err != nil {
			duration := time.Duration(math.Pow(float64(2), float64(i*2)) * 10)
			time.Sleep(duration * resolution)
			continue
		}

		return nil
	}

	return retry()
}
