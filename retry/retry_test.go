package retry

import (
	"errors"
	"testing"
)

func TestRetry(t *testing.T) {
	i := 0
	fn := func() error {
		if i == 2 {
			return nil
		}
		i++
		return errors.New("not 2")
	}
	err := Retry(fn, 3)
	if err != nil {
		panic(err)
	}
}
