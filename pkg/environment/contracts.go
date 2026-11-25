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

// Origin carries provenance without carrying a property's potentially secret
// value. Line and Column are optional one-based positions; zero means unknown.
// Implementations must keep Source, Resource and Key free of secret payloads.
type Origin struct {
	Source   string
	Resource string
	Key      string
	Line     int
	Column   int
}

// Property keeps a value and its provenance together through resolution and
// binding. A zero Property is found-null only when Lookup also returns true.
type Property struct {
	Value  Value
	Origin Origin
}

// PropertySource is an immutable, concurrently readable in-memory source of raw
// values. Lookup performs no transport I/O and never expands placeholders; it
// validates exact keys and distinguishes missing (false, nil) from found-null or
// empty (true, nil). Errors do not permit fallback to lower-precedence sources.
// Keys returns all addresses, sorted and detached from internal storage.
type PropertySource interface {
	Name() string
	Lookup(string) (Property, bool, error)
	Keys() []string
}

// PropertyResolver resolves immutable property sources in highest-precedence
// order. The first present exact key wins, including null/empty; source failures
// stop traversal. Keys are atomic addresses: an object at app neither deletes nor
// substitutes for an independently addressed app.port. Keys is sorted/detached.
// Lookup recursively resolves strings within lists/objects and retains origin.
// Resolve strictly expands ${key}, ${key:default}, recursive references and nested
// defaults and nested keys such as ${${stage}.url:localhost}. Keys expand before
// lookup; nested key expansion shares cycle, depth, output and cancellation
// limits with values and defaults. Escaped \${...} remains literal without reparsing. Defaults apply only to missing keys; null or
// composite references fail. Missing references, cycles, malformed syntax and
// documented finite nesting/output limits return classified errors preserving
// causes without printing values. Lookup rejects invalid keys.
type PropertyResolver interface {
	Lookup(string) (Property, bool, error)
	Resolve(string) (string, error)
	Keys() []string
}

// Environment combines resolution with bootstrap-selected profiles fixed for
// its lifetime. Profile accessors return unique, validated, detached names. Active
// profiles are effective when supplied, otherwise defaults; unspecified defaults
// are ["default"]. AcceptsProfiles ORs validated names or !name predicates; an
// empty predicate list returns false. All predicates are validated, even after a
// match. Expression parsing is outside this contract. Remote properties must not
// silently change topology. Implementations support concurrent readers.
type Environment interface {
	PropertyResolver
	ActiveProfiles() []string
	DefaultProfiles() []string
	AcceptsProfiles(...string) (bool, error)
}

// Snapshot pins a coherent immutable environment and its publication revision.
// Capture once for related reads; revisions identify publications within one
// owner's lifetime rather than promising distributed ordering.
type Snapshot interface {
	Environment
	Revision() uint64
}

// SnapshotProvider exposes an owner's current immutable view. Publication must
// be atomic: prepare and validate the entire candidate before replacing the last
// good snapshot. Failed refresh retains the previous revision. Provider calls and
// captured snapshots must remain safe for concurrent readers.
type SnapshotProvider interface {
	Snapshot() Snapshot
}

// Bindable admits handwritten or generated typed binders without prescribing
// reflection, tags or a registry. BindProperties rejects a nil target without
// panic and commits target changes only after all resolution/conversion/validation
// succeeds. Prefix keys follow ValidateKey when nonempty; an empty prefix selects
// the root. Binding must not mutate the resolver. Target synchronization belongs
// to the caller; binders report classified failures with provenance and causes.
type Bindable interface {
	BindProperties(PropertyResolver, string) error
}
