package iomerge

import "io"

// StreamMerger merges multiple io.Readers together in the order they are placed on the "in"
// channel and allows for reading them on the out channel.
type StreamMerger struct {
	*baseMerger
}

// NewStreamMerger creates and initializes a new StreamMerger.
func NewStreamMerger(inChanSize int) *StreamMerger {
	pipeReader, pipeWriter := io.Pipe()
	baseMerger := newBaseMerger(inChanSize)
	baseMerger.out = pipeReader

	go baseMerger.copyAll(pipeWriter)

	return &StreamMerger{
		baseMerger: baseMerger,
	}
}
