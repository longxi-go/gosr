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

package environment

import (
	"encoding/json"
	"maps"
	"slices"
	"strings"
)

// Kind identifies one of the closed portable value forms. A missing property is
// represented by Lookup's bool, never by an additional value kind.
type Kind uint8

const (
	// NullKind is the zero value, distinct from both missing and empty string.
	NullKind Kind = iota
	// StringKind retains text for strict resolution by a resolver.
	StringKind
	// BoolKind retains a boolean without imposing typed binding rules.
	BoolKind
	// NumberKind preserves validated JSON-number text without rounding.
	NumberKind
	// ListKind contains an ordered sequence of immutable values.
	ListKind
	// ObjectKind contains exact member names and immutable values.
	ObjectKind
)

// Value is a closed immutable property value safe to share across snapshots and
// goroutines. Its zero value is null. Collections are detached on construction
// and access; nested Values may share storage because no mutation is exposed.
// Number equality compares retained text, so 1 and 1.0 remain distinct.
type Value struct {
	kind   Kind
	scalar string
	items  []Value
	fields map[string]Value
}

// String preserves text verbatim, including empty strings and placeholders.
func String(text string) Value {
	return Value{
		kind:   StringKind,
		scalar: text,
	}
}

// Bool retains a boolean with its canonical scalar spelling.
func Bool(value bool) Value {
	text := "false"
	if value {
		text = "true"
	}
	return Value{
		kind:   BoolKind,
		scalar: text,
	}
}

// Number accepts exactly one JSON number without surrounding whitespace. It
// preserves the spelling and precision, including exponents outside float64's
// range. Invalid text returns null and InvalidValue without echoing the input.
func Number(text string) (Value, error) {
	if text == "" ||
		(text[0] != '-' && (text[0] < '0' || text[0] > '9')) ||
		strings.TrimSpace(text) != text || !json.Valid([]byte(text)) {
		return Value{}, &Error{
			Code:    InvalidValue,
			Message: "expected JSON number text",
		}
	}
	return Value{
		kind:   NumberKind,
		scalar: text,
	}, nil
}

// List detaches the caller's sequence. Nil and empty sequences are both empty
// lists; neither represents null.
func List(items ...Value) Value {
	return Value{
		kind:  ListKind,
		items: slices.Clone(items),
	}
}

// Object detaches the caller's map. Member names are data, not property
// addresses, so ValidateKey is not applied to them. Nil is an empty object.
func Object(fields map[string]Value) Value {
	return Value{
		kind:   ObjectKind,
		fields: maps.Clone(fields),
	}
}

// Kind exposes the value form without opening its representation to mutation.
func (v Value) Kind() Kind {
	return v.kind
}

// Scalar returns retained string/number text or canonical boolean text. Null,
// lists and objects have no scalar representation and return false.
func (v Value) Scalar() (string, bool) {
	switch v.kind {
	case StringKind, BoolKind, NumberKind:
		return v.scalar, true
	default:
		return "", false
	}
}

// Items returns a detached list or nil for other kinds. Inspect Kind to
// distinguish a nil-backed empty list from null.
func (v Value) Items() []Value {
	return slices.Clone(v.items)
}

// Fields returns a detached object or nil for other kinds. Inspect Kind to
// distinguish a nil-backed empty object from null.
func (v Value) Fields() map[string]Value {
	return maps.Clone(v.fields)
}

// Equal compares kind and contents recursively, ignoring object iteration order.
// It preserves scalar spelling and list order; empty nil-backed collections are
// equal to empty allocated collections of the same kind.
func (v Value) Equal(other Value) bool {
	if v.kind != other.kind || v.scalar != other.scalar {
		return false
	}
	if !slices.EqualFunc(v.items, other.items, Value.Equal) {
		return false
	}
	return maps.EqualFunc(v.fields, other.fields, Value.Equal)
}
