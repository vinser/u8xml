package u8xml

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectEncoding(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		expect string
		bomLen int
	}{
		{"UTF-8 with BOM", []byte{0xEF, 0xBB, 0xBF, 't', 'e', 's', 't'}, "UTF-8", 3},
		{"UTF-16BE with BOM", []byte{0xFE, 0xFF, 0, 't', 0, 'e', 0, 's', 0, 't'}, "UTF-16BE", 2},
		{"UTF-16LE with BOM", []byte{0xFF, 0xFE, 't', 0, 'e', 0, 's', 0, 't', 0}, "UTF-16LE", 2},
		{"UTF-32BE with BOM", []byte{0, 0, 0xFE, 0xFF, 0, 0, 't', 0, 0, 0, 'e', 0, 0, 's', 0, 0, 't', 0, 0, 0}, "UTF-32BE", 4},
		{"UTF-32LE with BOM", []byte{0xFF, 0xFE, 0, 0, 't', 0, 0, 0, 'e', 0, 0, 's', 0, 0, 't', 0, 0, 0, 0, 0}, "UTF-32LE", 4},
		{"UTF-8 without BOM", []byte{'t', 'e', 's', 't'}, "UTF-8", 0},
		{"XML declaration with \"ISO-8859-1\" encoding attribute", []byte("<?xml encoding=\"ISO-8859-1\"?>"), "ISO-8859-1", 0},
		{"XML declaration with empty encoding attribute", []byte("<?xml?>"), "UTF-8", 0},
		{"XML declaration with unclosed encoding attribute", []byte("<?xml encoding=\"ISO-8859-1"), "UTF-8", 0},
		{"Too small buffer", []byte("A"), "UTF-8", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, bomLen := DetectEncoding(tt.input)
			assert.Equal(t, tt.expect, enc)
			assert.Equal(t, tt.bomLen, bomLen)
		})
	}
}

var errEncodingNotSupported = errors.New("htmlindex: invalid encoding name")

func TestNewReader(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
		err    error
	}{
		{"UTF-8 with BOM", "\xEF\xBB\xBFtest", "test", nil},
		{"UTF-8 without BOM and without XMP declaration", "test", "test", nil},
		{"UTF-16LE", "\xFF\xFE\x74\x00\x65\x00\x73\x00\x74\x00", "test", nil},
		{"Windows-1251", "<?xml encoding=\"Windows-1251\"?>\xC1\xF3\xEB\xE3\xE0\xEA\xEE\xE2", "<?xml encoding=\"Windows-1251\"?>Булгаков", nil},
		{"Unsupported encoding", "<?xml encoding=\"Windows-1\"?>\xC1\xF3\xEB\xE3\xE0\xEA\xEE\xE2", "<?xml encoding=\"Windows-1251\"?>Булгаков", errEncodingNotSupported},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader([]byte(tt.input))
			reader, err := NewReader(r)
			if tt.err != nil {
				assert.Equal(t, tt.err, err)
				return
			}
			assert.Nil(t, err)
			b, err := io.ReadAll(reader)
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, string(b))
		})
	}
}

type person struct {
	Name string `xml:"Name"`
	Age  int    `xml:"Age"`
}

func TestNewDecoder(t *testing.T) {
	tests := []struct {
		name   string
		file   string
		expect string
	}{
		{"iso-8859-1", "test-samples/iso-8859-1.xml", "Gabriel García Márquez ISO-8859-1"},
		{"iso-8859-2", "test-samples/iso-8859-2.xml", "Ľudovít Štúr ISO-8859-2"},
		{"windows-1251", "test-samples/windows-1251.xml", "Михаил Афанасьевич Булгаков Windows-1251"},
		{"utf-16", "test-samples/utf-16.xml", "Ľudovít Štúr utf-16"},
		{"utf-16le.", "test-samples/utf-16le.xml", "Ľudovít Štúr utf-16le"},
		{"utf-16be", "test-samples/utf-16be.xml", "Ľudovít Štúr utf-16be"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.ReadFile(tt.file)
			assert.Nil(t, err)
			d := NewDecoder(bytes.NewReader(f))
			var p person
			err = d.Decode(&p)
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, p.Name)
		})
	}
}
