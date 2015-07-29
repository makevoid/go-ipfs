package chunk

import (
	"bytes"
	"fmt"
	"github.com/ipfs/go-ipfs/util"
	"testing"
)

func TestRabinChunking(t *testing.T) {
	data := make([]byte, 1024*1024*16)
	util.NewTimeSeededRand().Read(data)

	// really trying to get 256k
	r := NewRabin(1024 * 128)

	var chunks [][]byte
	blks, errs := r.Split(bytes.NewReader(data))

loop:
	for {
		select {
		case blk, ok := <-blks:
			if !ok {
				break loop
			}
			fmt.Printf("got block size %d\n", len(blk))
			chunks = append(chunks, blk)

		case err, ok := <-errs:
			if !ok {
				continue
			}
			t.Fatal(err)
		}
	}

	fmt.Printf("average block size: %d\n", len(data)/len(chunks))

	unchunked := bytes.Join(chunks, nil)
	if !bytes.Equal(unchunked, data) {
		fmt.Printf("%d %d\n", len(unchunked), len(data))
		t.Fatal("data was chunked incorrectly")
	}
}
