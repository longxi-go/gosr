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
	"strings"
)

const (
	errorDomain    = "environment"
	errorSeparator = ": "
	keyLabel       = "key="
	nilErrorText   = "<nil>"
	pathSeparator  = " -> "
	positionFormat = "%s:%d:%d"
	sourceLabel    = "source="
)

// Code classifies failures independently of an engine or diagnostic wording.
type Code string

const (
	// InvalidKey rejects an empty or whitespace-surrounded property address.
	InvalidKey Code = "invalid_key"
	// InvalidProfile rejects a profile name or predicate outside the grammar.
	InvalidProfile Code = "invalid_profile"
	// InvalidValue rejects data outside the portable value forms.
	InvalidValue Code = "invalid_value"
	// MissingProperty identifies an unresolved required reference.
	MissingProperty Code = "missing_property"
	// ResolutionCycle identifies a recursive reference chain.
	ResolutionCycle Code = "resolution_cycle"
	// ResolutionLimit identifies an engine's finite resource limit.
	ResolutionLimit Code = "resolution_limit"
	// InvalidSyntax identifies malformed placeholder text.
	InvalidSyntax Code = "invalid_syntax"
	// ConversionFailed identifies a typed binding conversion failure.
	ConversionFailed Code = "conversion_failed"
	// SourceFailed identifies failure to obtain an admitted source property.
	SourceFailed Code = "source_failed"
	// ValidationFailed identifies a rejected candidate or typed binding.
	ValidationFailed Code = "validation_failed"
	// InvalidState identifies an operation outside its admitted lifetime/state.
	InvalidState Code = "invalid_state"
)

// Error retains classification, provenance and a reference path without storing
// property values. Match nonempty codes with errors.Is and inspect details with
// errors.As. Message and identity fields must contain safe diagnostics, never
// interpolated property values. Cause may contain sensitive data: it is reachable
// through Unwrap but deliberately excluded from Error's printable form.
type Error struct {
	Code    Code
	Key     string
	Origin  Origin
	Path    []string
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e == nil {
		return nilErrorText
	}
	parts := []string{errorDomain, string(e.Code)}
	if e.Key != "" {
		parts = append(parts, keyLabel+e.Key)
	}
	if e.Origin.Source != "" {
		parts = append(parts, sourceLabel+e.Origin.Source)
	}
	if e.Origin.Resource != "" {
		parts = append(parts, fmt.Sprintf(positionFormat, e.Origin.Resource, e.Origin.Line, e.Origin.Column))
	}
	if len(e.Path) > 0 {
		parts = append(parts, strings.Join(e.Path, pathSeparator))
	}
	if e.Message != "" {
		parts = append(parts, e.Message)
	}
	return strings.Join(parts, errorSeparator)
}

// Is matches semantic classification, leaving provenance diagnostic-only.
func (e *Error) Is(target error) bool {
	other, ok := target.(*Error)
	return e != nil && ok && other != nil && e.Code != "" && e.Code == other.Code
}

// Unwrap preserves operational failures for errors.Is/errors.As without printing
// their potentially sensitive messages.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}
