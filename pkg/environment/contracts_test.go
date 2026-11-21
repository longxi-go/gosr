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

package environment_test

import "github.com/longxi-go/gosr/pkg/environment"

// These external-package consumers enforce method signatures and interface
// composition without introducing a resolver or publication engine into gosr.
var (
	_ func(environment.PropertySource) string                                        = environment.PropertySource.Name
	_ func(environment.PropertySource, string) (environment.Property, bool, error)   = environment.PropertySource.Lookup
	_ func(environment.PropertySource) []string                                      = environment.PropertySource.Keys
	_ func(environment.PropertyResolver, string) (environment.Property, bool, error) = environment.PropertyResolver.Lookup
	_ func(environment.PropertyResolver, string) (string, error)                     = environment.PropertyResolver.Resolve
	_ func(environment.PropertyResolver) []string                                    = environment.PropertyResolver.Keys
	_ func(environment.Environment) []string                                         = environment.Environment.ActiveProfiles
	_ func(environment.Environment) []string                                         = environment.Environment.DefaultProfiles
	_ func(environment.Environment, ...string) (bool, error)                         = environment.Environment.AcceptsProfiles
	_ func(environment.Snapshot) uint64                                              = environment.Snapshot.Revision
	_ func(environment.SnapshotProvider) environment.Snapshot                        = environment.SnapshotProvider.Snapshot
	_ func(environment.Bindable, environment.PropertyResolver, string) error         = environment.Bindable.BindProperties
	_ func(environment.Environment) environment.PropertyResolver                     = func(e environment.Environment) environment.PropertyResolver {
		return e
	}
	_ func(environment.Snapshot) environment.Environment = func(s environment.Snapshot) environment.Environment {
		return s
	}
)
