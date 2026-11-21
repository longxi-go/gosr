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

// Package environment defines portable property, resolution, profile and snapshot
// contracts independently of a container or configuration engine. It supplies
// closed immutable values, provenance, validation and value-redacting errors; it
// supplies no sources, placeholder engine, watcher, loader or binding runtime.
//
// Sources own immutable raw data, resolvers own strict expansion, and snapshot
// providers own atomic publication after candidate validation. A missing key is
// distinguished from an explicit null or empty string through Lookup's bool.
// Keys are exact atomic addresses, and first-present precedence never falls
// through on null, invalid data or source failure. Implementations must preserve
// these semantics while choosing their own transport and finite resource limits.
// Placeholder tokens and DefaultProfile name the shared grammar and fallback;
// ValueDecoderFunc admits typed conversions without prescribing their strategy.
//
// Value collections detach on input/output and may be shared concurrently.
// Published environments and snapshots support concurrent readers; typed binding
// is transactional, with caller-owned target synchronization. These interfaces
// specify runtime obligations for implementations, not an engine conformance
// claim established by the local helper tests.
package environment
