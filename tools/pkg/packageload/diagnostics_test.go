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

package packageload_test

import (
	"errors"
	"testing"

	"github.com/longxi-go/gosr/tools/pkg/packageload"
	"golang.org/x/tools/go/packages"
)

func TestStructuredDiagnosticOrderIsIndependentOfRootOrder(t *testing.T) {
	parse := &packages.Package{
		ID: "parse",
		Errors: []packages.Error{
			{
				Pos:  "a.go:1",
				Msg:  "broken",
				Kind: packages.ParseError,
			},
		},
	}
	typed := &packages.Package{
		ID: "types",
		Errors: []packages.Error{
			{
				Pos:  "a.go:1",
				Msg:  "broken",
				Kind: packages.TypeError,
			},
		},
	}
	for _, roots := range [][]*packages.Package{
		{
			parse,
			typed,
		},
		{
			typed,
			parse,
		},
	} {
		var first packages.Error
		err := packageload.Errors("check", roots)
		if !errors.As(err, &first) || first.Kind != packages.ParseError {
			t.Fatalf("first diagnostic = %#v, want parse error; %v", first, err)
		}
	}
}
