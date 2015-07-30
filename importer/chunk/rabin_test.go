package chunk

import (
	"bytes"
	"fmt"
	"github.com/ipfs/go-ipfs/blocks"
	"github.com/ipfs/go-ipfs/blocks/key"
	"github.com/ipfs/go-ipfs/util"
	"testing"
)

func TestRabinChunking(t *testing.T) {
	data := make([]byte, 1024*1024*16)
	util.NewTimeSeededRand().Read(data)

	r := NewRabin(1024 * 256)

	var chunks [][]byte
	blks, errs := r.Split(bytes.NewReader(data))

loop:
	for {
		select {
		case blk, ok := <-blks:
			if !ok {
				break loop
			}
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

func chunkData(t *testing.T, data []byte) map[key.Key]*blocks.Block {
	r := NewRabin(1024 * 256)

	blkmap := make(map[key.Key]*blocks.Block)
	blks, errs := r.Split(bytes.NewReader(data))

loop:
	for {
		select {
		case blk, ok := <-blks:
			if !ok {
				break loop
			}
			b := blocks.NewBlock(blk)
			blkmap[b.Key()] = b

		case err, ok := <-errs:
			if !ok {
				continue
			}
			t.Fatal(err)
		}
	}

	return blkmap
}

func TestRabinChunkReuse(t *testing.T) {
	data := make([]byte, 1024*1024*16)
	util.NewTimeSeededRand().Read(data)

	ch1 := chunkData(t, data[1000:])
	ch2 := chunkData(t, data)

	var extra int
	for k, _ := range ch2 {
		_, ok := ch1[k]
		if !ok {
			extra++
		}
	}

	if extra > 2 {
		t.Fatal("too many spare chunks made")
	}
	if extra == 2 {
		t.Log("why did we get two extra blocks?")
	}
}
