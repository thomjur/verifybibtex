# VerifyBibTeX
VerifyBibTeX-Go is CLI tool to verify the correctness of a BibTeX file. This includes parsing errors as well as potentially missing fields. This version of the tool has been written in Golang. There is also a Docker container that uses Python: [VerifyBibTeX-OS](https://github.com/phimisci/verifybibtex-os)

Currently, this repository only includes the `parser` library that can be used
independently of the later to build verifier. It can be imported via
`github.com/thomjur/verifybibtex/parser`.

The main usage is:

```go
// Relative path to BibTeX file
BibTeXFilePath := "bibliography.bib"

// Trying to open the file
file, err := os.Open(BibTeXFilePath)
if err != nil {
    fmt.Println("Something went terribly wrong :-(")
}
defer file.Close()

//Parse BibTeX file
bibtexFile, err := parser.ParseNewBibTeXFile(file)
if err != nil {
    fmt.Println("Something went terribly wrong :-(")
}

```

The `bibTeXFile` is a list of `Entry` structs.

```go
// Entry represents a bibliographic entry in a BibTeX file.
// It contains the type of the entry (e.g., article, book),
// a unique key to identify the entry, the raw entry string,
// and a map of fields with their corresponding values.
type Entry struct {
	EntryType  string            // The type of the entry (e.g., article, book).
	Key        string            // A unique key to identify the entry.
	RawEntry   string            // The raw entry string in BibTeX format.
	CleanEntry string            // The cleaned raw BibTeX input (RawEntry).
	Fields     map[string]string // A map of fields and their corresponding values.
}

// BibTeXFile represents a BibTeX file with its associated metadata.
// It contains the file path, name, and a list of entries.
type BibTeXFile struct {
	FilePath string   // The file path of the BibTeX file.
	Entries  []*Entry // A slice of Entry structs representing the entries in the BibTeX file.
}

```

## Version
2025-05-19
