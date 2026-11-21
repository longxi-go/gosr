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

package imports

import (
	"fmt"
	"go/token"
	"go/types"
	"sort"
)

// AliasFunc chooses an import alias only when a package is first used. An empty
// alias selects the package's declared name; a blank import cannot qualify a type.
// The index includes explicit imports; callers can preserve established naming.
type AliasFunc = func(importPath, packageName string, index int) string

// Import retains explicit aliases, including side-effect imports.
type Import struct {
	Path  string
	Alias string
}

type binding struct {
	alias      string
	identifier string
}

// Set owns the import namespace of one generated file. Errors are retained so
// go/types qualifier callbacks cannot silently discard a collision.
// It is invocation-owned and must not be mutated concurrently.
type Set struct {
	local   string
	alias   AliasFunc
	entries map[string]binding
	err     error
}

// New excludes local-package references and defers registration until actual use.
func New(local string, alias AliasFunc) *Set {
	return &Set{
		local:   local,
		alias:   alias,
		entries: make(map[string]binding),
	}
}

// Add registers an explicit alias and returns the known identifier. The result
// remains empty for a default import until Type establishes its declared name.
// A side-effect import may be promoted to a named use; incompatible named aliases
// are errors. Paths alone cannot identify package names.
func (s *Set) Add(importPath, alias string) string {
	if importPath != "" && importPath == s.local {
		return ""
	}
	if s.err != nil {
		return alias
	}
	if importPath == "" || alias == "." || alias == "init" ||
		(alias != "" && alias != "_" &&
			(!token.IsIdentifier(alias) || token.Lookup(alias).IsKeyword())) {
		s.err = fmt.Errorf("invalid import %q as %q", importPath, alias)
		return alias
	}
	if old, ok := s.entries[importPath]; ok {
		if old.alias == alias || alias == "_" {
			return old.identifier
		}
		if old.alias != "_" {
			s.err = fmt.Errorf("import %s has conflicting aliases %q and %q", importPath, old.alias, alias)
			return alias
		}
	}
	s.entries[importPath] = binding{
		alias:      alias,
		identifier: alias,
	}
	s.checkIdentifier(importPath, alias)
	return alias
}

// Type spells a type using go/types, registering only packages in the emitted
// expression. Diagnostic-only type spelling should use types.TypeString directly.
func (s *Set) Type(t types.Type) string {
	return types.TypeString(t, s.qualifier)
}

func (s *Set) qualifier(pkg *types.Package) string {
	if pkg.Path() == s.local {
		return ""
	}
	if s.err != nil {
		return pkg.Name()
	}
	if entry, ok := s.entries[pkg.Path()]; ok && entry.alias != "_" {
		if entry.alias == "" {
			s.checkIdentifier(pkg.Path(), pkg.Name())
			entry.identifier = pkg.Name()
			s.entries[pkg.Path()] = entry
			return pkg.Name()
		}
		return entry.alias
	}
	alias := pkg.Name()
	if s.alias != nil {
		alias = s.alias(pkg.Path(), pkg.Name(), len(s.entries))
	}
	if alias == "_" {
		s.err = fmt.Errorf("blank import %s cannot qualify a type", pkg.Path())
		return pkg.Name()
	}
	if alias == "" {
		s.Add(pkg.Path(), "")
		if s.err == nil {
			entry := s.entries[pkg.Path()]
			entry.identifier = pkg.Name()
			s.entries[pkg.Path()] = entry
			s.checkIdentifier(pkg.Path(), pkg.Name())
		}
		return pkg.Name()
	}
	return s.Add(pkg.Path(), alias)
}

func (s *Set) checkIdentifier(importPath, identifier string) {
	if s.err != nil || identifier == "" || identifier == "_" {
		return
	}
	for path, entry := range s.entries {
		if path != importPath && entry.identifier == identifier {
			s.err = fmt.Errorf("import alias %s conflicts between %s and %s", identifier, path, importPath)
			return
		}
	}
}

// List returns a deterministic copy, or the first namespace error.
func (s *Set) List() ([]Import, error) {
	if s.err != nil {
		return nil, s.err
	}
	result := make([]Import, 0, len(s.entries))
	paths := make([]string, 0, len(s.entries))
	for p := range s.entries {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		entry := s.entries[p]
		if entry.alias == "" && entry.identifier == "" {
			return nil, fmt.Errorf("default import %s requires a typed reference or explicit alias", p)
		}
		result = append(result, Import{
			Path:  p,
			Alias: entry.alias,
		})
	}
	return result, nil
}
