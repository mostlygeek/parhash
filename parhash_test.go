package parhash

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleParhash_usage() {
	p := New()

	// Add returns the same hash to make creation and assignment more concise
	hash1 := p.Add(md5.New())
	hash2 := p.Add(sha1.New())

	// Parhash is an io.Writer
	fmt.Fprintf(p, "Hello World")
	fmt.Printf("MD5 : %s\n", hex.EncodeToString(hash1.Sum(nil)))
	fmt.Printf("SHA1: %s", hex.EncodeToString(hash2.Sum(nil)))

	// Output:
	// MD5 : b10a8db164e0754105b7a99be72e3fe5
	// SHA1: 0a4d55a8d778e5022fab701977c5d840bbc486d0
}

func TestParhash(t *testing.T) {
	assert := assert.New(t)

	p := New()

	pMD5 := p.Add(md5.New())
	pSHA1 := p.Add(sha1.New())
	pFNV := p.Add(fnv.New64a())

	// for checking...
	m := md5.New()
	s1 := sha1.New()
	f := fnv.New64a()

	testData := []string{
		"hello", "these are ", "a bunch",
		"of strings to write into the hashes",
	}

	for _, s := range testData {
		data := []byte(s)

		// write all the data into the hashes
		p.Write(data)

		m.Write(data)
		s1.Write(data)
		f.Write(data)
	}

	assert.Equal(m.Sum(nil), pMD5.Sum(nil))
	assert.Equal(s1.Sum(nil), pSHA1.Sum(nil))
	assert.Equal(f.Sum(nil), pFNV.Sum(nil))
}

// writeSerial is only used for benchmarking to contrast
// performance differences between serial and parallel hashing
func (p *Parhash) writeSerial(b []byte) (n int, err error) {
	for _, hasher := range p.hashes {
		hasher.hash.Write(b)
	}

	return len(b), nil
}

func BenchmarkSerial(b *testing.B) {
	p := New()
	p.Add(md5.New())
	p.Add(sha1.New())
	p.Add(fnv.New64a())

	data := []byte(strings.Repeat("A", 1024*1024))
	for i := 0; i < b.N; i++ {
		p.writeSerial(data)
	}
}

func BenchmarkParallel(b *testing.B) {

	p := New()
	p.Add(md5.New())
	p.Add(sha1.New())
	p.Add(fnv.New64a())

	data := []byte(strings.Repeat("A", 1024*1024))

	for i := 0; i < b.N; i++ {
		p.Write(data)
	}
}
