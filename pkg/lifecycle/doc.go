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

// Package lifecycle defines independent resource capabilities shared across
// domains. Implementing one hook never requires implementing every other hook.
// Owners initialize before exposure and activate before readiness; they stop
// before destroying in reverse dependency ownership order. Failed Init permits
// Destroy, and any attempted Start permits Stop, so cleanup tolerates partial
// completion. A constructor that returns no owned result must clean itself up.
// Scheduling, cancellation, exactly-once tracking and aggregation belong to the
// owner. This package deliberately supplies no lifecycle runner.
package lifecycle
