package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"io"
	"os"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/sha3"

	"github.com/mostlygeek/parhash"
	"github.com/urfave/cli"
)

func main() {

	availableHashers := map[string]func() hash.Hash{
		"md5":         md5.New,
		"sha1":        sha1.New,
		"sha256":      sha256.New,
		"sha512":      sha512.New,
		"sha3:224":    sha3.New224,
		"sha3:256":    sha3.New256,
		"sha3:384":    sha3.New384,
		"sha3:512":    sha3.New512,
		"adler32":     func() hash.Hash { return adler32.New() },
		"crc32":       func() hash.Hash { return crc32.NewIEEE() },
		"crc64:iso":   func() hash.Hash { return crc64.New(crc64.MakeTable(crc64.ISO)) },
		"crc64:ecma":  func() hash.Hash { return crc64.New(crc64.MakeTable(crc64.ECMA)) },
		"fnv64":       func() hash.Hash { return fnv.New64() },
		"fnv64a":      func() hash.Hash { return fnv.New64a() },
		"fnv32":       func() hash.Hash { return fnv.New32() },
		"fnv32a":      func() hash.Hash { return fnv.New32a() },
		"blake2b:256": func() hash.Hash { x, _ := blake2b.New256(nil); return x },
		"blake2b:384": func() hash.Hash { x, _ := blake2b.New384(nil); return x },
		"blake2b:512": func() hash.Hash { x, _ := blake2b.New512(nil); return x },
		"blake2s:256": func() hash.Hash { x, _ := blake2s.New256(nil); return x },
		// key required for blake2s:128 so excluding it
	}

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "all, a", Usage: "Use all available hash functions"},
	}
	for key, _ := range availableHashers {
		app.Flags = append(app.Flags, cli.BoolFlag{Name: key})
	}

	app.Action = func(c *cli.Context) error {
		parhasher := parhash.New()
		hashes := make(map[string]hash.Hash)

		for key, fn := range availableHashers {
			if c.Bool("all") || c.Bool(key) {
				hashes[key] = parhasher.Add(fn())
			}
		}

		filename := c.Args().Get(0)
		if filename == "" {
			return cli.NewExitError("No filename specified", 1)
		}

		if _, err := os.Stat(filename); err != nil {
			if os.IsNotExist(err) {
				return cli.NewExitError("File does not exist", 1)
			} else {
				return cli.NewExitError(err.Error(), 1)
			}
		}

		file, err := os.Open(filename)
		if err != nil {
			return cli.NewExitError("Could not open file", 1)
		}

		// hash all the data
		if _, err := io.Copy(parhasher, file); err != nil {
			return cli.NewExitError("Could not read file", 1)
		}

		for key, hash := range hashes {
			fmt.Printf("%s %s\n", key, hex.EncodeToString(hash.Sum(nil)))
		}

		return nil
	}

	app.Run(os.Args)
}
