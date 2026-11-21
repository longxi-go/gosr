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
	"os"

	"github.com/longxi-go/gosr/pkg/codegen/artifact"
)

func (s stagedFiles) replace(ctx context.Context, files []artifact.File, owns OwnershipFunc) error {
	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, _, err := readOwned(file.Path, owns); err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := os.Rename(s[file.Path], file.Path); err != nil {
			return err
		}
		delete(s, file.Path)
	}
	return nil
}

func removeStale(ctx context.Context, paths []string, owns OwnershipFunc) error {
	for _, path := range paths {
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, exists, err := readOwned(path, owns); err != nil {
			return err
		} else if !exists {
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
