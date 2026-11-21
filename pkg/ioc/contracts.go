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

import "context"

// DefinitionReader returns detached candidate metadata. Mutating a view must
// never change a registry or an application's bound graph. View order is not a
// dependency order; implementations must not derive binding from enumeration.
type DefinitionReader interface {
	Definitions() []BeanDefinition
}

// DefinitionRegistry atomically accepts metadata batches, rejects duplicate IDs
// and owns copies. It never creates objects; registries are optional capabilities.
// Names and aliases must be unique among activated definitions; raw candidates
// may overlap before activation. Contexts snapshot this metadata independently.
type DefinitionRegistry interface {
	DefinitionReader
	Register(...BeanDefinition) error
}

// BeanFactory resolves active visible candidates. Lookup selects a unique match
// or unique primary; missing/ambiguous matches have distinct error codes. An
// explicit name with a wrong contract is ContractMismatch. LookupAll retains all
// matches ordered by Order then ID, and an empty collection succeeds. After
// lifecycle/context gates, malformed selectors return InvalidSelector. Optional
// dependencies suppress only MissingBean caused by zero matching candidates,
// never ambiguity, construction failure or a named contract mismatch.
type BeanFactory interface {
	Lookup(context.Context, Selector) (any, error)
	LookupAll(context.Context, Selector) ([]any, error)
}

// State exposes lifecycle gates without prescribing a scheduler or locking model.
type State string

const (
	Created  State = "created"
	Starting State = "starting"
	Running  State = "running"
	Failed   State = "failed"
	Closing  State = "closing"
	Closed   State = "closed"
)

// ApplicationContext owns independent instances. External lookup is allowed only
// while Running; Start is single-attempt and Close is repeatable from any state.
// Failed startup compensates owned work; canceled shutdown waits do not abandon it.
type ApplicationContext interface {
	BeanFactory
	DefinitionReader
	Start(context.Context) error
	Close(context.Context) error
	State() State
}

// Scope owns one explicit lifetime and releases its objects on repeatable Close.
type Scope interface {
	Close(context.Context) error
}

// ScopeManager is optional. The returned context carries a capability tied to
// this owner. A definition is reused within one scope and isolated across
// scopes. Missing, foreign, closed or unsupported scopes fail with ScopeViolation.
type ScopeManager interface {
	OpenScope(context.Context, ScopeID) (context.Context, Scope, error)
}

// Provider resolves against its invocation context, never a captured request.
type Provider[T any] func(context.Context) (T, error)

// Optional distinguishes absence from a present typed zero value.
type Optional[T any] struct {
	Value   T
	Present bool
}

// Lease makes caller-owned resource lifetimes explicit. Acquisition is an
// implementation extension rather than a requirement on every BeanFactory.
type Lease[T any] interface {
	Value() T
	Close(context.Context) error
}
