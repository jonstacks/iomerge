package iomerge

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamGzipMerger(t *testing.T) {
	m := NewStreamGzipMerger(10)

	m.Read(func(r io.Reader) error {
		reader, err := gzip.NewReader(r)
		assert.NoError(t, err)
		out, err := ioutil.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, "I am File 1, Line 1\nI am File 2, Line 1\n", string(out))
		return err
	})

	m.Write(func(in chan<- io.Reader) {
		in <- bytes.NewReader([]byte("I am File 1, Line 1\n"))
		in <- bytes.NewReader([]byte("I am File 2, Line 1\n"))
	})

	assert.NoError(t, m.Wait())
}

func TestStreamGzipMergerImplementsMerger(t *testing.T) {
	assert.Implements(t, (*Merger)(nil), NewStreamGzipMerger(10))
}
