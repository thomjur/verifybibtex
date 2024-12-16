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

type ErrEmptyString struct {
	Message string
}

func (e *ErrParsingEntry) Error() string {
	return fmt.Sprintf("Error parsing a BibTeX entry: %s", e.Message)
}

func (e *ErrEmptyString) Error() string {
	return fmt.Sprintf("Error processing a BibTeX entry: %s", e.Message)
}

// Package vars
var regExRemoveWhiteSpace = regexp.MustCompile(`\s{2,}`)
var regExRemoveComments = regexp.MustCompile(`(^|[^\\])%[^\n\r]*`)

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
	FilePath string  // The file path of the BibTeX file.
	Name     string  // The name of the BibTeX file.
	Entries  []Entry // A slice of Entry structs representing the entries in the BibTeX file.
}

// ParseNewEntry parses a raw string in BibTeX format and tries to create an Entry struct.
// The expected format of the RawEntry string is a valid BibTeX entry, which includes the entry type,
// a unique key, and a set of fields with their corresponding values. The function cleans the raw entry
// by removing unnecessary white spaces and line breaks, and then checks if the cleaned entry is empty.
// ParseNewEntry also gracefull removes TeX comments starting with % (also using % for comments in BibTeX should generally be avoided).
// If the cleaned entry is not empty, it returns a new Entry struct with the raw entry string.
func ParseNewEntry(RawEntry string) (*Entry, error) {
	newEntry := &Entry{
		RawEntry: RawEntry,
	}
	// Clean raw entry for processing
	cleanEntry := cleanRawEntry(RawEntry)
	// Check if entry is empty
	if len(cleanEntry) == 0 {
		return nil, &ErrParsingEntry{Message: "Entry is empty after cleaning."}
	}
	newEntry.CleanEntry = cleanEntry
	// Parse entry type
	entryType, err := parseEntryType(cleanEntry)
	if err != nil {
		return nil, err
	}
	newEntry.EntryType = entryType
	return newEntry, nil

}

// Helper functions

// cleanRawEntry tries to clean a BibTeX raw string.
// Stripping the text of unnecessary white spaces and line breaks.
func cleanRawEntry(input string) string {
	// Trim leading and trailing white spaces
	trimmed := strings.TrimSpace(input)
	// Remove % comments
	oneLine := regExRemoveComments.ReplaceAllString(trimmed, "")
	// Remove line breaks, tabs, and carriage returns
	replacer := strings.NewReplacer("\n", "", "\r", "", "\t", "")
	oneLine = replacer.Replace(oneLine)
	// Replace multiple white spaces with single white space
	oneLine = regExRemoveWhiteSpace.ReplaceAllString(oneLine, " ")
	return oneLine
}

// parseEntryType parses the entry type of a BibTeX entry string.
func parseEntryType(bibtexEntry string) (string, error) {
	if len(bibtexEntry) == 0 {
		return "", &ErrEmptyString{Message: "The string is empty."}
	}
	// Split on the first appearing '{'
	substringList := strings.Split(bibtexEntry, "{")
	entryType, ok := safeGet(substringList, 0)
	if !ok {
		return "", &ErrParsingEntry{Message: fmt.Sprintf("Could not split entry on '{': %s", bibtexEntry)}
	}
	// Trim
	trimmedEntryType := strings.TrimSpace(entryType)
	if len(trimmedEntryType) == 0 {
		return "", &ErrEmptyString{Message: "The string is empty."}
	}
	// Check if type starts with an @
	if trimmedEntryType[0] != '@' {
		return "", &ErrParsingEntry{Message: fmt.Sprintf("Cannot parse entry type from this entry: %s", bibtexEntry)}
	}

	return trimmedEntryType[1:], nil
}

// safeGet retrieves the element at the specified index from the slice.
// It returns the element and a boolean indicating whether the access was successful.
func safeGet[T any](slice []T, index int) (T, bool) {
	var zeroValue T
	if index >= 0 && index < len(slice) {
		return slice[index], true
	}
	return zeroValue, false
}
