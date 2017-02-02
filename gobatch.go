package gobatch

import (
	"context"

	"github.com/MasterOfBinary/gobatch/processor"
	"github.com/MasterOfBinary/gobatch/source"
)

type Batch interface {
	Go(ctx context.Context, s source.Source, p processor.Processor) <-chan error

	Done() <-chan struct{}
}