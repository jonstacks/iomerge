package iomerge

import (
	"io"
	"sync"
)

// Merger interface is implemented by various mergers. It blocks until all data is read
// and returns an error/nil
type Merger interface {
	// Close is a convenience function for closing the input channel. It should only
	// be called once.
	Close()
	// In returns the input channel that can be used to write readers to. It is the
	// user's responsibility to close the
	In() chan<- io.Reader
	// Out returns an io.Reader which contains the merged data.
	Out() io.Reader
	// Wait blocks until all input data has been read and merged together and then
	// returns an error.
	Wait() error
}

type baseMerger struct {
	in          chan io.Reader
	inCloseOnce sync.Once
	out         io.Reader

	err chan error
}

func newBaseMerger(inChanSize int) *baseMerger {
	return &baseMerger{
		in:  make(chan io.Reader, inChanSize),
		err: make(chan error, 1),
	}
}

func (bm *baseMerger) Close() {
	bm.inCloseOnce.Do(func() {
		close(bm.in)
	})
}

func (bm *baseMerger) In() chan<- io.Reader {
	return bm.in
}

func (bm *baseMerger) Out() io.Reader {
	return bm.out
}

func (bm *baseMerger) Wait() error {
	return <-bm.err
}

func (bm *baseMerger) copyAll(w io.WriteCloser) {
	defer w.Close()
	for reader := range bm.in {
		_, err := io.Copy(w, reader)
		if err != nil {
			bm.err <- err
		}
	}
	bm.err <- nil
}
