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

package output_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/longxi-go/gosr/pkg/codegen/artifact"
	"github.com/longxi-go/gosr/pkg/codegen/output"
)

func TestChangedOutputPreservesExistingPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not implement POSIX file permission bits")
	}
	for _, mode := range []os.FileMode{0o600, 0o640, 0o755} {
		t.Run(mode.String(), func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "generated.data")
			if err := os.WriteFile(path, []byte("owned:old"), mode); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(path, mode); err != nil {
				t.Fatal(err)
			}
			set, err := artifact.New(artifact.File{
				Path:    path,
				Content: []byte("owned:new"),
				Kind:    artifact.Data,
			})
			if err != nil {
				t.Fatal(err)
			}
			if err := output.Commit(context.Background(), set, output.Options{
				Owns: ownsPermissionFixture,
			}); err != nil {
				t.Fatal(err)
			}
			info, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode().Perm() != mode {
				t.Fatalf("permissions = %o, want %o", info.Mode().Perm(), mode)
			}
			data, err := os.ReadFile(path)
			if err != nil || string(data) != "owned:new" {
				t.Fatalf("updated content = %q, error = %v", data, err)
			}
		})
	}
}

func ownsPermissionFixture(_ string, data []byte) bool {
	return bytes.HasPrefix(data, []byte("owned:"))
}
