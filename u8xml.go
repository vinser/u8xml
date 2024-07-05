// Copyright 2024 Serguei Vine. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// The u8xml package implements NewDecoder which can be used to parse
// XML files with IANA character encodings such as Windows-1252, ISO-8859-1, unicode,etc.
// It can be used to decode XML files/strings with Go Standard Library xml package
// Decoder type methods like Decode(), Token(), etc.
//
// XML files must contain a BOM at the beginning in the case of unicode characters or
// an XML declaration with an encoding attribute otherwise.
//
// XML files with UTF-8 content may be detected either by BOM or XML declaration.
// XML files with no BOM or XML declaration will be treated as UTF-8.
package u8xml

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"

	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

var boms = []struct {
	bom []byte
	utf string
}{
	{[]byte{0xFF, 0xFE, 0x00, 0x00}, "UTF-32LE"},
	{[]byte{0x00, 0x00, 0xFE, 0xFF}, "UTF-32BE"},
	{[]byte{0xEF, 0xBB, 0xBF}, "UTF-8"},
	{[]byte{0xFF, 0xFE}, "UTF-16LE"},
	{[]byte{0xFE, 0xFF}, "UTF-16BE"},
}

// detectEncoding detects the encoding of a byte slice.
//
// Parameters:
// - buf: a byte slice to detect the encoding of.
//
// Returns:
// - string: the detected encoding, or default "UTF-8" if no BOM or XML declaration encoding attribute is found.
// - int: the length of the BOM if a BOM is found, or 0 otherwise.
func detectEncoding(buf []byte) (string, int) {
	// Check for a byte order mark (BOM) in the buffer.
	// If found, return the corresponding encoding and the length of the BOM.
	for _, b := range boms {
		if len(buf) < len(b.bom) {
			continue
		}
		if bytes.Equal(buf[:len(b.bom)], b.bom) {
			return b.utf, len(b.bom)
		}
	}

	// Check for an XML declaration with an encoding attribute.
	// If found, return the encoding specified in the XML declaration.
	if len(buf) < 6 || !bytes.HasPrefix(buf, []byte("<?xml")) {
		return "UTF-8", 0
	}
	encStart := bytes.Index(buf, []byte("encoding=\""))
	if encStart == -1 {
		return "UTF-8", 0
	}
	encEnd := bytes.Index(buf[encStart+11:], []byte("\""))
	if encEnd == -1 {
		return "UTF-8", 0
	}
	return string(buf[encStart+10 : encStart+encEnd+11]), 0
}

const bufCapacity = 128

// newReader implements an io reader that converts source bytes to UTF-8.
//
// r - input io.Reader
// Returns io.Reader, error
func newReader(r io.Reader) (io.Reader, error) {
	t := bufio.NewReader(r)
	buf, _ := t.Peek(bufCapacity)
	strEnc, bomLen := detectEncoding(buf)
	if bomLen > 0 {
		t.Discard(bomLen) // skip BOM
	}
	if strEnc == "UTF-8" {
		return t, nil
	}
	enc, err := ianaindex.IANA.Encoding(strEnc)
	if err != nil {
		return t, err // do not transform in the case of error
	}
	return transform.NewReader(t, enc.NewDecoder()), nil
}

// NewDecoder creates a new XML parser reading from r.
// Decoder converts source bytes to UTF-8
//
// r - input io.Reader
// Returns *xml.Decoder
func NewDecoder(r io.Reader) *xml.Decoder {
	u8r, _ := newReader(r)
	d := xml.NewDecoder(u8r)
	d.CharsetReader = func(chset string, input io.Reader) (io.Reader, error) {
		return input, nil
	}
	return d
}
