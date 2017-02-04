package source

import (
	"context"
	"time"
)

type nilSource struct {
	duration time.Duration
}

// Nil returns a Source that doesn't read any data. Instead it closes
// the channels after specified duration. It can be used as a mock Source.
func Nil(duration time.Duration) Source {
	return &nilSource{
		duration: duration,
	}
}

// Read doesn't read anything.
func (s *nilSource) Read(ctx context.Context, items chan<- interface{}, errs chan<- error) {
	time.Sleep(s.duration)
	close(items)
	close(errs)
}