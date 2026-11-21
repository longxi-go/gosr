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

package output

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/longxi-go/gosr/pkg/codegen/artifact"
)

type commitPlan struct {
	changes []artifact.File
	stale   []string
}

func prepareCommit(ctx context.Context, set *artifact.Set, options Options) (commitPlan, error) {
	var plan commitPlan
	desired := make(map[string]bool)
	files := set.Files()
	for _, file := range files {
		desired[file.Path] = true
	}
	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return plan, err
		}
		if err := plannedParents(file.Path, desired); err != nil {
			return plan, err
		}
		if err := parentDirectories(file.Path); err != nil {
			return plan, err
		}
		same, err := unchanged(file.Path, file.Content, options.Owns)
		if err != nil {
			return plan, err
		}
		if !same {
			plan.changes = append(plan.changes, file)
		}
	}
	stale := make(map[string]bool)
	for _, path := range options.Stale {
		if path == "" {
			return plan, fmt.Errorf("stale path is empty")
		}
		path = filepath.Clean(path)
		if !desired[path] {
			stale[path] = true
		}
	}
	paths := make([]string, 0, len(stale))
	for path := range stale {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		if err := ctx.Err(); err != nil {
			return plan, err
		}
		if err := parentDirectories(path); err != nil {
			return plan, err
		}
		if _, exists, err := readOwned(path, options.Owns); err != nil {
			return plan, err
		} else if exists {
			plan.stale = append(plan.stale, path)
		}
	}
	return plan, nil
}
