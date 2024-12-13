// Unittests for ents.go
package parser

import (
	"testing"
)

func TestCleanRawEntry(t *testing.T) {
	testCase1 := `    @article{id1234,
	
	author={Jurczyk, Thomas},
	date={20.12.2023}

	}
	
		
	`
	testCase2 := `  
             
			
	`

	// Case 1
	expected := `@article{id1234,author={Jurczyk, Thomas},date={20.12.2023}}`
	result := cleanRawEntry(testCase1)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
	// Case 2: Empty entry after preprocessing
	expected2 := ""
	result2 := cleanRawEntry(testCase2)
	if result2 != expected2 {
		t.Errorf("Expected '%s', but got '%s'", expected2, result2)
	}
}
