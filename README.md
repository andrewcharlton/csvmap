[![GoDoc](https://godoc.org/github.com/andrewcharlton/csvmap?status.svg)](https://godoc.org/github.com/andrewcharlton/csvmap)
[![Build Status](https://travis-ci.org/andrewcharlton/csvmap.svg?branch=master)](https://travis-ci.org/andrewcharlton/csvmap)
[![Coverage Status](https://coveralls.io/repos/github/andrewcharlton/csvmap/badge.svg?branch=master)](https://coveralls.io/github/andrewcharlton/csvmap?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/andrewcharlton/csvmap)](https://goreportcard.com/report/github.com/andrewcharlton/csvmap)

# CSV Map

CSV Map is a wrapper for the csv package in go's standard library, designed to facilitate
easy map access to csv files with header rows.

## Installation

This package can be installed with the go get command

```
go get github.com/andrewcharlton/csvmap
```


## Documentation

API documentation can be found on [GoDoc](https://godoc.org/github.com/andrewcharlton/csvmap).
Where possible, the API has been designed to stick as closely to that of the original csv package as possible, with the exception that maps are returned instead of slices.


## Example Usage

### Reading 

``` go
func ExampleReader() {

	in := `name,alias,superpower
Logan,Wolverine,"Super healing and adamantium claws"
Charles Xavier,Professor X,Telepathy
`

	r := csvmap.NewReader(strings.NewReader(in))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Name:", record["name"])
		fmt.Println("Alias:", record["alias"])
		fmt.Println("Superpower:", record["superpower"])
		fmt.Println("")
	}

	// Output:
	// Name: Logan
	// Alias: Wolverine
	// Superpower: Super healing and adamantium claws
	//
	// Name: Charles Xavier
	// Alias: Professor X
	// Superpower: Telepathy
	//
}
```


### Writing 

``` go
func ExampleWriter() {

	headers := []string{"Name", "Alias", "Superpower"}
	data := []map[string]string{
		{"Name": "Logan", "Alias": "Wolverine", "Superpower": "Super healing"},
		{"Name": "Charles Xavier", "Alias": "Professor X", "Superpower": "Telepathy"},
	}

	out := &bytes.Buffer{}
	w := csvmap.NewWriter(out, headers)

	err := w.WriteAll(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out.String())

	// Output:
	// Name,Alias,Superpower
	// Logan,Wolverine,Super healing
	// Charles Xavier,Professor X,Telepathy
	//

}
```

