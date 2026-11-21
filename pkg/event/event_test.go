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

package event_test

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/longxi-go/gosr/pkg/event"
)

type (
	embeddedEvent struct {
		event.BaseEvent
	}

	directEvent struct{}
)

func (directEvent) EventMarker() {
}

func TestTypedNilRetainsPointerEventIdentity(t *testing.T) {
	var value event.Event = (*embeddedEvent)(nil)
	if _, ok := value.(*embeddedEvent); !ok {
		t.Fatal("pointer event identity changed")
	}
	if _, ok := value.(embeddedEvent); ok {
		t.Fatal("pointer event was implicitly dereferenced")
	}
	var direct event.Event = directEvent{}
	if _, ok := direct.(directEvent); !ok {
		t.Fatal("direct implementation requires optional embedding")
	}
}

// Contract compilation detects weakening Event to an empty interface or losing
// the Listener type constraint. It checks the real production declarations,
// rather than a copied definition or a runtime pretending to enforce typing.
func TestCompilerEnforcesEventBoundary(t *testing.T) {
	fset := token.NewFileSet()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	var files []*ast.File
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, entry.Name(), nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, file)
	}
	standard := importer.Default()
	configuration := types.Config{
		Importer: standard,
	}
	contract, err := configuration.Check("github.com/longxi-go/gosr/pkg/event", fset, files, nil)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name, fixture, diagnostic string
	}{
		{
			"external implementation",
			"consumer/external-implementation.go.tmpl",
			"",
		},
		{
			"embedded implementation",
			"consumer/embedded-implementation.go.tmpl",
			"",
		},
		{
			"reject ordinary generic argument",
			"consumer/reject-ordinary-generic-argument.go.tmpl",
			"does not satisfy",
		},
		{
			"reject ordinary publication",
			"consumer/reject-ordinary-publication.go.tmpl",
			"does not implement",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			source, err := os.ReadFile(filepath.Join("testdata", tc.fixture))
			if err != nil {
				t.Fatal(err)
			}
			file, err := parser.ParseFile(fset, tc.fixture, source, 0)
			if err != nil {
				t.Fatal(err)
			}
			configuration := types.Config{
				Importer: contractImporter{
					contract: contract,
					fallback: standard,
				},
			}
			_, err = configuration.Check("consumer", fset, []*ast.File{file}, nil)
			if tc.diagnostic == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.diagnostic) {
				t.Fatalf("want %q compile rejection, got %v", tc.diagnostic, err)
			}
		})
	}
}

type contractImporter struct {
	contract *types.Package
	fallback types.Importer
}

func (i contractImporter) Import(path string) (*types.Package, error) {
	if path == i.contract.Path() {
		return i.contract, nil
	}
	return i.fallback.Import(path)
}
