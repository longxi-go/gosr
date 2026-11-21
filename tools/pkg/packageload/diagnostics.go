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

package packageload

import (
	"errors"
	"go/token"
	"sort"

	"github.com/longxi-go/gosr/pkg/codegen/diagnostic"
	"golang.org/x/tools/go/packages"
)

// Errors visits dependencies as well as roots and reports errors deterministically
// while preserving original packages.Error values for errors.As.
func Errors(phase string, pkgs []*packages.Package) error {
	var failures []packages.Error
	packages.Visit(pkgs, func(p *packages.Package) bool {
		failures = append(failures, p.Errors...)
		return true
	}, nil)
	if len(failures) == 0 {
		return nil
	}
	sort.SliceStable(failures, func(i, j int) bool {
		left, right := failures[i], failures[j]
		if left.Error() != right.Error() {
			return left.Error() < right.Error()
		}
		if left.Kind != right.Kind {
			return left.Kind < right.Kind
		}
		return left.Pos < right.Pos
	})
	causes := make([]error, len(failures))
	for i, err := range failures {
		causes[i] = err
	}
	return diagnostic.Wrap(phase, token.Position{}, "", errors.Join(causes...))
}
