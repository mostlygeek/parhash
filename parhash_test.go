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

	p.Add("md5", md5.New())
	p.Add("sha1", sha1.New())
	p.Add("fnv64a", fnv.New64a())

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

	assert.Equal(m.Sum(nil), p.GetSum("md5"))
	assert.Equal(s1.Sum(nil), p.GetSum("sha1"))
	assert.Equal(f.Sum(nil), p.GetSum("fnv64a"))
}

func BenchmarkSerial(b *testing.B) {
	p := New()
	p.Add("md5", md5.New())
	p.Add("sha1", sha1.New())
	p.Add("fnv64a", fnv.New64a())

	data := []byte(strings.Repeat("A", 1024*1024))
	for i := 0; i < b.N; i++ {
		p.writeSerial(data)
	}
}

func BenchmarkParallel(b *testing.B) {

	p := New()
	p.Add("md5", md5.New())
	p.Add("sha1", sha1.New())
	p.Add("fnv64a", fnv.New64a())

	data := []byte(strings.Repeat("A", 1024*1024))

	for i := 0; i < b.N; i++ {
		p.Write(data)
	}
}
