package retry

import "errors"

var (
	ErrRetryTooManyTimes = errors.New("too many times retried")
)

func Retry(fn func() error, times int) error {
	for i := 0; i < times; i++ {
		if err := fn(); err != nil {
			continue
		}
		return nil
	}
	return ErrRetryTooManyTimes
}
