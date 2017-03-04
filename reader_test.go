package csvmap

import (
	"reflect"
	"strings"
	"testing"
)

var readTests = []struct {
	Name   string
	Input  string
	Comma  rune
	Output []map[string]string
	Error  string
}{
	{
		Name:  "Simple",
		Input: "a,b,c\n1,2,3",
		Output: []map[string]string{
			{"a": "1", "b": "2", "c": "3"},
		},
	},
	{
		Name:  "Mutliple rows",
		Input: "A,B,C,D\n1,2,3,4\n5,6,7,8",
		Output: []map[string]string{
			{"A": "1", "B": "2", "C": "3", "D": "4"},
			{"A": "5", "B": "6", "C": "7", "D": "8"},
		},
	},
	{
		Name:   "Empty File",
		Input:  "",
		Output: []map[string]string{},
		Error:  "",
	},
	{
		Name:  "Long row",
		Input: "A,B,C\n1,2,3,4",
		Error: "wrong number of fields in line",
	},
	{
		Name:  "Short row",
		Input: "A,B,C\n1,2",
		Error: "wrong number of fields in line",
	},
	{
		Name:  "Duplicate Headers",
		Input: "A,B,C,A\n1,2,3,4",
		Error: "duplicate headers found",
	},
	{
		Name:  "| Delimiter",
		Input: "A|B|C\n1|2|3",
		Comma: '|',
		Output: []map[string]string{
			{"A": "1", "B": "2", "C": "3"},
		},
	},
}

func TestRead(t *testing.T) {

	for _, tt := range readTests {
		r := NewReader(strings.NewReader(tt.Input))
		if tt.Comma != 0 {
			r.Comma = tt.Comma
		}
		out, err := r.ReadAll()

		if !reflect.DeepEqual(out, tt.Output) {
			t.Errorf("%v: out=%q want=%q", tt.Name, out, tt.Output)
		}

		if tt.Error != "" {
			if err == nil || !strings.Contains(err.Error(), tt.Error) {
				t.Errorf("%s: error %v, want error %q", tt.Name, err, tt.Error)
			}
		}
	}
}

func TestReadHeaders(t *testing.T) {

	input := strings.NewReader("A,B,C\n1,2,3")
	r := NewReader(input)

	headers, err := r.ReadHeaders()
	exp := []string{"A", "B", "C"}
	if !reflect.DeepEqual(headers, exp) {
		t.Errorf("out=%q, want=%q", headers, exp)
	}

	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}

	headers, err = r.ReadHeaders()
	if headers != nil {
		t.Errorf("out=%q, expected=nil", headers)
	}
	if err != ErrHeadersSet {
		t.Errorf("err %q, want error %q", err, ErrHeadersSet)
	}

}
