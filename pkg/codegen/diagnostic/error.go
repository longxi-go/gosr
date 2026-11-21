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

package diagnostic

import (
	"go/token"
	"strings"
)

// Error preserves a machine-readable phase, source position and original cause
// without requiring consumers to classify failures by matching display strings.
type Error struct {
	Phase    string
	Position token.Position
	Message  string
	Cause    error
}

// Wrap supplies boundary context without changing errors.Is/As behavior.
func Wrap(phase string, position token.Position, message string, cause error) error {
	return &Error{
		Phase:    phase,
		Position: position,
		Message:  message,
		Cause:    cause,
	}
}

func (e *Error) Error() string {
	var parts []string
	if e.Position.IsValid() {
		parts = append(parts, e.Position.String())
	}
	if e.Phase != "" {
		parts = append(parts, e.Phase)
	}
	if e.Message != "" {
		parts = append(parts, e.Message)
	}
	if e.Cause != nil {
		parts = append(parts, e.Cause.Error())
	}
	return strings.Join(parts, ": ")
}

// Unwrap retains underlying parser, template and filesystem error identities.
func (e *Error) Unwrap() error {
	return e.Cause
}
