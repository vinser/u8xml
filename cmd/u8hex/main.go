// With u8hex CLI utility they can get the hex representation of a string with a given character set.
// It may be useful for debugging
package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"golang.org/x/text/encoding/ianaindex"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run", path.Base(os.Args[0]), "<input_string> <charset>")
		return
	}

	inputString := os.Args[1]
	charset := os.Args[2]
	encoding, err := ianaindex.IANA.Encoding(charset)
	if err != nil {
		log.Fatalln(err)
	}
	if encoding == nil {
		log.Fatalln("ianaindex: unsupported charset")
	}

	outputBytes := make([]byte, 1024)
	transformer := encoding.NewEncoder()
	n, _, err := transformer.Transform(outputBytes, []byte(inputString), true)
	if err != nil {
		log.Fatalln(err)
	}

	hexCodes := make([]string, 0)
	for _, b := range outputBytes[:n] {
		hexCodes = append(hexCodes, fmt.Sprintf("\\x%X", b))
	}
	hexSequence := strings.Join(hexCodes, "")

	fmt.Println(hexSequence)
}
