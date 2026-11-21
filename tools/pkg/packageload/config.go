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
	"bytes"
	"context"

	"golang.org/x/tools/go/packages"
)

// Config establishes an invocation-owned input snapshot. Callers explicitly
// choose modes, flags and environment; no download/vendor/workspace policy is
// silently imposed on the Go command. Context must be non-nil.
func Config(ctx context.Context, dir string, overlay map[string][]byte) *packages.Config {
	return &packages.Config{
		Context: ctx,
		Dir:     dir,
		Overlay: cloneOverlay(overlay),
	}
}

func cloneOverlay(source map[string][]byte) map[string][]byte {
	if source == nil {
		return nil
	}
	copy := make(map[string][]byte, len(source))
	for path, data := range source {
		copy[path] = bytes.Clone(data)
	}
	return copy
}
