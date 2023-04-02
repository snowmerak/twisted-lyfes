package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/snowmerak/twisted-lyfes/src/db/fsys"
	"golang.org/x/crypto/blake2b"
	"log"
	"os"
)

func main() {
	home := os.Getenv("HOME")
	fs := fsys.New(home+"/.twisted-lyfes", func(bytes []byte) []byte {
		bs := blake2b.Sum512(bytes)
		return bs[:]
	}, hex.EncodeToString)

	importFlag := flag.String("import", "", "import a file into lyfes")
	exportFlag := flag.String("export", "", "export a file to home path")

	flag.Parse()

	if *importFlag != "" {
		if err := fs.Import(*importFlag); err != nil {
			panic(err)
		}
		log.Println("imported", *importFlag)
		return
	}

	if *exportFlag != "" {
		if err := fs.Export(home, *exportFlag); err != nil {
			panic(err)
		}
		log.Println("exported", *exportFlag)
		return
	}

	fmt.Println("Usage: fsys -import <file> or fsys -export <file>")
}
