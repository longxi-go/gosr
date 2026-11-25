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
	"fmt"
	"iter"
	"slices"
	"testing"
)

func TestValueReadOnlyIteration(t *testing.T) {
	list := List(String("first"), String("2"))
	if got := list.Len(); got != 2 {
		t.Fatalf("list length = %d, want 2", got)
	}
	if got := collectElements(list.Elements()); !slices.Equal(got, []int{0, 1}) {
		t.Fatalf("list indexes = %v", got)
	}
	object := Object(map[string]Value{
		"first":  String("one"),
		"second": String("two"),
	})
	if got := object.Len(); got != 2 {
		t.Fatalf("object length = %d, want 2", got)
	}
	members := collectMembers(object.Members())
	if !slices.Equal(members, []string{"first", "second"}) && !slices.Equal(members, []string{"second", "first"}) {
		t.Fatalf("object members = %v", members)
	}
	if (Value{}).Len() != 0 || String("value").Len() != 0 {
		t.Fatal("scalar length should be zero")
	}
}

func TestValueIterationCanStopEarlyAndRepeat(t *testing.T) {
	list := List(String("first"), String("second"), String("third"))
	count := 0
	for range list.Elements() {
		count++
		break
	}
	if count != 1 {
		t.Fatalf("early-stop count = %d, want 1", count)
	}
	if got := collectElements(list.Elements()); !slices.Equal(got, []int{0, 1, 2}) {
		t.Fatalf("repeated iteration indexes = %v", got)
	}
}

func collectElements(sequence iter.Seq2[int, Value]) []int {
	values := []int{}
	for index := range sequence {
		values = append(values, index)
	}
	return values
}

func collectMembers(sequence iter.Seq2[string, Value]) []string {
	values := []string{}
	for name := range sequence {
		values = append(values, name)
	}
	return values
}

func TestValueIterationKindsAndConcurrentReaders(t *testing.T) {
	number, err := Number("12")
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range []Value{
		{},
		String("text"),
		Bool(true),
		number,
		List(),
		Object(nil),
	} {
		if value.Len() != 0 ||
			len(collectElements(value.Elements())) != 0 ||
			len(collectMembers(value.Members())) != 0 {
			t.Fatal("empty or noncollection traversal produced members")
		}
	}
	value := Object(map[string]Value{"items": List(String("one"), String("two"))})
	for range 8 {
		t.Run("reader", func(t *testing.T) {
			t.Parallel()
			for range 100 {
				count := 0
				for _, item := range value.Members() {
					if item.Len() != 2 || len(collectMembers(item.Members())) != 0 {
						t.Fatal("nested collection changed")
					}
					count++
					break
				}
				if count != 1 || len(collectElements(value.Elements())) != 0 {
					t.Fatal("object traversal changed")
				}
			}
		})
	}
}

func BenchmarkValueIterationEarlyStop(b *testing.B) {
	for _, size := range []int{10, 100000} {
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			items := make([]Value, size)
			fields := make(map[string]Value, size)
			for index := range items {
				fields[fmt.Sprint(index)] = String("value")
			}
			list, object := List(items...), Object(fields)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				for range list.Elements() {
					break
				}
				for range object.Members() {
					break
				}
			}
		})
	}
}
