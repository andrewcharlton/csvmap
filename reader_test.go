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
		Name:  "Multiple rows",
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

	headers, err := r.Headers()
	exp := []string{"A", "B", "C"}
	if !reflect.DeepEqual(headers, exp) {
		t.Errorf("out=%q, want=%q", headers, exp)
	}

	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}
}

func TestDuplicateHeaders(t *testing.T) {

	input := strings.NewReader("A,B,C,A\n1,2,3,4")
	r := NewReader(input)

	headers, err := r.Headers()
	if headers != nil {
		t.Errorf("Unexpected headers: %v", headers)
	}

	if err != ErrDuplicateHeaders {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestDiscard(t *testing.T) {

	in := "\n\n\n\nA,B,C\n1,2,3\n"

	testcases := []struct {
		Name    string
		N       int
		Err     string
		Headers []string
		HeadErr string
	}{
		{
			Name:    "Simple",
			N:       4,
			Err:     "",
			Headers: []string{"A", "B", "C"},
			HeadErr: "",
		},
		{
			Name:    "EOF",
			N:       10,
			Err:     "EOF",
			Headers: nil,
			HeadErr: "EOF",
		},
	}

	for _, tc := range testcases {

		r := NewReader(strings.NewReader(in))

		err := r.Discard(tc.N)
		if err == nil && tc.Err != "" {
			t.Errorf("%v. Missing error. Got %q, Want %q", tc.Name, err, tc.Err)
		}
		if err != nil && (tc.Err == "" || !strings.Contains(err.Error(), tc.Err)) {
			t.Errorf("%v. Got %q, want error %q", tc.Name, err, tc.Err)
		}

		headers, err := r.Headers()
		if !reflect.DeepEqual(headers, tc.Headers) {
			t.Errorf("%v. Wrong headers. Got %q, Want %q", tc.Name, headers, tc.Headers)
		}

		if err == nil && tc.HeadErr != "" {
			t.Errorf("%v. Missing error. Got %q, Want %q", tc.Name, err, tc.HeadErr)
		}
		if err != nil && (tc.HeadErr == "" || !strings.Contains(err.Error(), tc.HeadErr)) {
			t.Errorf("%v. Got %q, want error %q", tc.Name, err, tc.HeadErr)
		}

	}
}

func TestDiscardErr(t *testing.T) {

	r := NewReader(strings.NewReader("A,B,C\n1,2,3"))
	r.Headers()

	err := r.Discard(2)
	if err != ErrHeaderSet {
		t.Errorf("Unexpected error: Got %q, want %q", err, ErrHeaderSet)
	}
}
