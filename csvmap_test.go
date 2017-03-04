package csvmap_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/andrewcharlton/csvmap"
)

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
