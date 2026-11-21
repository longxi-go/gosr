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
	"fmt"
	"go/token"
	"strings"
	"unicode"

	"sigs.k8s.io/controller-tools/pkg/markers"
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
	repeat   map[string]bool
}

// New rejects duplicate schemas instead of relying on map overwrite order.
// The prefix starts with '+' and contains no whitespace; '+' selects all markers.
func New(prefix string, definitions []Definition) (*Registry, error) {
	if prefix == "" {
		return nil, fmt.Errorf("annotation prefix is empty")
	}
	if !strings.HasPrefix(prefix, "+") || strings.IndexFunc(prefix, unicode.IsSpace) >= 0 {
		return nil, fmt.Errorf("annotation prefix must start with '+' and contain no whitespace")
	}
	r := &Registry{
		prefix:   prefix,
		registry: &markers.Registry{},
		repeat:   make(map[string]bool),
	}
	seen := make(map[string]bool)
	for _, d := range definitions {
		if d.Sample == nil {
			return nil, fmt.Errorf("marker %s sample is nil; use struct{} for a marker without arguments", d.Name)
		}
		key := fmt.Sprintf("%d:%s", d.Target, d.Name)
		if seen[key] {
			return nil, fmt.Errorf("duplicate marker definition %s", d.Name)
		}
		seen[key] = true
		if err := r.registry.Define(d.Name, d.Target, d.Sample); err != nil {
			return nil, err
		}
		r.repeat[key] = d.Repeatable
	}
	return r, nil
}
