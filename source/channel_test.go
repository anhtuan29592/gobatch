package source

import (
	"context"
	"sync"
	"testing"

	"github.com/MasterOfBinary/gobatch/item"
)

func TestChannelSource_Read(t *testing.T) {
	size := 10
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	defer wg.Wait()

	itemsIn := make(chan interface{}, size)
	itemsOut := make(chan item.Item)
	errsOut := make(chan error)

	itemGen := item.NewMockGenerator()
	defer itemGen.Close()

	s := Channel(itemsIn)

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Read(ctx, itemGen.GetCh(), itemsOut, errsOut)
	}()

	numItems := 10
	for i := 0; i < numItems; i++ {
		itemsIn <- i
	}
	close(itemsIn)

	i := 0
	for item := range itemsOut {
		if i > numItems-1 {
			t.Fatalf("items in itemsOut > %v", i)
		}

		if item.Get() != i {
			t.Errorf("itemsOut <- %v, want %v", item, i)
		}

		i++
	}

	if i < numItems {
		t.Errorf("items in itemsOut < %v", i)
	}
}
