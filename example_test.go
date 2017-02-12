package gobatch_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MasterOfBinary/gobatch"
	"github.com/MasterOfBinary/gobatch/item"
	"github.com/MasterOfBinary/gobatch/source"
)

// printProcessor implements the processor.Processor interface.
type printProcessor struct{}

// Process processes items in batches.
func (p printProcessor) Process(ctx context.Context, items []item.Item, errs chan<- error) {
	// Process needs to close the error channel after it's done
	defer close(errs)

	toPrint := make([]interface{}, 0, len(items))

	// This processor prints all the items in a line. If items includes 5 it removes
	// it and throws an error for no reason
	for i := 0; i < len(items); i++ {
		if items[i].Get() == 5 {
			errs <- errors.New("Cannot process 5")
		} else {
			toPrint = append(toPrint, items[i].Get())
		}
	}

	fmt.Println(toPrint)
}

func Example() {
	// Create a batch processor that processes items 5 at a time
	config := gobatch.ConstantConfig(&gobatch.ConfigValues{
		MinItems: 5,
	})
	b := gobatch.New(config, 1)
	p := &printProcessor{}

	// The channel Source reads from a channel until it's closed
	ch := make(chan interface{})
	s := source.Channel(ch)

	// Go runs in the background while the main goroutine processes errors
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errs := b.Go(ctx, s, p)

	// Spawn a goroutine that simulates loading data from somewhere
	go func() {
		for i := 0; i < 20; i++ {
			time.Sleep(time.Millisecond * 10)
			ch <- i
		}
		close(ch)
	}()

	// Wait for errors. When the error channel is closed the pipeline has been
	// completely drained. Alternatively, we could wait for Done.
	var lastErr error
	for err := range errs {
		lastErr = err
	}

	fmt.Println("Finished processing.")
	if lastErr != nil {
		fmt.Println("Found error:", lastErr.Error())
	}
	// Output:
	// [0 1 2 3 4]
	// [6 7 8 9]
	// [10 11 12 13 14]
	// [15 16 17 18 19]
	// Finished processing.
	// Found error: Cannot process 5
}
