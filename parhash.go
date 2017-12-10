package parhash

import (
	"hash"
	"runtime"
	"sync"
)

type workMessage struct {
	hash hash.Hash
	data []byte
}

type Parhash struct {
	sync.Mutex
	wg       sync.WaitGroup
	hashes   map[string]hash.Hash
	workChan chan *workMessage
}

func New() *Parhash {
	parhash := &Parhash{
		hashes:   make(map[string]hash.Hash),
		workChan: make(chan *workMessage, 10),
	}

	// start some goroutines to process parhash's work
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(p *Parhash) {
			for msg := range p.workChan {
				msg.hash.Write(msg.data)
				p.wg.Done()
			}
		}(parhash)
	}

	return parhash
}

// Stop shutdown Parhashes goroutines
func (p *Parhash) Stop() {
	close(p.workChan)
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
		msg := &workMessage{hash: h, data: b}
		p.wg.Add(1)
		p.workChan <- msg
	}

	p.wg.Wait()
	return len(b), nil
}
