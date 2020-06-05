package iomerge

import (
	"compress/gzip"
	"io"
)

// StreamGzipMerger merges multiple io.Readers together in the order they are placed on the "in"
// channel, compresses them with gzip, and allows for reading the compressed output on the out channel.
type StreamGzipMerger struct {
	*baseMerger
}

// NewStreamGzipMerger creates and initializes a new StreamGzipMerger.
func NewStreamGzipMerger(inChanSize int) *StreamGzipMerger {
	pipeReader, pipeWriter := io.Pipe()
	baseMerger := newBaseMerger(inChanSize)
	baseMerger.out = pipeReader

	go func() {
		defer pipeWriter.Close()
		baseMerger.copyAll(gzip.NewWriter(pipeWriter))
	}()

	return &StreamGzipMerger{
		baseMerger: baseMerger,
	}
}
