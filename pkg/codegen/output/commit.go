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
	"errors"

	"github.com/longxi-go/gosr/pkg/codegen/artifact"
)

const (
	artifactSetIsNilMessage = "artifact set is nil"
)

// Commit preflights every destination and stale file, stages all changed bytes,
// then replaces files in path order and removes explicitly selected stale files.
// It skips unchanged files. Cancellation and cleanup preserve their causes.
// Changed files retain existing permission bits; new files use 0644. Other file
// metadata is not preserved by replacement.
// Replacement follows os.Rename's platform guarantees. On Unix it is atomic per
// file, never across the set. A late I/O failure or
// cancellation can leave earlier files committed. Callers must serialize writers;
// this is not protection against a hostile concurrent filesystem actor.
func Commit(ctx context.Context, set *artifact.Set, options Options) (resultErr error) {
	if err := ctx.Err(); err != nil {
		return err
	}
	if set == nil {
		return errors.New(artifactSetIsNilMessage)
	}
	plan, err := prepareCommit(ctx, set, options)
	if err != nil {
		return err
	}
	staged := make(stagedFiles)
	defer func() {
		resultErr = errors.Join(resultErr, staged.cleanup())
	}()
	if err := staged.stage(ctx, plan.changes); err != nil {
		return err
	}
	if err := staged.replace(ctx, plan.changes, options.Owns); err != nil {
		return err
	}
	return removeStale(ctx, plan.stale, options.Owns)
}
