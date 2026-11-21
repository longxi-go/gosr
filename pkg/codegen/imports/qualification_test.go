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
	"strings"
	"testing"

	"github.com/longxi-go/gosr/pkg/codegen/imports"
)

func TestPreferredDefaultAliasKeepsTypedReferencesStable(t *testing.T) {
	set := imports.New("example/app", defaultAlias)
	pkg := types.NewPackage("example/lib/v2", "domain")
	typ := types.NewNamed(types.NewTypeName(0, pkg, "Value", nil), types.Typ[types.Int], nil)
	for range 2 {
		if got := set.Type(typ); got != "domain.Value" {
			t.Fatalf("typed reference = %q, want domain.Value", got)
		}
	}
	if got := set.Add(pkg.Path(), ""); got != "domain" {
		t.Fatalf("registered identifier = %q, want domain", got)
	}
	list, err := set.List()
	if err != nil || len(list) != 1 || list[0].Alias != "" {
		t.Fatalf("default import changed: %v, %v", list, err)
	}
	set.Add("example/other", "domain")
	if _, err := set.List(); err == nil {
		t.Fatal("default identifier collision accepted")
	}
}

func TestPreferredSideEffectAliasCannotQualifyAType(t *testing.T) {
	set := imports.New("example/app", sideEffectAlias)
	pkg := types.NewPackage("example/lib", "domain")
	typ := types.NewNamed(types.NewTypeName(0, pkg, "Value", nil), types.Typ[types.Int], nil)
	set.Type(typ)
	if _, err := set.List(); err == nil {
		t.Fatal("side-effect import accepted for a typed reference")
	}
}

func TestInitCannotBeAnImportIdentifier(t *testing.T) {
	set := imports.New("example/app", nil)
	set.Add("context", "init")
	if _, err := set.List(); err == nil {
		t.Fatal("reserved init import identifier accepted")
	}
}

func TestUnresolvedDefaultImportDiagnosticsUsePathOrder(t *testing.T) {
	set := imports.New("example/app", nil)
	set.Add("example/z", "")
	set.Add("example/a", "")
	for range 256 {
		_, err := set.List()
		if err == nil || !strings.Contains(err.Error(), "example/a") {
			t.Fatalf("first unresolved import = %v, want example/a", err)
		}
	}
}

func defaultAlias(string, string, int) string {
	return ""
}

func sideEffectAlias(string, string, int) string {
	return "_"
}
