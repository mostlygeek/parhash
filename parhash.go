// Package parhash presents tools for making it easy to generate multiple
// hash sums over the same data. It does this in parallel, spreading the
// computation work over all available CPUs.
//
// Parhash has a speed advantage over serial operations when working
// with multiple hash algorithms over large data. Run the benchmark to see
// the difference over serial operation on your machine:
//     go test -bench .
//
// On my machine, i7 3.5Ghz dual core I get a speed up of about 240%
//   BenchmarkSerial-4            500           3577258 ns/op
//   BenchmarkParallel-4         1000           1458583 ns/op
//
// For smaller data sets parhash won't see significant performance
// advantages due to the overhead of using channels and goroutines
// distribute the work.
package parhash

import (
	"hash"
	"io"
	"runtime"
	"sync"
)

var (
	workQueue chan *hasher
)

// to parallelize work the package maintains one worker per CPU
func init() {
	numCPUs := runtime.NumCPU()
	workQueue = make(chan *hasher, numCPUs+1)

	// create a worker for each CPU
	for i := 0; i < numCPUs; i++ {
		go func() {
			for h := range workQueue {
				h.hash.Write(*h.data)
				h.wg.Done()
			}
		}()
	}
}

type hasher struct {
	hash hash.Hash

	// points to data to be written to the hash
	data *[]byte

	// share a WaitGroup with to trigger the Done()
	wg *sync.WaitGroup
}

// Parhash distributes hash sum computation over all CPU cores. It
// is an io.Writer so it works well with data handling tools in the
// Go standard library
type Parhash struct {
	io.Writer
	wg     sync.WaitGroup
	hashes []*hasher
}

// Create a new parallel hasher.
func New() *Parhash {
	return &Parhash{hashes: make([]*hasher, 0, 2)}
}

// Add a new hash to be updated in parallel. It returns the same
// hash that was provided to allow for this usage pattern:
//    h1 := p.Add(md5.New())
//    h2 := p.Add(sha1.New())
func (p *Parhash) Add(h hash.Hash) hash.Hash {
	p.hashes = append(p.hashes, &hasher{
		wg:   &p.wg,
		hash: h,
		data: nil,
	})

	return h
}

func (p *Parhash) Write(b []byte) (n int, err error) {
	for _, hasher := range p.hashes {
		p.wg.Add(1)
		hasher.data = &b
		workQueue <- hasher
	}

	p.wg.Wait()
	return len(b), nil
}
