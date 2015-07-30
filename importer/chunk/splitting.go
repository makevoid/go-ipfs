// package chunk implements streaming block splitters
package chunk

import (
	"io"

	"github.com/ipfs/go-ipfs/util"
)

var log = util.Logger("chunk")

var DefaultBlockSize = 1024 * 256
var DefaultSplitter = &SizeSplitter{Size: DefaultBlockSize}

type BlockSplitter interface {
	Split(r io.Reader) (<-chan []byte, <-chan error)
}

type SizeSplitter struct {
	Size int
}

func (ss *SizeSplitter) Split(r io.Reader) (<-chan []byte, <-chan error) {
	out := make(chan []byte)
	errs := make(chan error)
	go func() {
		defer close(out)
		defer close(errs)

		// all-chunks loop (keep creating chunks)
		for {
			// log.Infof("making chunk with size: %d", ss.Size)
			chunk := make([]byte, ss.Size)
			nread, err := io.ReadFull(r, chunk)
			if nread > 0 {
				// log.Infof("sending out chunk with size: %d", sofar)
				out <- chunk[:nread]
			}
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return
			}
			if err != nil {
				errs <- err
				return
			}
		}
	}()
	return out, errs
}
