// Copyright © 2026-present the longxi-go/gosr authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package annotation

import (
	"errors"
	"fmt"
	"go/token"
	"strings"
	"unicode"

	"sigs.k8s.io/controller-tools/pkg/markers"
)

const (
	markerPrefix                    = "+"
	nilMarkerSampleFormat           = "marker %s sample is nil; use struct{} for a marker without arguments"
	annotationPrefixIsEmptyMessage  = "annotation prefix is empty"
	duplicateMarkerDefinitionFormat = "duplicate marker definition %s"
	invalidAnnotationPrefixMessage  = "annotation prefix must start with '+' and contain no whitespace"
)

// Definition leaves marker schema, target adaptation and repetition policy in the
// consuming generator while reusing controller-tools validation and parsing.
// Sample must carry a type; use struct{} for a marker with no arguments.
// Numeric conversion follows controller-tools; prefer wide or string samples
// and validate domain ranges before converting to narrower integer types.
type Definition struct {
	Name       string
	Target     markers.TargetType
	Sample     any
	Repeatable bool
}

// Value retains the exact source location and full registered marker name.
type Value struct {
	Name     string
	Value    any
	Position token.Position
	Raw      string
}

// Registry is invocation-owned. Construct one for each schema/target combination;
// function doc groups may explicitly use DescribesType as controller-tools lacks
// a function target. Domain restrictions remain the consumer's responsibility.
type Registry struct {
	prefix   string
	registry *markers.Registry
	repeat   map[markerKey]bool
}

// markerKey preserves the schema's target/name identity without string encoding.
type markerKey struct {
	target markers.TargetType
	name   string
}

// New rejects duplicate schemas instead of relying on map overwrite order.
// The prefix starts with '+' and contains no whitespace; '+' selects all markers.
func New(prefix string, definitions []Definition) (*Registry, error) {
	if prefix == "" {
		return nil, errors.New(annotationPrefixIsEmptyMessage)
	}
	if !strings.HasPrefix(prefix, markerPrefix) || strings.IndexFunc(prefix, unicode.IsSpace) >= 0 {
		return nil, errors.New(invalidAnnotationPrefixMessage)
	}
	r := &Registry{
		prefix:   prefix,
		registry: &markers.Registry{},
		repeat:   make(map[markerKey]bool),
	}
	seen := make(map[markerKey]bool)
	for _, d := range definitions {
		if d.Sample == nil {
			return nil, fmt.Errorf(nilMarkerSampleFormat, d.Name)
		}
		key := markerKey{
			target: d.Target,
			name:   d.Name,
		}
		if seen[key] {
			return nil, fmt.Errorf(duplicateMarkerDefinitionFormat, d.Name)
		}
		seen[key] = true
		if err := r.registry.Define(d.Name, d.Target, d.Sample); err != nil {
			return nil, err
		}
		r.repeat[key] = d.Repeatable
	}
	return r, nil
}
