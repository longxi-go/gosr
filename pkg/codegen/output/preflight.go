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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

func readOwned(path string, owns OwnershipFunc) ([]byte, bool, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if !info.Mode().IsRegular() {
		return nil, false, fmt.Errorf("generated output destination is not a regular file: %s", path)
	}
	current, err := os.ReadFile(path)
	if err != nil {
		return nil, false, err
	}
	if owns == nil || !owns(path, bytes.Clone(current)) {
		return nil, false, fmt.Errorf("refuse to modify unowned file %s", path)
	}
	return current, true, nil
}

func parentDirectories(path string) error {
	for dir := filepath.Dir(path); ; dir = filepath.Dir(dir) {
		info, err := os.Lstat(dir)
		if err == nil && !info.IsDir() {
			return fmt.Errorf("output parent is not a directory: %s", dir)
		}
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if next := filepath.Dir(dir); next == dir {
			return nil
		}
	}
}

func unchanged(path string, data []byte, owns OwnershipFunc) (bool, error) {
	current, exists, err := readOwned(path, owns)
	return exists && bytes.Equal(current, data), err
}
