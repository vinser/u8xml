package u8xml

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"

	"golang.org/x/text/encoding/htmlindex"
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

// https://www.ibm.com/docs/en/db2/10.5?topic=encoding-internally-encoded-xml-data

// DetectEncoding detects the encoding of a byte slice.
//
// Parameters:
// - buf: a byte slice to detect the encoding of.
//
// Returns:
// - string: the detected encoding.
// - int: the length of the BOM if a BOM is found, or 0 otherwise.
func DetectEncoding(buf []byte) (string, int) {
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

// NewReader implements io reader that converts source bytes to UTF-8.
//
// r - input io.Reader
// Returns io.Reader, error
func NewReader(r io.Reader) (io.Reader, error) {
	t := bufio.NewReader(r)
	buf, _ := t.Peek(bufCapacity)
	strEnc, bomLen := DetectEncoding(buf)
	if bomLen > 0 {
		t.Read(make([]byte, bomLen)) // skip BOM
	}
	// fmt.Printf("Encoding: %s\n", strEnc)
	if strEnc == "UTF-8" {
		return t, nil
	}
	e, err := htmlindex.Get(strEnc)
	if err != nil {
		return t, err
	}
	er := transform.NewReader(t, e.NewDecoder())
	return er, nil
}

// NewDecoder creates a new XML parser reading from r.
// Decoder converts source bytes to UTF-8
//
// r - input io.Reader
// Returns *xml.Decoder
func NewDecoder(r io.Reader) *xml.Decoder {
	u8r, _ := NewReader(r)
	d := xml.NewDecoder(u8r)
	d.CharsetReader = func(chset string, input io.Reader) (io.Reader, error) {
		return input, nil
	}
	return d
}
