// Unittests for ents.go
package parser

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestCleanRawEntry(t *testing.T) {
	// Case 1: Valid raw string with many empty spaces
	testCase1 := `    @article{id1234,
	
	author={Jurczyk, Thomas},
	date={20.12.2023}

	}
	
		
	`
	expected := `@article{id1234,author={Jurczyk, Thomas},date={20.12.2023}}`
	result := cleanRawEntry(testCase1)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
	// Case 2: Empty entry after preprocessing
	testCase2 := `  
             
			
	`
	expected2 := ""
	result2 := cleanRawEntry(testCase2)
	if result2 != expected2 {
		t.Errorf("Expected '%s', but got '%s'", expected2, result2)
	}
	// Case 3: Removing comments
	testCase3 := `% This is a comment
	@article{id1234, % yet another comment!
	% yet another comment
	title={Drugs \% Comments}, % Remove this!!!!
	author={Jurczyk, Thomas},
	date={20.12.2023}

	}
	
		
	`
	result3 := cleanRawEntry(testCase3)
	expected3 := `@article{id1234,title={Drugs \% Comments},author={Jurczyk, Thomas},date={20.12.2023}}`
	if result3 != expected3 {
		t.Errorf("Expected '%s', but got '%s'", expected3, result3)
	}
}

func TestParseEntryType(t *testing.T) {
	// Case 1: Valid raw entry
	testCase1 := `    @article{id1234,
	
	author={Jurczyk, Thomas},
	date={20.12.2023}

	}
	
		
	`
	expected := `article`
	result, _ := parseEntryType(testCase1)
	if expected != result {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
	// Case 2: Valid clean string
	testCase2 := `@article{id1234,author={Jurczyk, Thomas},date={20.12.2023}}`
	expected2 := `article`
	result2, _ := parseEntryType(testCase2)
	if expected2 != result2 {
		t.Errorf("Expected '%s', but got '%s'", expected2, result2)
	}
	// Case 3: Only white spaces string
	expected3 := &ErrEmptyString{Message: "The string is empty."}
	_, err := parseEntryType(`    `)
	if err == nil || expected3.Error() != err.Error() {
		t.Errorf("Expected '%#v', but got '%#v'", expected3, err)
	}
	// Case 4: Empty string
	expected4 := &ErrEmptyString{Message: "The string is empty."}
	_, err2 := parseEntryType("")
	if err2 == nil || expected4.Error() != err2.Error() {
		t.Errorf("Expected '%#v', but got '%#v'", expected4, err2)
	}
	// Case 5: Invalid BibTeX entry
	expected5 := &ErrParsingEntry{Message: fmt.Sprintf("Cannot parse entry type from this entry: %s", "Hey guys. This is absolutely no valid BibTeX entry, but just a normal text!!!!")}
	_, err3 := parseEntryType("Hey guys. This is absolutely no valid BibTeX entry, but just a normal text!!!!")
	if err3 == nil || expected5.Error() != err3.Error() {
		t.Errorf("Expected '%#v', but got '%#v'", expected5, err3)
	}
	// Case 6: Invalid BibTeX entry
	expected6 := &ErrParsingEntry{Message: fmt.Sprintf("Cannot parse entry type from this entry: %s", "article{id1234,author={Jurczyk, Thomas},date={20.12.2023}}")}
	_, err4 := parseEntryType("article{id1234,author={Jurczyk, Thomas},date={20.12.2023}}")
	if err4 == nil || expected6.Error() != err4.Error() {
		t.Errorf("Expected '%#v', but got '%#v'", expected6, err4)
	}
}

func TestParseFields(t *testing.T) {
	entry1 := `{schmidt2024,
  author       = {Schmidt, Anna and Müller, Bernd and {O'Connor}, Claire and García, Diego},
  editor       = "Weber, Eva and {D'Amico}, Fabio",
  title        = {Fortgeschrittene Datenanalyse (mit, = in dem Text) mit Python: Methoden und Anwendungen}, publisher    = {Technik Verlag}  ,
  year         = {2024},
  volume       = {3},
  SERIES       = {Datenwissenschaftliche Studien},
  address      = {München},
  edition      = {2., überarbeitete und erweiterte Auflage},
  month        = {März},
  isbn         = {978-3-16-148410-0},
  doi          = {10.1000/182},
  url          = {https://www.technik-verlag.de/buecher/fortgeschrittene-datenanalyse},
  note         = {Beinhaltet ein Kapitel über maschinelles Lernen},
  abstract     = {Dieses Buch bietet eine umfassende Einführung in fortgeschrittene Methoden der Datenanalyse mit Python, einschließlich praxisnaher Anwendungen und Fallstudien.},
  keywords     = {Datenanalyse, Python, maschinelles Lernen, Statistik},
  language     = {Deutsch}
}`

	// Case 1: Complex entry, just checking if all fields have been found
	expected1 := []string{"publisher", "author", "editor", "year", "title", "volume", "series", "address", "edition", "month", "isbn", "doi", "url", "note", "abstract", "keywords", "language"}
	// Sort list for comparison
	sort.Strings(expected1)

	fields, err := parseFields(entry1)
	// Collect field names
	fieldNameList := make([]string, 0, 8)
	for k := range fields {
		fieldNameList = append(fieldNameList, k)
	}
	// Sort fieldnameList for comparison
	sort.Strings(fieldNameList)

	if !reflect.DeepEqual(expected1, fieldNameList) {
		t.Errorf("Expected '%#v', but got '%#v'", expected1, fieldNameList)
		t.Errorf("Error: %s", err.Error())
	}

	// Case 2: Simple bibliography
	entry2 := `@book{schmidt2024,author = {Schmidt, Anna and Müller, Bernd and {O'Connor}, Claire and García, Diego},language = "Deutsch"}`
	expected2 := map[string]string{"author": "Schmidt, Anna and Müller, Bernd and {O'Connor}, Claire and García, Diego", "language": "Deutsch"}

	fields2, _ := parseFields(entry2)

	if !reflect.DeepEqual(expected2, fields2) {
		t.Errorf("Expected '%#v', but got '%#v'", expected2, fields2)
	}

	// Case 3: Valid simple entry
	entry3 := `@article{muster2024,
  author  = {Max Mustermann},
  title   = {Einführung in die Datenwissenschaft},
  journal = {Journal für Informatik},
  year    = {2024},
  volume  = {42},
  number  = {3},
  pages   = {123--145}
}
`
	expected3 := map[string]string{"author": "Max Mustermann", "title": "Einführung in die Datenwissenschaft", "journal": "Journal für Informatik", "year": "2024", "volume": "42", "number": "3", "pages": "123--145"}

	fields3, _ := parseFields(entry3)

	if !reflect.DeepEqual(expected3, fields3) {
		t.Errorf("Expected '%#v', but got '%#v'", expected3, fields3)
	}

}

func TestParseID(t *testing.T) {
	// Case 1: Valid ID as it should be
	entry1 := `@book{muster2024,
	author  = {Max Mustermann},
	title   = {Einführung in die Datenwissenschaft},
	journal = {Journal für Informatik},
	year    = {2024},
	volume  = {42},
	number  = {3},
	pages   = {123--145}
}
  `
	expected1 := "muster2024"

	parsedIDString, err := parseID(entry1)

	if err != nil {
		debugLog.Println(err.Error())
	}

	if expected1 != parsedIDString {
		t.Errorf("Expected '%#v', but got '%#v'", expected1, parsedIDString)
	}

	// Case 2: ID not in first position
	entry2 := `@book{
	author  = {Max Mustermann},
	muster2024,
	title   = {Einführung in die Datenwissenschaft},
	journal = {Journal für Informatik},
	year    = {2024},
	volume  = {42},
	number  = {3},
	pages   = {123--145}
 } 
  `
	expected2 := "muster2024"

	parsedIDString2, err2 := parseID(entry2)

	if err != nil {
		debugLog.Println(err2.Error())
	}

	if expected2 != parsedIDString2 {
		t.Errorf("Expected '%#v', but got '%#v'", expected2, parsedIDString2)
	}

	// Case 3: ID last position
	entry3 := `@book{
			author  = {Max Mustermann},
			title   = {Einführung in die Datenwissenschaft},
			journal = {Journal für Informatik},
			year    = {2024},
			volume  = {42},
			number  = {3},
			pages   = {123--145},
			muster2024
		 } 
		  `
	expected3 := "muster2024"

	parsedIDString3, err3 := parseID(entry3)

	if err != nil {
		debugLog.Println(err3.Error())
	}

	if expected3 != parsedIDString3 {
		t.Errorf("Expected '%#v', but got '%#v'", expected3, parsedIDString3)
	}

}
