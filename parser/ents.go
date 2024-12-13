// The ents.go source file includes structures and interfaces which store information for BibTeX entities
//
// Entry: struct to store information about a BibTeX entry
//
// Author: Thomas Jurczyk
// Date: December 12, 2024
package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// Define errors
type ErrParsingEntry struct {
	Message string
}

func (e *ErrParsingEntry) Error() string {
	return fmt.Sprintf("Error parsing a BibTeX entry: %s", e.Message)
}

// Package vars
var regExRemoveWhiteSpace = regexp.MustCompile(`\s{2,}`)

// Entry represents a bibliographic entry in a BibTeX file.
// It contains the type of the entry (e.g., article, book),
// a unique key to identify the entry, the raw entry string,
// and a map of fields with their corresponding values.
type Entry struct {
	EntryType string            // The type of the entry (e.g., article, book).
	Key       string            // A unique key to identify the entry.
	RawEntry  string            // The raw entry string in BibTeX format.
	Fields    map[string]string // A map of fields and their corresponding values.
}

// BibTeXFile represents a BibTeX file with its associated metadata.
// It contains the file path, name, and a list of entries.
type BibTeXFile struct {
	FilePath string  // The file path of the BibTeX file.
	Name     string  // The name of the BibTeX file.
	Entries  []Entry // A slice of Entry structs representing the entries in the BibTeX file.
}

// CreateNewEntry parses a raw string in BibTeX format and tries to create an Entry struct.
// The expected format of the RawEntry string is a valid BibTeX entry, which includes the entry type,
// a unique key, and a set of fields with their corresponding values. The function cleans the raw entry
// by removing unnecessary white spaces and line breaks, and then checks if the cleaned entry is empty.
// If the cleaned entry is not empty, it returns a new Entry struct with the raw entry string.
func CreateNewEntry(RawEntry string) (*Entry, error) {
	newEntry := &Entry{
		RawEntry: RawEntry,
	}
	// Clean raw entry for processing
	cleanEntry := cleanRawEntry(RawEntry)
	// Check if entry is empty
	if len(cleanEntry) == 0 {
		return nil, &ErrParsingEntry{Message: "Entry is empty after cleaning."}
	}
	return newEntry, nil

}

// Helper functions

// cleanRawEntry tries to clean a BibTeX raw string.
// Stripping the text of unnecessary white spaces and line breaks
func cleanRawEntry(input string) string {
	// Trim leading and trailing white spaces
	trimmed := strings.TrimSpace(input)
	// Remove line breaks, tabs, and carriage returns
	replacer := strings.NewReplacer("\n", "", "\r", "", "\t", "")
	oneLine := replacer.Replace(trimmed)
	// Replace multiple white spaces with single white space
	oneLine = regExRemoveWhiteSpace.ReplaceAllString(oneLine, " ")
	return oneLine
}
