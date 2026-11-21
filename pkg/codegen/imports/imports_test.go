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

package imports_test

import (
	"go/types"
	"testing"

	"github.com/longxi-go/gosr/pkg/codegen/imports"
)

func TestReferencesAndConflicts(t *testing.T) {
	s := imports.New("example/local", func(path, name string, index int) string {
		return name
	})
	local := types.NewPackage("example/local", "local")
	if got := s.Type(types.NewNamed(types.NewTypeName(0, local, "Own", nil), types.Typ[types.Int], nil)); got != "Own" {
		t.Fatal(got)
	}
	s.Add("example/z", "_")
	s.Add("example/z", "z")
	s.Type(types.NewSlice(types.NewNamed(types.NewTypeName(0, types.NewPackage("example/a", "a"), "Value", nil), types.Typ[types.Int], nil)))
	list, err := s.List()
	if err != nil || len(list) != 2 || list[0].Path != "example/a" || list[1].Alias != "z" {
		t.Fatalf("%v %v", list, err)
	}
	s.Add("example/other", "a")
	if _, err := s.List(); err == nil {
		t.Fatal("alias collision accepted")
	}
}

func TestSamePathConflictAndNoImplicitImports(t *testing.T) {
	s := imports.New("", nil)
	if list, err := s.List(); err != nil || len(list) != 0 {
		t.Fatal(list, err)
	}
	s.Add("example/a", "one")
	s.Add("example/a", "two")
	if _, err := s.List(); err == nil {
		t.Fatal("same path alias conflict")
	}
}

func TestSideEffectPromotionAndDefaultQualifier(t *testing.T) {
	s := imports.New("example/app", nil)
	s.Add("example/value", "_")
	pkg := types.NewPackage("example/value", "value")
	typ := types.NewNamed(types.NewTypeName(0, pkg, "Value", nil), types.Typ[types.Int], nil)
	if got := s.Type(types.NewMap(types.Typ[types.String], types.NewPointer(typ))); got != "map[string]*value.Value" {
		t.Fatal(got)
	}
	s.Add("example/value", "_")
	list, err := s.List()
	if err != nil || len(list) != 1 || list[0].Alias != "value" {
		t.Fatal(list, err)
	}
}

func TestInvalidImportsAreReported(t *testing.T) {
	for _, entry := range []struct {
		path, alias string
	}{
		{
			"",
			"x",
		},
		{
			"example/a",
			".",
		},
		{
			"example/a",
			"for",
		},
	} {
		set := imports.New("", nil)
		set.Add(entry.path, entry.alias)
		if _, err := set.List(); err == nil {
			t.Fatalf("accepted invalid import %q as %q", entry.path, entry.alias)
		}
	}
}

func TestDefaultImportUsesDeclaredNameAndDetectsCollision(t *testing.T) {
	set := imports.New("", nil)
	set.Add("example/lib/v2", "")
	pkg := types.NewPackage("example/lib/v2", "domain")
	typ := types.NewNamed(types.NewTypeName(0, pkg, "Value", nil), types.Typ[types.Int], nil)
	if got := set.Type(typ); got != "domain.Value" {
		t.Fatal(got)
	}
	list, err := set.List()
	if err != nil || len(list) != 1 || list[0].Alias != "" {
		t.Fatal("default alias changed", list, err)
	}
	set.Add("example/other", "domain")
	if _, err := set.List(); err == nil {
		t.Fatal("declared package name collision accepted")
	}
}
