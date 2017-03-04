package csvmap

import (
	"reflect"
	"strings"
	"testing"
)

var readTests = []struct {
	Name   string
	Input  string
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
}

func TestRead(t *testing.T) {

	for _, tt := range readTests {
		r := NewReader(strings.NewReader(tt.Input))
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
