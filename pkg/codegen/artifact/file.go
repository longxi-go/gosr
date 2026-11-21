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

package artifact

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
)

// Kind distinguishes Go overlays from non-Go publication assets.
type Kind string

const (
	Go   Kind = "go"
	Data Kind = "data"
)

// File carries only generation mechanics. GoVersion is the output owner's
// language version (for example go1.27.0), not the generator's version.
type File struct {
	Path      string
	Content   []byte
	Kind      Kind
	GoVersion string
}

// Set owns immutable byte snapshots indexed by normalized path. Callers must
// choose one path base; converting relative paths to absolute belongs to them.
type Set struct {
	files []File
}

// New rejects duplicate normalized paths rather than silently replacing output.
func New(files ...File) (*Set, error) {
	s := &Set{}
	seen := make(map[string]bool)
	for _, file := range files {
		if file.Path == "" {
			return nil, fmt.Errorf("artifact path is empty")
		}
		file.Path = filepath.Clean(file.Path)
		if seen[file.Path] {
			return nil, fmt.Errorf("duplicate artifact %s", file.Path)
		}
		seen[file.Path] = true
		file.Content = bytes.Clone(file.Content)
		s.files = append(s.files, file)
	}
	sort.Slice(s.files, func(i, j int) bool {
		return s.files[i].Path < s.files[j].Path
	})
	return s, nil
}

// Files returns independent bytes so later formatting cannot mutate a snapshot.
func (s *Set) Files() []File {
	result := make([]File, len(s.files))
	for i, file := range s.files {
		file.Content = bytes.Clone(file.Content)
		result[i] = file
	}
	return result
}

// Overlay returns only Go sources. Every call owns its map and byte slices.
func (s *Set) Overlay() map[string][]byte {
	result := make(map[string][]byte)
	for _, file := range s.files {
		if file.Kind == Go {
			result[file.Path] = bytes.Clone(file.Content)
		}
	}
	return result
}
