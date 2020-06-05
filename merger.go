package iomerge

import (
	"io"
	"sync"
)

// Merger interface is implemented by various mergers. It blocks until all data is read
// and returns an error/nil
type Merger interface {
	// Read takes a function and provides it with the io.Reader to read from. The
	// function returns an error on whether or not read was successful. The subsequent
	// call to Wait blocks until the read is done.
	Read(func(io.Reader) error)
	// Write takes a function and provides it with a in channel which can be used to write
	// io.Readers to.
	Write(func(in chan<- io.Reader))
	// Wait blocks until all input data has been merged together read by the function provided
	// to Read and returns an error.
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

func (bm *baseMerger) close() {
	bm.inCloseOnce.Do(func() {
		close(bm.in)
	})
}

func (bm *baseMerger) Wait() error {
	return <-bm.err
}

func (bm *baseMerger) Write(f func(chan<- io.Reader)) {
	go func() {
		defer bm.close()
		f(bm.in)
	}()
}

func (bm *baseMerger) Read(f func(io.Reader) error) {
	go func() {
		bm.err <- f(bm.out)
	}()
}

func (bm *baseMerger) copyAll(w io.WriteCloser) {
	defer w.Close()
	for reader := range bm.in {
		_, err := io.Copy(w, reader)
		if err != nil {
			bm.err <- err
		}
	}
}
