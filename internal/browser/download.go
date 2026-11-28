package browser

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"path"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// A download into the grant's store (browser batch 1): a file the page
// gave, through the session's egress (the grant's sites only), saved by the
// browser and never opened by it, then typed by its content, not its name,
// checked and size-capped before planque or the taf sees it. Only PDFs, CSV
// text, PNG and JPEG images are taken; the taf reads them as data, never
// instructions. The check is structural (what could run or open is
// refused), not an antivirus engine.

// MaxDownload is the largest file taken; a larger one is stopped as it
// arrives and dropped.
const MaxDownload = 5 << 20

// Download is one file the page gave. Refused says why it wasn't taken
// (Data is then empty); Type is its type as read from its bytes.
type Download struct {
	Name    string `json:"name"`
	Type    string `json:"type,omitempty"`
	Size    int    `json:"size"`
	SHA256  string `json:"sha256,omitempty"`
	Data    []byte `json:"data,omitempty"`
	Refused string `json:"refused,omitempty"`
	// Formulas: a CSV with a cell that a spreadsheet would run as a
	// formula (it starts with =, +, -, @, or a tab or carriage return), for
	// the apps to say so before the owner opens it.
	Formulas bool `json:"formulas,omitempty"`
}

// Downloader is a session that hands over the files its pages gave since
// the last ask, each taken or refused, and forgets them.
type Downloader interface {
	Downloads(ctx context.Context) ([]Download, error)
}

// The types a download may be.
const (
	TypePDF  = "application/pdf"
	TypeCSV  = "text/csv"
	TypePNG  = "image/png"
	TypeJPEG = "image/jpeg"
)

var (
	// ErrDownloadType: the bytes are neither a PDF nor CSV text.
	ErrDownloadType = errors.New("only a PDF, CSV, PNG or JPEG file is taken")
	// ErrDownloadActive: a PDF that carries something that runs or opens
	// (a script, a launch, an embedded file, a form's submit), or that
	// hides its objects from the check (compressed object streams,
	// encryption).
	ErrDownloadActive = errors.New("the PDF carries a script, a launch action, an embedded file or a form submit, or hides its objects (compressed object streams or encryption), and isn't taken")
	// ErrDownloadPixels: an image larger than MaxPixels, refused before
	// it's decoded (a decompression bomb declares a huge size).
	ErrDownloadPixels = errors.New("the image is larger than 40 million pixels")
)

// MaxPixels is the largest image taken, by its declared size, checked
// before anything is decoded.
const MaxPixels = 40_000_000

// pdfActive are a PDF's names for content that runs or opens something.
// ObjStm and Encrypt are refused too: inside a compressed object stream, or
// an encrypted file, a dictionary isn't readable by this check.
var pdfActive = regexp.MustCompile(`/(JavaScript|JS|Launch|EmbeddedFile|RichMedia|SubmitForm|ImportData|GoToE|XFA|ObjStm|Encrypt)\b`)

// pdfEscape is a #xx escape in a PDF name ("/Java#53cript" is /JavaScript).
var pdfEscape = regexp.MustCompile(`#([0-9A-Fa-f]{2})`)

// pdfNames is b with every name's #xx escapes decoded, so an escaped name
// reads as the one it stands for.
func pdfNames(b []byte) []byte {
	return pdfEscape.ReplaceAllFunc(b, func(m []byte) []byte {
		v, _ := strconv.ParseUint(string(m[1:]), 16, 8)
		return []byte{byte(v)}
	})
}

// TypeOf is a download's type, read from its bytes: a PDF (its header,
// and none of the names that run or open something), a PNG or JPEG that
// decodes as one, or CSV text (UTF-8, no control bytes but tab and line
// ends, and a header row split by commas or semicolons). Anything else is
// refused.
func TypeOf(b []byte) (string, error) {
	switch {
	case len(b) == 0:
		return "", ErrDownloadType
	case bytes.HasPrefix(b, []byte("%PDF-")):
		if pdfActive.Match(b) || pdfActive.Match(pdfNames(b)) {
			return "", ErrDownloadActive
		}
		return TypePDF, nil
	case bytes.HasPrefix(b, []byte("\x89PNG\r\n\x1a\n")):
		if err := imageFits(b, png.DecodeConfig); err != nil {
			return "", err
		}
		if _, err := png.Decode(bytes.NewReader(b)); err != nil {
			return "", ErrDownloadType
		}
		return TypePNG, nil
	case bytes.HasPrefix(b, []byte{0xff, 0xd8, 0xff}):
		if err := imageFits(b, jpeg.DecodeConfig); err != nil {
			return "", err
		}
		if _, err := jpeg.Decode(bytes.NewReader(b)); err != nil {
			return "", ErrDownloadType
		}
		return TypeJPEG, nil
	case csvText(b):
		return TypeCSV, nil
	}
	return "", ErrDownloadType
}

// imageFits reads an image's declared size, and refuses one that doesn't
// read, or is larger than MaxPixels, before it's decoded.
func imageFits(b []byte, config func(io.Reader) (image.Config, error)) error {
	c, err := config(bytes.NewReader(b))
	if err != nil {
		return ErrDownloadType
	}
	if c.Width <= 0 || c.Height <= 0 || int64(c.Width)*int64(c.Height) > MaxPixels {
		return ErrDownloadPixels
	}
	return nil
}

// CSVFormulas says a CSV has a cell a spreadsheet would run as a formula.
func CSVFormulas(b []byte) bool {
	for _, line := range strings.Split(string(b), "\n") {
		for _, cell := range strings.FieldsFunc(line, func(r rune) bool { return r == ',' || r == ';' }) {
			cell = strings.TrimLeft(cell, " \"")
			if cell != "" && strings.ContainsRune("=+-@\t\r", rune(cell[0])) {
				return true
			}
		}
	}
	return false
}

func csvText(b []byte) bool {
	b = bytes.TrimPrefix(b, []byte("\xef\xbb\xbf")) // a byte-order mark
	if !utf8.Valid(b) {
		return false
	}
	for _, c := range b {
		if c < 0x20 && c != '\t' && c != '\n' && c != '\r' {
			return false
		}
	}
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	sep := ","
	if strings.Count(lines[0], ";") > strings.Count(lines[0], ",") {
		sep = ";"
	}
	if !strings.Contains(lines[0], sep) {
		return false
	}
	return true
}

var nameUnsafe = regexp.MustCompile(`[^A-Za-z0-9._ -]+`)

// DownloadName is a file's name as kept: its base only, with nothing but
// letters, digits, dots, dashes, underscores and spaces, at most 100
// characters, and the extension its type has.
func DownloadName(suggested, typ string) string {
	base := path.Base(strings.ReplaceAll(suggested, "\\", "/"))
	base = strings.Trim(nameUnsafe.ReplaceAllString(base, "_"), " ._")
	ext := map[string]string{TypePDF: ".pdf", TypeCSV: ".csv", TypePNG: ".png", TypeJPEG: ".jpg"}[typ]
	for _, e := range []string{".pdf", ".csv", ".png", ".jpg", ".jpeg"} {
		base = strings.TrimSuffix(base, e)
	}
	if base == "" {
		base = "download"
	}
	if len(base) > 96 {
		base = base[:96]
	}
	return base + ext
}
