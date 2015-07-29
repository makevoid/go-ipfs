package chunk

import (
	"fmt"
	"hash/fnv"
	"io"
	"math"

	"github.com/whyrusleeping/chunker"
)

type Rabin struct {
	avChunkSize int
}

func NewRabin(avgBlkSize int) *Rabin {
	return &Rabin{avChunkSize: avgBlkSize}
}

func (mr *Rabin) Split(r io.Reader) (<-chan []byte, <-chan error) {
	errs := make(chan error, 1)

	pol, err := chunker.RandomPolynomial()
	if err != nil {
		errs <- err
		close(errs)
		return nil, errs
	}

	nbits := uint(math.Log2(float64(mr.avChunkSize)))
	h := fnv.New32a()
	ch := chunker.New(r, pol, h, nbits)

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
			fmt.Println("chunk returned", chunk.Length, len(chunk.Data))
			out <- chunk.Data
		}
	}()
	return out, errs
}
