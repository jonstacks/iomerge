package iomerge

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamMerger(t *testing.T) {
	m := NewStreamMerger(10)
	readAllDone := make(chan bool)

	go func() {
		out, err := ioutil.ReadAll(m.Out())
		assert.NoError(t, err)
		assert.Equal(t, `{"msg": "hello world"}{"msg": "another hello"}`, string(out))
		readAllDone <- true
	}()
	in := m.In()
	in <- bytes.NewReader([]byte(`{"msg": "hello world"}`))
	in <- bytes.NewReader([]byte(`{"msg": "another hello"}`))
	m.Close()

	assert.NoError(t, m.Wait())
	<-readAllDone
}

func TestStreamMergerImplementsMerger(t *testing.T) {
	assert.Implements(t, (*Merger)(nil), NewStreamMerger(10))
}
