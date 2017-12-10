package parhash

import (
	"hash"
	"io"
)

type Parhash struct {
	hashes map[string]hash.Hash
}

func New() *Parhash {
	return &Parhash{hashes: make(map[string]hash.Hash)}
}

func (p *Parhash) Add(id string, h hash.Hash) error {
	p.hashes[id] = h
	return nil
}

func (p *Parhash) Reset() {
	for _, h := range p.hashes {
		h.Reset()
	}
}

func (p *Parhash) GetSum(id string) []byte {

	hash, ok := p.hashes[id]
	if !ok {
		return []byte{}
	}

	return hash.Sum(nil)
}

func (p *Parhash) Write(b []byte) (n int, err error) {
	for _, h := range p.hashes {
		n, err = h.Write(b)
		if err != nil {
			return
		}

		if n != len(b) {
			err = io.ErrShortWrite
			return
		}
	}

	return len(b), nil
}
