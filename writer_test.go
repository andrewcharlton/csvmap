package csvmap

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

var writerTests = []struct {
	Name    string
	Headers []string
	In      []map[string]string
	Comma   rune
	Out     string
}{
	{
		Name:    "Simple",
		Headers: []string{"A", "B", "C"},
		In: []map[string]string{
			{"A": "1", "B": "2", "C": "3"},
			{"A": "x", "B": "y", "C": "z"},
		},
		Out: "A,B,C\n1,2,3\nx,y,z\n",
	},
	{
		Name:    "Missing headers",
		Headers: []string{"A", "B"},
		In: []map[string]string{
			{"A": "1", "B": "2", "C": "3"},
			{"A": "x", "B": "y", "C": "z"},
		},
		Out: "A,B\n1,2\nx,y\n",
	},
	{
		Name:    "Too many headers",
		Headers: []string{"A", "B", "C", "D"},
		In: []map[string]string{
			{"A": "1", "B": "2", "C": "3"},
			{"A": "x", "B": "y", "C": "z"},
		},
		Out: "A,B,C,D\n1,2,3,\nx,y,z,\n",
	},
}

func TestWriter(t *testing.T) {

	for _, tt := range writerTests {

		buf := &bytes.Buffer{}
		w := NewWriter(buf, tt.Headers)
		if tt.Comma != 0 {
			w.Comma = tt.Comma
		}

		err := w.WriteAll(tt.In)

		if err != nil {
			t.Errorf("%v. Unexpected error: %v", tt.Name, err)
		}

		out := buf.String()
		if out != tt.Out {
			t.Errorf("%v. Out=%q Want=%q", tt.Name, out, tt.Out)
		}

	}
}

func TestUnwrittenWriter(t *testing.T) {

	buf := &bytes.Buffer{}
	w := NewWriter(buf, []string{})
	w.Flush()
	err := w.Error()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestBrokenWriter(t *testing.T) {

	b := BrokenWriter{}
	w := NewWriter(b, []string{"A", "B"})
	record := []map[string]string{{"A": "1", "B": "2"}}

	err := w.WriteAll(record)
	if err == nil || !strings.Contains(err.Error(), "Error writing") {
		t.Errorf("Got %q, Want %q", err, "Error writing")
	}
}

type BrokenWriter struct{}

func (w BrokenWriter) Write(p []byte) (int, error) {

	buf := &bytes.Buffer{}
	buf.Write(p)

	return 0, errors.New("Error writing")

}
