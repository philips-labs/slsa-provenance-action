package intoto

import (
	"encoding/json"
	"fmt"
	"io"
)

// WithMaterials adds additional materials to the predicate
func WithMaterials(materials []Item) StatementOption {
	return func(s *Statement) {
		s.Predicate.Materials = append(s.Predicate.Materials, materials...)
	}
}

// ReadMaterials reads the material from file
func ReadMaterials(r io.Reader) ([]Item, error) {
	var materials []Item

	if err := json.NewDecoder(r).Decode(&materials); err != nil {
		return nil, err
	}

	for _, material := range materials {
		if material.URI == "" {
			return nil, fmt.Errorf("empty or missing \"uri\" for material")
		}
		if len(material.Digest) == 0 {
			return nil, fmt.Errorf("empty or missing \"digest\" for material")
		}
	}

	return materials, nil
}
