package goparty

import (
	"errors"
	"time"
)

var (
	ErrTimeout = errors.New("execute timeout")
)

func ExeWithTimeout(timeout time.Duration, task func() error) error {
	ch := make(chan error, 1)
	go func() {ch <- task()}()
	select {
	case ret := <-ch:
		return ret
	case <-time.After(timeout):
		return ErrTimeout
	}
}