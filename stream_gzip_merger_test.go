package iomerge

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamGzipMerger(t *testing.T) {
	m := NewStreamGzipMerger(10)
	readAllDone := make(chan bool)

	go func() {
		reader, err := gzip.NewReader(m.Out())
		assert.NoError(t, err)
		out, err := ioutil.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, "I am File 1, Line 1\nI am File 2, Line 1\n", string(out))
		readAllDone <- true
	}()

	go func() {
		in := m.In()
		in <- bytes.NewReader([]byte("I am File 1, Line 1\n"))
		in <- bytes.NewReader([]byte("I am File 2, Line 1\n"))
		m.Close()
	}()

	assert.NoError(t, m.Wait())
	<-readAllDone
}

func TestStreamGzipMergerImplementsMerger(t *testing.T) {
	assert.Implements(t, (*Merger)(nil), NewStreamGzipMerger(10))
}
