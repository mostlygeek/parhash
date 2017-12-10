package parhash

import (
	"crypto/md5"
	"crypto/sha1"
	"hash/fnv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
