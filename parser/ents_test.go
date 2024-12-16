// Unittests for ents.go
package parser

import (
	"fmt"
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
