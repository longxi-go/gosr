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

// Package ioc defines the IoC domain of Go Shared Rules. It provides portable
// declaration values and small capability interfaces, not a container engine.
//
// Definitions contain no constructors, reflection tokens or instance state.
// Implementations own graph binding, execution and resources; consumers depend
// only on the lookup, registry, context or scope capability they actually need.
// Named initialization prerequisites are separate from injected dependencies.
// Their completion covers construction and Init, not Start or business readiness;
// implementations validate ownership and release same-owner dependents first.
// ContextRefreshed and ContextClosing describe optional event integration through
// the independent event package; ApplicationContext itself is not expanded.
// See docs/design/docs/contracts/ioc-contracts.md for normative lifecycle and failure rules.
package ioc
