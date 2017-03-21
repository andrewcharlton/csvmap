package csvmap

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
)

var (
	// ErrDuplicateHeaders is returned when there are duplicated items in the
	// header row.
	ErrDuplicateHeaders = errors.New("duplicate headers found")
)

// A Reader returns records (a map of values) from a csv-encoded file.
//
// As returned by NewReader, a Reader expects input conforming to RFC 4180.
// The exported fields can be changed to customize the details before the
// first call to Headers/Read/ReadAll.
//
// The header row will be read on the first call to Headers/Read/ReadAll.
// If there are duplicated keys in the header, an ErrDuplicateHeaders error
// will be returned at this point.
type Reader struct {
	// Comma is the field delimiter.
	// It is set to comma (',') by NewReader.
	Comma rune
	// Comment, if not 0, is the comment character. Lines beginning with the
	// Comment character without preceding whitespace are ignored.
	// With leading whitespace the Comment character becomes part of the
	// field, even if TrimLeadingSpace is true.
	Comment rune
	// If LazyQuotes is true, a quote may appear in an unquoted field and a
	// non-doubled quote may appear in a quoted field.
	LazyQuotes bool
	// If TrimLeadingSpace is true, leading white space in a field is ignored.
	// This is done even if the field delimiter, Comma, is white space.
	TrimLeadingSpace bool

	// List of headers
	headers []string

	// input reader
	in io.Reader

	// csv reader
	csvReader *csv.Reader
}

// NewReader returns a reader that will read from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		Comma:     ',',
		in:        r,
		headers:   nil,
		csvReader: nil,
	}
}

// Discard ignores the first n lines of the input reader before
// reading the headers. If should be called before the first call
// to Headers/Read/ReadAll, otherwise it will return an ErrHeaderSet
// error.
//
// If there are insufficient lines to discard, it will return an
// io.EOF error.
func (r *Reader) Discard(n int) error {

	if r.csvReader != nil {
		return ErrHeaderSet
	}

	buf := bufio.NewReader(r.in)

	for i := 0; i < n; i++ {
		_, err := buf.ReadString('\n')
		if err != nil {
			return err
		}
	}

	r.in = buf
	return nil
}

// getReader creates the underlying csv.Reader prior to being used.
func (r *Reader) getReader() {

	r.csvReader = csv.NewReader(r.in)
	r.csvReader.Comma = r.Comma
	r.csvReader.Comment = r.Comment
	r.csvReader.FieldsPerRecord = 0
	r.csvReader.LazyQuotes = r.LazyQuotes
	r.csvReader.TrimLeadingSpace = r.TrimLeadingSpace

}

// readHeaders reads the first line of the file and sets the headers.
// If there are duplicated headers in the file, it will return  ErrDuplicatedHeaders.
//
// readHeaders should only be called once. If it is called again, it will
// return nil, ErrHeadersSet.
func (r *Reader) readHeaders() error {

	if r.csvReader == nil {
		r.getReader()
	}

	if r.headers != nil {
		return ErrHeaderSet
	}

	headers, err := r.csvReader.Read()
	if err != nil {
		return err
	}

	check := map[string]struct{}{}
	for _, h := range headers {
		_, exists := check[h]
		if exists {
			return ErrDuplicateHeaders
		}
		check[h] = struct{}{}
	}

	r.headers = headers
	return nil
}

// Headers returns the column headers
func (r *Reader) Headers() ([]string, error) {

	if r.headers == nil {
		err := r.readHeaders()
		if err != nil {
			return nil, err
		}
	}

	headers := append([]string{}, r.headers...)
	return headers, nil
}

// Read reads one record (a map of fields to values). If the record
// has the unexpected number of fields, Read returns a map of the values
// present, along with a csv.ErrFieldCount error. Except for that case,
// Read always returns either a non-nil record or a non-nil error, but not both.
// If there is no data left to be read, Read returns nil, io.EOF.
//
// On the first call of Read/ReadAll, if headers have not been set by
// ReadHeaders, this will be done automatically.
func (r *Reader) Read() (map[string]string, error) {

	if r.headers == nil {
		err := r.readHeaders()
		if err != nil {
			return nil, err
		}
	}

	fields, err := r.csvReader.Read()
	if err != nil && (err == io.EOF || len(fields) == len(r.headers)) {
		return nil, err
	}

	if len(fields) > len(r.headers) {
		fields = fields[:len(r.headers)]
	}

	record := map[string]string{}
	for n, f := range fields {
		record[r.headers[n]] = f
	}

	return record, err
}

// ReadAll reads all the remaining records from r. Each record is a slice of
// fields. A successful call returns err == nil, not err == io.EOF. Because
// ReadAll is defined to read until EOF, it does not treat end of file as an error to be reported.
func (r *Reader) ReadAll() ([]map[string]string, error) {

	records := []map[string]string{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			return records, nil
		}
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}
}
