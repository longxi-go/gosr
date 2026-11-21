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
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/longxi-go/gosr/pkg/event"
)

func TestFailureRetainsClassificationCauseAndDeclaration(t *testing.T) {
	cause := context.DeadlineExceeded
	failure := &event.Error{
		Code:     event.ListenerFailed,
		Listener: "orders/audit.OnEvent",
		Message:  "invoke",
		Cause:    cause,
	}
	err := fmt.Errorf("publish order: %w", failure)
	if !errors.Is(err, cause) || !errors.Is(err, &event.Error{
		Code: event.ListenerFailed,
	}) {
		t.Fatal("wrapping lost cause or event classification")
	}
	for _, target := range []error{&event.Error{
		Code: event.InvalidEvent,
	}, &event.Error{}, errors.New("foreign")} {
		if errors.Is(err, target) {
			t.Fatalf("matched unrelated target %v", target)
		}
	}
	var typed *event.Error
	if !errors.As(err, &typed) || typed.Listener != "orders/audit.OnEvent" {
		t.Fatal("listener provenance lost")
	}
	for _, diagnostic := range []string{"listener_failed", "orders/audit.OnEvent", "invoke", "context deadline exceeded"} {
		if !strings.Contains(err.Error(), diagnostic) {
			t.Fatalf("missing diagnostic %q: %v", diagnostic, err)
		}
	}
}

func TestCycleRetainsActiveDeclarationPath(t *testing.T) {
	err := &event.Error{
		Code:     event.PublicationCycle,
		Listener: "a.OnA",
		Path:     []string{"a.OnA", "b.OnB", "a.OnA"},
	}
	if !errors.Is(err, &event.Error{
		Code: event.PublicationCycle,
	}) || !strings.Contains(err.Error(), "a.OnA -> b.OnB -> a.OnA") {
		t.Fatalf("cycle diagnostics lost: %v", err)
	}
}

func TestNilAndUnclassifiedErrorsDoNotMatchCodes(t *testing.T) {
	var absent *event.Error
	if absent.Error() != "<nil>" || absent.Unwrap() != nil || absent.Is(&event.Error{
		Code: event.InvalidState,
	}) {
		t.Fatal("nil error receiver is unsafe")
	}
	if (&event.Error{
		Code: event.InvalidState,
	}).Is(absent) || (&event.Error{}).Is(&event.Error{}) || (&event.Error{}).Unwrap() != nil {
		t.Fatal("nil or empty target was classified")
	}
}
