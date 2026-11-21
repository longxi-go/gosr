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

// Package pipeline defines a small generic contract for executing ordered steps
// over invocation-owned state. It provides deterministic sequential execution,
// cancellation checks between steps and contextual errors without owning the
// steps, their state or the work they perform.
//
// Pass a non-nil context and valid Step values; a Step interface must not hold a
// typed-nil pointer. Names must be nonblank and stable, but need not be unique:
// a compiler may deliberately execute the same formatting operation twice.
// Nil step interfaces, blank names and missing adapter functions return errors.
// Run does not synchronize caller state, recover panics or roll back work.
// An empty sequence succeeds, and cancellation after the final step does not
// replace that step's result. Steps own cancellation within their operation.
package pipeline
