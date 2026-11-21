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

import "context"

// Event explicitly admits a value to publication. The exported marker permits
// independent modules to implement the contract without embedding a framework
// type. Publishers inspect types and must not call the marker for validation.
type Event interface {
	EventMarker()
}

// BaseEvent supplies optional embedding for business events. It imposes no
// timestamp, source, application reference or runtime-owned mutable state.
type BaseEvent struct{}

// EventMarker admits both BaseEvent values and pointers through embedding.
// A typed nil pointer is still an Event; publication must not call this method.
func (BaseEvent) EventMarker() {
}

// Publisher separates producers from listener discovery and ownership. Delivery
// is synchronous and sequential within one call; success means all matching
// declarations completed, including success when none match. Preparation failure
// starts no listener callbacks. Callback failure stops subsequent listeners and
// preserves its cause without undoing earlier side effects.
//
// Implementations support concurrent publication, propagate cancellation and
// track active declaration identities through the invocation context. Same-chain
// reentry returns PublicationCycle; independent calls share no visited state.
// A nil interface is InvalidEvent, while a typed nil follows normal Go matching.
type Publisher interface {
	PublishEvent(context.Context, Event) error
}

// Listener declares one typed receiving capability without prescribing discovery
// or registration. E follows Go type identity and assignability; values and
// pointers are not implicitly converted. Listener[Event] observes every event.
// Implementations must not retain the invocation context and must synchronize
// mutable state when publications can overlap. Payloads are read-only by contract.
type Listener[E Event] interface {
	OnEvent(context.Context, E) error
}
