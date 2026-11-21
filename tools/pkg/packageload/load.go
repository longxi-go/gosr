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

package packageload

import (
	"slices"

	"golang.org/x/tools/go/packages"
)

// Load copies mutable input slices/maps again at the execution boundary. A nil
// config uses go/packages defaults. It does not reject package diagnostics;
// discovery may intentionally load incomplete first-generation source. Call
// Errors when type correctness is required.
func Load(config *packages.Config, patterns ...string) ([]*packages.Package, error) {
	var snapshot *packages.Config
	if config != nil {
		copy := *config
		copy.Env = slices.Clone(config.Env)
		copy.BuildFlags = slices.Clone(config.BuildFlags)
		copy.Overlay = cloneOverlay(config.Overlay)
		snapshot = &copy
	}
	if snapshot != nil && snapshot.Context != nil {
		if err := snapshot.Context.Err(); err != nil {
			return nil, err
		}
	}
	return packages.Load(snapshot, patterns...)
}
