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

package environment_test

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/longxi-go/gosr/pkg/environment"
)

func TestValueKindsAndScalars(t *testing.T) {
	number, err := environment.Number("-123456789012345678901.50e+200")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		value  environment.Value
		kind   environment.Kind
		text   string
		scalar bool
	}{
		{
			environment.Value{},
			environment.NullKind,
			"",
			false,
		},
		{
			environment.String(""),
			environment.StringKind,
			"",
			true,
		},
		{
			environment.String("hello"),
			environment.StringKind,
			"hello",
			true,
		},
		{
			environment.Bool(true),
			environment.BoolKind,
			"true",
			true,
		},
		{
			environment.Bool(false),
			environment.BoolKind,
			"false",
			true,
		},
		{
			number,
			environment.NumberKind,
			"-123456789012345678901.50e+200",
			true,
		},
		{
			environment.List(),
			environment.ListKind,
			"",
			false,
		},
		{
			environment.Object(nil),
			environment.ObjectKind,
			"",
			false,
		},
	}
	for _, tc := range cases {
		if tc.value.Kind() != tc.kind {
			t.Errorf("kind = %v, want %v", tc.value.Kind(), tc.kind)
		}
		if text, scalar := tc.value.Scalar(); text != tc.text || scalar != tc.scalar {
			t.Errorf("scalar = %q, %v; want %q, %v", text, scalar, tc.text, tc.scalar)
		}
		if tc.kind != environment.ListKind && tc.value.Items() != nil {
			t.Error("non-list exposes items")
		}
		if tc.kind != environment.ObjectKind && tc.value.Fields() != nil {
			t.Error("non-object exposes fields")
		}
	}
}

func TestNumberGrammar(t *testing.T) {
	for _, text := range []string{
		"0", "-0", "10", "-10", "0.01", "1e0", "1E-9", "1e+999999",
		"123456789012345678901234567890",
	} {
		v, err := environment.Number(text)
		if err != nil {
			t.Errorf("valid number %q: %v", text, err)
			continue
		}
		if got, ok := v.Scalar(); !ok || got != text {
			t.Errorf("number changed: %q", got)
		}
	}
	for _, text := range []string{
		"", "+1", "01", "-01", " 1", "1 ", "1\n", ".1", "1.", "-", "1e", "1e+",
		"NaN", "Infinity", "null", "true", "[1]", "1 2", "1_000", "١", "0x10",
		"secret-token",
	} {
		v, err := environment.Number(text)
		if !errors.Is(err, &environment.Error{
			Code: environment.InvalidValue,
		}) {
			t.Errorf("invalid number %q: %v", text, err)
		}
		if !v.Equal(environment.Value{}) {
			t.Errorf("invalid number %q returned non-null", text)
		}
		if err != nil && text == "secret-token" && strings.Contains(err.Error(), text) {
			t.Error("number error prints input")
		}
	}
}

func TestCollectionIsolation(t *testing.T) {
	original := environment.String("original")
	fields := map[string]environment.Value{"secret": original}
	items := []environment.Value{environment.Object(fields)}
	v := environment.List(items...)
	fields["secret"] = environment.String("mutated")
	items[0] = environment.Value{}
	first := v.Items()
	first[0].Fields()["secret"] = environment.String("changed")
	first[0] = environment.Value{}
	got := v.Items()[0].Fields()["secret"]
	if !got.Equal(original) {
		t.Fatal("constructor or accessor aliases caller state")
	}
	object := environment.Object(map[string]environment.Value{"nested": v})
	exposed := object.Fields()
	delete(exposed, "nested")
	if !object.Fields()["nested"].Equal(v) {
		t.Fatal("object accessor aliases state")
	}
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 50 {
				detached := object.Fields()
				detached["nested"] = environment.Value{}
				_ = object.Equal(object)
			}
		})
	}
	wg.Wait()
}

func TestValueEquality(t *testing.T) {
	one, err := environment.Number("1")
	if err != nil {
		t.Fatal(err)
	}
	decimal, err := environment.Number("1.0")
	if err != nil {
		t.Fatal(err)
	}
	equal := [][2]environment.Value{
		{
			{},
			{},
		},
		{
			environment.String(""),
			environment.String(""),
		},
		{
			environment.List(),
			environment.List(nil...),
		},
		{
			environment.Object(nil),
			environment.Object(map[string]environment.Value{}),
		},
		{
			environment.Object(map[string]environment.Value{
				"a": one,
				"b": environment.List(environment.Bool(true)),
			}),
			environment.Object(map[string]environment.Value{
				"b": environment.List(environment.Bool(true)),
				"a": one,
			}),
		},
	}
	for _, pair := range equal {
		if !pair[0].Equal(pair[1]) || !pair[1].Equal(pair[0]) {
			t.Error("equal values differ")
		}
	}
	unequal := [][2]environment.Value{
		{
			{},
			environment.String(""),
		},
		{
			one,
			environment.String("1"),
		},
		{
			one,
			decimal,
		},
		{
			environment.Bool(false),
			environment.Bool(true),
		},
		{
			environment.String("a"),
			environment.String("b"),
		},
		{
			environment.List(),
			environment.List(environment.Value{}),
		},
		{
			environment.List(one, decimal),
			environment.List(decimal, one),
		},
		{
			environment.Object(nil),
			environment.Object(map[string]environment.Value{
				"x": {},
			}),
		},
		{
			environment.Object(map[string]environment.Value{
				"x": {},
			}),
			environment.Object(map[string]environment.Value{
				"y": {},
			}),
		},
		{
			environment.Object(map[string]environment.Value{"x": one}),
			environment.Object(map[string]environment.Value{"x": decimal}),
		},
	}
	for _, pair := range unequal {
		if pair[0].Equal(pair[1]) || pair[1].Equal(pair[0]) {
			t.Error("different values compare equal")
		}
	}
}
