// Package csvmap wraps the standard library's encoding/csv package to provide
// reading and writing maps to csv files.
//
// Because this package only wraps encoding/csv, it only supports the csv file
// format specified in RFC 4180. Please see the documentation for encoding/csv
// for more details.
//
// This package assumes that the first row of the csv file contains header data
// which provides the keys for the following map access. Reading:
//
//     Header1,Header2,Header3
//     Field1,Field2,Field3
//
// results in:
//
//     {"Header1":"Field1", "Header2":"Field2", "Header3":"Field3"}
//
// For files with lines of additional header information, the Discard function is
// provided to remove these before reading the header row.
//
// A Writer method is also provided for writing mapped data to file as well.
//
//		w := NewWriter(..., []string{"Header1", "Header", "Header3"})
//		w.Write({"Header1":"Field1", "Header2":"Field2", "Header3":"Field3"})
//		w.Flush()
//
// results in:
//     Header1,Header2,Header3
//     Field1,Field2,Field3
//
package csvmap
