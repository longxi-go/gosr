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
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/longxi-go/gosr/pkg/codegen/artifact"
)

type stagedFiles map[string]string

func (s stagedFiles) stage(ctx context.Context, files []artifact.File) error {
	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(file.Path), 0o755); err != nil {
			return err
		}
		mode, err := destinationMode(file.Path)
		if err != nil {
			return err
		}
		temp, err := os.CreateTemp(filepath.Dir(file.Path), ".codegen-*")
		if err != nil {
			return err
		}
		s[file.Path] = temp.Name()
		_, writeErr := temp.Write(file.Content)
		modeErr := temp.Chmod(mode)
		closeErr := temp.Close()
		if err := errors.Join(writeErr, modeErr, closeErr); err != nil {
			return fmt.Errorf("stage %s: %w", file.Path, err)
		}
	}
	return nil
}

func (s stagedFiles) cleanup() error {
	var result error
	paths := make([]string, 0, len(s))
	for path := range s {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		if err := os.Remove(s[path]); err != nil && !os.IsNotExist(err) {
			result = errors.Join(result, fmt.Errorf("staging cleanup: %w", err))
		}
	}
	return result
}
