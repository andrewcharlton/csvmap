package csvmap

import (
	"encoding/csv"
	"io"
)

// A Writer writes records to a csv-encoded file.
//
// As returned by NewWriter, a Writer writes records
// terminated by a newline and uses ',' as the field delimiter.
// The exported fields can be changed to customize the details
// before the call to WriteHeaders
//
// Comma is the field delimiter.
//
// If UseCRLF is true, the Writer ends each record with \r\n instead of \n.
type Writer struct {
	Comma   rune // Field delimiter (set to ',' by NewWriter)
	UseCRLF bool // True to use \r\n as the line terminator
	// contains filtered or unexported fields

	headers   []string
	out       io.Writer
	csvWriter *csv.Writer
}

// NewWriter returns a new writer that writes to w.
//
// The file headers must be specified to provide the order in which
// record fields will be written to the writer. If headers are provided
// which don't match the fields in the records, these columns will be
// left blank.
func NewWriter(w io.Writer, headers []string) *Writer {
	return &Writer{
		Comma:   ',',
		headers: headers,
		out:     w,
	}
}

func (w *Writer) getCSVWriter() {

	w.csvWriter = csv.NewWriter(w.out)
	w.csvWriter.Comma = w.Comma
	w.csvWriter.UseCRLF = w.UseCRLF
}

func (w *Writer) writeHeaders() error {

	if w.csvWriter == nil {
		w.getCSVWriter()
	}

	return w.csvWriter.Write(w.headers)
}

// Error reports any error that has occurred during a previous Write or Flush.
func (w *Writer) Error() error {

	// If csvWriter is nil, no data has been written so ignore.
	if w.csvWriter == nil {
		return nil
	}

	return w.csvWriter.Error()
}

// Flush writes any buffered data to the underlying io.Writer.
// To check if an error occurred during the Flush, call Error.
func (w *Writer) Flush() {

	// If csvWriter is nil, no data has been written so ignore.
	if w.csvWriter == nil {
		return
	}

	w.csvWriter.Flush()
}

// Write writes a single record to w along with any necessary quoting.
//
// Write will only write fields whose keys are in the writer's headers.
// If a record has keys missing, an empty column will be written.
func (w *Writer) Write(record map[string]string) error {

	if w.csvWriter == nil {
		err := w.writeHeaders()
		if err != nil {
			return err
		}
	}

	cols := []string{}
	for _, h := range w.headers {
		cols = append(cols, record[h])
	}

	return w.csvWriter.Write(cols)
}

// WriteAll writes multiple records to w using Write and then calls flush
func (w *Writer) WriteAll(records []map[string]string) error {

	for _, record := range records {
		err := w.Write(record)
		if err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}
