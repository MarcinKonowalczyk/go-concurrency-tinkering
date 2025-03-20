package main

import "fmt"

type ringBuffer[T any] struct {
	input  chan T
	Dropped chan T
}

func newRingBuffer[T any](size int) *ringBuffer[T] {
	return &ringBuffer[T]{
		input:  make(chan T, size),
		Dropped: nil,
	}
}

func (rb *ringBuffer[T]) add(item T) {
	select {
	case rb.input <- item:
	default:
		// failed to make space in the buffer, so we need to read from the output channel
		value := <-rb.input
		select {
			// non-blocking send to the dropped channel
			case rb.Dropped <- value:
			default:
		}
		rb.input <- item
	}
}

func (rb *ringBuffer[T]) into_slice() []T {
	slice := make([]T, 0)
	loop: for {
		select {
		case item := <-rb.input:
			slice = append(slice, item)
		default:
			break loop
		}
	}
	return slice
}

func (rb *ringBuffer[T]) slice() []T {
	slice := rb.into_slice()
	// enqueue the items back to the input channel
	for _, item := range slice {
		rb.input <- item
	}
	return slice
}

func (rb *ringBuffer[T]) close() {
	close(rb.input)
}

func (rb *ringBuffer[T]) recieveDropped() {
	// Setup the dropped channel
	rb.Dropped = make(chan T, 1)
}

// Non-blocking read from the dropped channel
func (rb *ringBuffer[T]) maybeDropped() *T {
	select {
	case item := <-rb.Dropped:
		return &item
	default:
	}
	return nil
}

func main() {
	rb := newRingBuffer[int](4)
	rb.recieveDropped()

	defer rb.close()

	for i := 0; i < 10; i++ {
		rb.add(i);
		fmt.Println(rb.slice())
		dropped := rb.maybeDropped()
		if dropped != nil {
			fmt.Println("Dropped:", *dropped)
		}
	}

}
