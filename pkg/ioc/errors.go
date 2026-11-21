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

package ioc

import (
	"fmt"
	"strings"
)

// Code classifies failures independently of diagnostic wording or an engine.
type Code string

const (
	InvalidDefinition   Code = "invalid_definition"
	InvalidSelector     Code = "invalid_selector"
	DuplicateDefinition Code = "duplicate_definition"
	MissingBean         Code = "missing_bean"
	AmbiguousBean       Code = "ambiguous_bean"
	ContractMismatch    Code = "contract_mismatch"
	DependencyCycle     Code = "dependency_cycle"
	ScopeViolation      Code = "scope_violation"
	InvalidState        Code = "invalid_state"
	UnsupportedContract Code = "unsupported_contract"
	CreationFailed      Code = "creation_failed"
)

// Error preserves semantic classification, provenance and the underlying cause.
// Match a nonempty Code with errors.Is; names and messages are diagnostics only.
type Error struct {
	Code          Code
	Bean, Message string
	Source        Source
	Path          []string
	Cause         error
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	parts := []string{"ioc", string(e.Code)}
	if e.Bean != "" {
		parts = append(parts, "bean="+e.Bean)
	}
	if e.Source.File != "" {
		parts = append(parts, fmt.Sprintf("%s:%d", e.Source.File, e.Source.Line))
	}
	if len(e.Path) > 0 {
		parts = append(parts, strings.Join(e.Path, " -> "))
	}
	if e.Message != "" {
		parts = append(parts, e.Message)
	}
	if e.Cause != nil {
		parts = append(parts, e.Cause.Error())
	}
	return strings.Join(parts, ": ")
}

// Is compares semantic codes without depending on a particular engine's text.
func (e *Error) Is(target error) bool {
	other, ok := target.(*Error)
	return e != nil && ok && other != nil && e.Code != "" && e.Code == other.Code
}

// Unwrap retains cancellation and original operational failures for callers.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}
