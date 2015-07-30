package chunk

import (
	"hash/fnv"
	"io"

	"github.com/whyrusleeping/chunker"
)

var IpfsRabinPoly = chunker.Pol(17437180132763653)

type Rabin struct {
	avChunkSize uint64
}

func NewRabin(avgBlkSize uint64) *Rabin {
	return &Rabin{avChunkSize: avgBlkSize}
}

func (mr *Rabin) Split(r io.Reader) (<-chan []byte, <-chan error) {
	errs := make(chan error, 1)

	h := fnv.New32a()
	ch := chunker.New(r, IpfsRabinPoly, h, mr.avChunkSize)
	ch.MinSize = mr.avChunkSize / 3 //tweaking to get a better average size

	out := make(chan []byte, 16)
	go func() {
		defer close(out)
		defer close(errs)

		for {
			chunk, err := ch.Next()
			if err != nil {
				if err != io.EOF {
					errs <- err
				}
				return
			}

			out <- chunk.Data
		}
	}()
	return out, errs
}
