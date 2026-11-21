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

package artifact_test

import (
	"testing"

	"github.com/longxi-go/gosr/pkg/codegen/artifact"
)

func TestSnapshotsAndDuplicates(t *testing.T) {
	data := []byte("original")
	s, err := artifact.New(artifact.File{
		Path:      "b.go",
		Content:   data,
		Kind:      artifact.Go,
		GoVersion: "go1.27.0",
	}, artifact.File{
		Path:    "a.json",
		Content: []byte("{}"),
		Kind:    artifact.Data,
	})
	if err != nil {
		t.Fatal(err)
	}
	data[0] = 'X'
	list := s.Files()
	if list[0].Path != "a.json" || string(list[1].Content) != "original" {
		t.Fatal(list)
	}
	list[1].Content[0] = 'Y'
	overlay := s.Overlay()
	overlay["b.go"][0] = 'Z'
	if string(s.Files()[1].Content) != "original" || len(overlay) != 1 {
		t.Fatal("snapshot aliases caller bytes")
	}
	if _, err := artifact.New(artifact.File{
		Path: "a/../b.go",
	}, artifact.File{
		Path: "b.go",
	}); err == nil {
		t.Fatal("normalized duplicate accepted")
	}
	if _, err := artifact.New(artifact.File{}); err == nil {
		t.Fatal("empty path accepted")
	}
}
