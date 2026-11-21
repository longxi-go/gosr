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

package ordered_test

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/longxi-go/gosr/pkg/ordered"
)

func Example() {
	// This extension consumer owns its stable identities and samples before
	// sorting. The standard does not impose a container or sorting utility.
	objects := []struct {
		id    string
		value any
	}{
		{
			"ordinary",
			ordinary(ordered.HighestPrecedence),
		},
		{
			"unordered",
			struct{}{},
		},
		{
			"priority",
			priority{
				ordinary(ordered.LowestPrecedence),
			},
		},
	}
	type key struct {
		id           string
		group, order int
	}
	keys := make([]key, 0, len(objects))
	for _, object := range objects {
		k := key{
			id:    object.id,
			group: 1,
			order: ordered.LowestPrecedence,
		}
		if _, ok := object.value.(ordered.PriorityOrdered); ok {
			k.group = 0
		}
		if value, ok := object.value.(ordered.Ordered); ok {
			k.order = value.Order()
		}
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a, b key) int {
		return cmp.Or(cmp.Compare(a.group, b.group), cmp.Compare(a.order, b.order), cmp.Compare(a.id, b.id))
	})
	for _, k := range keys {
		fmt.Println(k.id)
	}
	// Output:
	// priority
	// ordinary
	// unordered
}
