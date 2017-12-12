# About

Parhash presents tools for making it easy to generate multiple
hash sums over the same data. It does this in parallel, spreading the
computation work over all available CPUs.

Parhash has a speed advantage over serial operations when working
with multiple hash algorithms over large data. Run the benchmark to see
the difference over serial operation on your machine:

```
    go test -bench .
```

On my machine, i7 3.5Ghz dual core I get a speed up of about 240%

```
  BenchmarkSerial-4            500           3577258 ns/op
  BenchmarkParallel-4         1000           1458583 ns/op
```

For smaller data sets parhash won't see significant performance
advantages due to the overhead of using channels and goroutines
distribute the work.

## Installation

```
go get github.com/mostlygeek/parhash

```

## Usage Example


```
p := parhash.New()

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
```

## CLI

A command line tool is included as an example. It includes many of the hash
algorithms available from go's standard library. Try it out, it's fun!

To install it do this:

```
cd cmd
go build -o $GOPATH/bin/parhasher parhasher.go
```

Usage:

* `parhasher -h` - prints out help
* `parhasher -a <file>` - computes sums using all available hash algorithms for <file>
* `parhasher -sha1 -md5 <file>` - computes a sha1 and md5 sum for <file>

## License

```
MIT License

Copyright (c) 2017 Benson Wong

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
