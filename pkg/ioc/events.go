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

// ContextRefreshed is an internal startup barrier for contexts providing event
// integration. Init and Start have completed, but Running is not yet committed:
// failure compensates startup. The supplied callback context admits internal
// startup operations; external lookup and business publication remain closed.
// It is not an application or Web-ready notification and carries no owner.
// Values and pointers are reserved against business publication.
type ContextRefreshed struct{}

// EventMarker admits the startup notification without runtime-owned payload.
func (ContextRefreshed) EventMarker() {
}

// ContextClosing is sent internally at most once after a successful Running
// commit, after admission closes and active work drains, before Stop and Destroy.
// All eligible existing listeners are attempted; failures join cleanup errors.
// Callbacks cannot resolve Beans, publish or create listeners. Essential cleanup
// belongs to lifecycle callbacks even if notification fails. Canceled Close
// waiters do not abandon owned cleanup. Values and pointers are reserved against
// business publication; Created or failed contexts do not send this notification.
type ContextClosing struct{}

// EventMarker admits the closing notification without retaining an application.
func (ContextClosing) EventMarker() {
}
