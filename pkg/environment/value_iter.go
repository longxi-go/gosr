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
	"iter"
	"maps"
	"slices"
)

// Len reports collection cardinality without copying storage. Non-collection
// values have length zero; Kind distinguishes them from empty collections.
func (v Value) Len() int {
	switch v.kind {
	case ListKind:
		return len(v.items)
	case ObjectKind:
		return len(v.fields)
	default:
		return 0
	}
}

// Elements traverses a list in index order without exposing mutable storage.
// Other kinds produce no elements. Iteration may stop early and be repeated;
// the immutable backing values also permit concurrent read-only traversal.
func (v Value) Elements() iter.Seq2[int, Value] {
	return slices.All(v.items)
}

// Members traverses an object without exposing mutable storage or copying its
// complete map. Member order is unspecified. Other kinds produce no members.
func (v Value) Members() iter.Seq2[string, Value] {
	return maps.All(v.fields)
}
