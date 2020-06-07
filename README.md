# iomerge

![CI](https://github.com/jonstacks/iomerge/workflows/CI/badge.svg)
[![codecov](https://codecov.io/gh/jonstacks/iomerge/branch/master/graph/badge.svg)](https://codecov.io/gh/jonstacks/iomerge)


A golang package for merging multiple io streams together.

## Usage

### StreamMerger

StreamMerger merges in any number of `io.Reader`s, in order, and makes
them available on a single `io.Reader`.

```golang
// Creates a new stream merger whose input channel has size 10, any
// more writes after that block until they are read from the channel.
m := NewStreamMerger(10)
m.Read(func(r io.Reader) error {
    // Read from reader and return an error when done reading
})
m.Write(func(in chan<- io.Reader) {
    // Add readers to the in channel
    in <- myReader1
    in <- myReader2
    in <- myReader3
})
err := m.Wait()
```

### StreamGzipMerger

StreamGzipMerger merges in any number of `io.Reader`s, in order, compresses
that stream data, and makes the compressed stream avaialable on a single `io.Reader`.


```golang
m := NewStreamGzipMerger(10)
m.Read(func(r io.Reader) error {
    // Read the compressed stream from the reader and return an error when done
})
m.Write(func(in chan<- io.Reader) {
    in <- myReader1
    in <- myReader2
    in <- myReader3
})
err := m.Wait()
```