package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/LindsayBradford/go-dbf/godbf"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func toASCII(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	
	// Further cleanup: remove any remaining non-ASCII characters and filter to just letters, spaces, hyphens
	var b strings.Builder
	for _, r := range result {
		if r <= unicode.MaxASCII && (unicode.IsLetter(r) || r == ' ' || r == '-') {
			b.WriteRune(r)
		}
	}
	return strings.TrimSpace(b.String())
}

type Place struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func main() {
	base := "/home/meet/code/Natural_Earth_quick_start/packages/Natural_Earth_quick_start"
	outputFile := "../data/places.json"

	// Map Name -> Type (Country overrides State, State overrides City if duplicate names exist, or keep all? 
	// For simplicity, let's keep the "highest" level (Country > State > City)
	placeMap := make(map[string]string)

	sources := []struct {
		SubPath string
		File    string
		Col     string
		Type    string
	}{
		// 10m - High resolution
		{"10m_cultural", "ne_10m_populated_places.dbf", "NAME", "City"},
		{"10m_cultural", "ne_10m_admin_1_states_provinces.dbf", "name", "State"}, 
		{"10m_cultural", "ne_10m_admin_0_countries.dbf", "NAME", "Country"},
	}

	for _, src := range sources {
		fullPath := filepath.Join(base, src.SubPath, src.File)
		
		err := readDBF(fullPath, src.Col, src.Type, placeMap)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Skipping %s: File not found.\n", src.Type)
			} else {
				fmt.Printf("Error processing %s: %v\n", src.Type, err)
			}
		} else {
			fmt.Printf("Processed %s\n", src.Type)
		}
	}

	// Convert to List
	var finalList []Place
	for name, pType := range placeMap {
		if len(name) > 1 {
			finalList = append(finalList, Place{Name: name, Type: pType})
		}
	}
	
	// Sort by name
	sort.Slice(finalList, func(i, j int) bool {
		return finalList[i].Name < finalList[j].Name
	})

	data, err := json.MarshalIndent(finalList, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Successfully wrote %d places to %s\n", len(finalList), outputFile)
}

func readDBF(path string, colName string, typeName string, storage map[string]string) error {
	dbfTable, err := godbf.NewFromFile(path, "UTF8")
	if err != nil {
		return err
	}

	colIdx := -1
	fields := dbfTable.Fields()
	for j, field := range fields {
		if strings.EqualFold(field.Name(), colName) {
			colIdx = j
			break
		}
	}

	if colIdx == -1 {
		return fmt.Errorf("column '%s' not found", colName)
	}

	for i := 0; i < dbfTable.NumberOfRecords(); i++ {
		row := dbfTable.GetRowAsSlice(i)
		if colIdx < len(row) {
			val := strings.TrimSpace(row[colIdx])
			val = clean(val)
			if val != "" && isValid(val) {
				// Clean key for case-insensitive storage? 
				// No, let's store Display Name in Map, but check lower
				// Actually the game logic lowercases everything. 
				// Let's just store the Display Name.
				
				// Priority: Country > State > City
				// If already exists...
				currentType, exists := storage[val]
				if exists {
					if typeName == "Country" {
						storage[val] = typeName
					} else if typeName == "State" && currentType == "City" {
						storage[val] = typeName
					}
				} else {
					storage[val] = typeName
				}
			}
		}
	}
	return nil
}

func clean(s string) string {
	if idx := strings.Index(s, "("); idx != -1 {
		s = s[:idx]
	}
	return toASCII(s)
}

func isValid(s string) bool {
	if len(s) < 2 { return false }
	if strings.ContainsAny(s, "0123456789") { return false }
	return true
}
