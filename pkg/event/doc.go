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

// Package event defines synchronous application event contracts independently of
// IoC, transport and annotation machinery. Event is an explicit nonempty boundary;
// BaseEvent embedding is optional. Producers depend on Publisher, receivers on
// Listener[E], and implementations own registration, matching and execution.
//
// Matching uses Go type identity and interface assignability. All matching
// declarations prepare before callbacks begin, then execute in ordered priority
// groups with one sampled order per object per round and stable declaration IDs
// for ties. An unrelated lazy receiver need not be created. The first business
// callback failure stops the call; preparation and already-completed side effects
// are not transactional. Nil interfaces are invalid; typed nil values can match.
//
// Invocation contexts carry cancellation and same-chain reentry tracking. One
// active declaration cannot reenter itself; different methods on one object can
// form nested calls. Concurrent calls have no global ordering. Implementations
// keep admitted callbacks and borrowed resources alive until completion; callers
// and listeners own synchronization of shared mutable state. Events are read-only
// by contract, with no deep copying. See docs/design/docs/contracts/events-ordering.md for
// normative semantics and the boundary between contract and runtime evidence.
package event
