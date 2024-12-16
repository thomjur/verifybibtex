// The ents.go source file includes structures and interfaces which store information for BibTeX entities
//
// Entry: struct to store information about a BibTeX entry
//
// Author: Thomas Jurczyk
// Date: December 12, 2024
package parser

import (
	"fmt"
	"log"
	"os"
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

// Debug logger
var debugLog = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

// Package vars
var regexRemoveWhiteSpace = regexp.MustCompile(`\s{2,}`)

// Regex to remove comments in BibTeX entry starting with %
// Should not remove escaped percentages like \%
var regexRemoveComments = regexp.MustCompile(`(^|[^\\])%[^\n\r]*`)

// Regex to find all valid field names
var regexFindFieldNames = regexp.MustCompile(`([a-zA-Z\s]+)=(?:\s*[{"]+)`)

// Regex to find BibTeX entry ID
var regexFindID = regexp.MustCompile(`(^|,)\s*[a-zA-Z-:_0-9]+\s*(,|$)`)

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
	// Parse fields
	newEntry.Fields, err = parseFields(cleanEntry)
	if err != nil {
		debugLog.Println(err)
	}
	// Parse ID
	newEntry.Key, err = parseID(cleanEntry)
	if err != nil {
		debugLog.Println(err)
	}
	return newEntry, nil
}

// Helper functions

// cleanRawEntry tries to clean a BibTeX raw string.
// Stripping the text of unnecessary white spaces and line breaks.
func cleanRawEntry(input string) string {
	// Trim leading and trailing white spaces
	trimmed := strings.TrimSpace(input)
	// Remove % comments
	oneLine := regexRemoveComments.ReplaceAllString(trimmed, "")
	// Remove line breaks, tabs, and carriage returns
	replacer := strings.NewReplacer("\n", "", "\r", "", "\t", "")
	oneLine = replacer.Replace(oneLine)
	// Replace multiple white spaces with single white space
	oneLine = regexRemoveWhiteSpace.ReplaceAllString(oneLine, " ")
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

// parseFields parses all fields from a clean (!) BibTeX entry.
// For cleaning a BibTeX entry, see cleanRawEntry().
func parseFields(cleanBibtexEntry string) (map[string]string, error) {
	fieldsHashMap := make(map[string]string)
	// Get the inner field first.
	// Example: @article{id, author={Thomas Jurczy},...}
	// Here, the inner field is id, author={Thomas Jurczy},...
	_, innerField, found := strings.Cut(cleanBibtexEntry, "{")
	if !found {
		return nil, &ErrParsingEntry{Message: fmt.Sprintf("Could not split on '{': %s", cleanBibtexEntry)}
	}
	// Check if innerField is empty
	innerField = strings.TrimSpace(innerField)
	if len(innerField) == 0 {
		return nil, &ErrEmptyString{Message: "The string is empty."}
	}
	// Verify trailing '}'
	if innerField[len(innerField)-1] != '}' {
		return nil, &ErrParsingEntry{Message: "The last char in fields list should be '}'."}
	}
	// Remove trailing '}'
	innerField = innerField[:len(innerField)-1]
	// Trying to find all valid fields via their field name indices
	matches := regexFindFieldNames.FindAllStringIndex(innerField, -1)
	// Storing field information in list
	// Difficult and needs better documentation
	lastIndex := 0
	previousFieldName := ""
	// Iterating over all matches
	for _, match := range matches {
		// Add previous text as value for the field
		if match[0] > lastIndex {
			if previousFieldName != "" {
				// Adding value to field
				// Clean field value
				v := innerField[lastIndex:match[0]]
				v = strings.TrimSpace(v)
				// Create []rune slice
				vrunes := []rune(v)
				// Check that v is not empty
				if len(vrunes) == 0 {
					continue
				}
				// Check if last char is ',' and remove if this is the case
				if vrunes[len(vrunes)-1] == ',' {
					vrunes = vrunes[:len(vrunes)-1]
					vrunes = []rune(strings.TrimSpace(string(vrunes)))
				}
				// Remove trailing and leading '{}' or '""'
				if (vrunes[0] == '"' && vrunes[len(vrunes)-1] == '"') || (vrunes[0] == '{' && vrunes[len(vrunes)-1] == '}') {
					vrunes = vrunes[1 : len(vrunes)-1]
				} else {
					return nil, &ErrParsingEntry{Message: fmt.Sprintf(`The first and last char in field value should either be {} or "": %s`, v)}
				}
				fieldsHashMap[previousFieldName] = string(vrunes)
			}
		}
		// Adding the field name as key to HashMap
		// Clean field name
		fieldName := innerField[match[0] : match[1]-1]
		fieldName = strings.ReplaceAll(fieldName, "=", "")
		fieldName = strings.TrimSpace(fieldName)
		fieldName = strings.ToLower(fieldName)
		if fieldName != "" {
			fieldsHashMap[fieldName] = ""
			previousFieldName = fieldName
		}
		lastIndex = match[1] - 1
	}
	// Add remaining value
	if lastIndex < len(innerField) {
		if previousFieldName != "" {
			v := innerField[lastIndex:]
			v = strings.TrimSpace(v)
			vrunes := []rune(v)
			if len(vrunes) > 0 {
				if vrunes[len(vrunes)-1] == ',' {
					vrunes = vrunes[:len(vrunes)-1]
					vrunes = []rune(strings.TrimSpace(string(vrunes)))
				}
				if (vrunes[0] == '"' && vrunes[len(vrunes)-1] == '"') || (vrunes[0] == '{' && vrunes[len(vrunes)-1] == '}') {
					vrunes = vrunes[1 : len(vrunes)-1]
				} else {
					return nil, &ErrParsingEntry{Message: fmt.Sprintf(`The first and last char in field value should either be {} or "": %s`, v)}
				}
				fieldsHashMap[previousFieldName] = string(vrunes)
			}
		}
	}
	return fieldsHashMap, nil
}

// parseID searches for a BibTeX ID in a clean (!) BibTeX entry.
// For cleaning a BibTeX entry, see cleanRawEntry().
func parseID(cleanBibtexEntry string) (string, error) {
	// Get the inner field first.
	// Example: @article{id, author={Thomas Jurczy},...}
	// Here, the inner field is id, author={Thomas Jurczy},...
	_, innerField, found := strings.Cut(cleanBibtexEntry, "{")
	if !found {
		return "", &ErrParsingEntry{Message: fmt.Sprintf("Could not split on '{': %s", cleanBibtexEntry)}
	}
	// Check if innerField is empty
	innerField = strings.TrimSpace(innerField)
	if len(innerField) == 0 {
		return "", &ErrEmptyString{Message: "The string is empty."}
	}
	// Verify trailing '}'
	if innerField[len(innerField)-1] != '}' {
		return "", &ErrParsingEntry{Message: "The last char in fields list should be '}'."}
	}
	// Remove trailing '}'
	innerField = innerField[:len(innerField)-1]
	id := regexFindID.FindString(innerField)
	// Remove potential ',' in the beginning
	if len(id) > 0 && id[0] == ',' {
		id = id[1:]
	}
	// Remove potential ',' in the end
	if len(id) > 0 && id[len(id)-1] == ',' {
		id = id[:len(id)-1]
	}
	idTrimmed := strings.TrimSpace(id)
	if idTrimmed == "" {
		return "", &ErrParsingEntry{Message: "Could not find ID in BibTeX entry."}
	}
	return idTrimmed, nil
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
