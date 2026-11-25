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
	"fmt"
	"path/filepath"
)

const (
	generatedFileIsAParentOfFormat = "generated file %s is a parent of %s"
)

// A generated file cannot also serve as a directory for another planned file.
// Detect this before staging can create directories or publish unrelated bytes.
func plannedParents(path string, desired map[string]bool) error {
	for parent := filepath.Dir(path); ; parent = filepath.Dir(parent) {
		if desired[parent] {
			return fmt.Errorf(generatedFileIsAParentOfFormat, parent, path)
		}
		if filepath.Dir(parent) == parent {
			return nil
		}
	}
}
