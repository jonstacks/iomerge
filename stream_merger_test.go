package iomerge

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// m.Read(func(r io.Reader) error {
//     // Do something with reader
//     out, err := ioutil.ReadAll(r)
//     return err
//})

// m.Write(func(in chan<- io.Reader) {
//
// })

func TestStreamMerger(t *testing.T) {
	m := NewStreamMerger(10)

	m.Read(func(r io.Reader) error {
		out, err := ioutil.ReadAll(r)
		assert.Equal(t, `{"msg": "hello world"}{"msg": "another hello"}`, string(out))
		return err
	})

	m.Write(func(in chan<- io.Reader) {
		in <- bytes.NewReader([]byte(`{"msg": "hello world"}`))
		in <- bytes.NewReader([]byte(`{"msg": "another hello"}`))
	})

	assert.NoError(t, m.Wait())
}

func TestStreamMergerImplementsMerger(t *testing.T) {
	assert.Implements(t, (*Merger)(nil), NewStreamMerger(10))
}
