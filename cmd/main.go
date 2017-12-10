package main

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/mostlygeek/parhash"
)

func main() {

	// 1MB blob of data
	data := []byte(strings.Repeat("A", 1024*1024))

	MBs := 1024

	p := parhash.New()
	p.Add("md5", md5.New())
	p.Add("sha1", sha1.New())
	p.Add("fnv64a", fnv.New64a())
	start := time.Now()
	for i := 0; i < MBs; i++ {
		p.Write(data)
	}
	end := time.Now().Sub(start)
	fmt.Printf("Parallel Took: %v to process %d MBs of data\n", end, MBs)
}
