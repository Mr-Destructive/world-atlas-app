package game

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

type PlaceInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Dictionary struct {
	places map[string]PlaceInfo
	mu     sync.RWMutex
}

func NewDictionary(filepath string) (*Dictionary, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var placeList []PlaceInfo
	// Try parsing as array of objects
	if err := json.Unmarshal(data, &placeList); err != nil {
		// Fallback for simple string array during migration
		var simpleList []string
		if err2 := json.Unmarshal(data, &simpleList); err2 == nil {
			for _, s := range simpleList {
				placeList = append(placeList, PlaceInfo{Name: s, Type: "Place"})
			}
		} else {
			return nil, err
		}
	}

	places := make(map[string]PlaceInfo)
	for _, p := range placeList {
		places[strings.ToLower(p.Name)] = p
	}

	return &Dictionary{
		places: places,
	}, nil
}

func (d *Dictionary) IsValid(place string) (bool, string, string) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	info, ok := d.places[strings.ToLower(place)]
	if ok {
		return true, info.Type, info.Name
	}
	return false, "", ""
}

func (d *Dictionary) GetInfo(place string) PlaceInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.places[strings.ToLower(place)]
}

func (d *Dictionary) GetUnusedPlaceStartingWith(letter byte, used map[string]bool) PlaceInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()
	target := strings.ToLower(string(letter))
	for lowerName, info := range d.places {
		if strings.HasPrefix(lowerName, target) {
			if !used[lowerName] {
				return info
			}
		}
	}
	return PlaceInfo{}
}