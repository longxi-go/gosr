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

package event

import "strings"

// Code classifies publication failures independently of a runtime's diagnostics.
type Code string

const (
	// InvalidEvent rejects an inadmissible payload, including a nil interface.
	InvalidEvent Code = "invalid_event"
	// InvalidState rejects publication outside an owner's admitted lifetime.
	InvalidState Code = "invalid_state"
	// ListenerFailed includes preparation, ordering, invocation and recovered
	// callback panic failures. An underlying Bean or context cause is preserved.
	ListenerFailed Code = "listener_failed"
	// PublicationCycle rejects reentry into an active declaration on one chain.
	PublicationCycle Code = "publication_cycle"
)

// Error retains declaration identity, active path and cause without depending on
// an IoC engine or its source metadata. Listener identifies a declaration, not
// an entire object. Match nonempty codes with errors.Is and inspect with errors.As;
// messages and identities remain diagnostics rather than comparison keys.
type Error struct {
	Code     Code
	Listener string
	Path     []string
	Message  string
	Cause    error
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	parts := []string{"event", string(e.Code)}
	if e.Listener != "" {
		parts = append(parts, "listener="+e.Listener)
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

// Is compares semantic classification while leaving provenance diagnostic-only.
func (e *Error) Is(target error) bool {
	other, ok := target.(*Error)
	return e != nil && ok && other != nil && e.Code != "" && e.Code == other.Code
}

// Unwrap keeps cancellation and implementation failures reachable to callers.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}
