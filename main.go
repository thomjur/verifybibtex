package main

import (
	"fmt"
	"os"

	"github.com/thomjur/verifybibtex/parser"
)

func main() {
	BibTeXFilePath := "bibliography.bib"
	// Trying to open the file
	file, err := os.Open(BibTeXFilePath)
	if err != nil {
		fmt.Println("Something went terribly wrong :-(")
	}
	defer file.Close()
	bibtexFile, err2 := parser.ParseNewBibTeXFile(file)
	if err2 != nil {
		fmt.Println("Something went terribly wrong :-(")
	}
	// Don't forget to add filename afterwards
	bibtexFile.FilePath = BibTeXFilePath

}
