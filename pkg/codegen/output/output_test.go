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
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/longxi-go/gosr/pkg/codegen/artifact"
	"github.com/longxi-go/gosr/pkg/codegen/output"
)

func owned(path string, data []byte) bool {
	return bytes.HasPrefix(data, []byte("generated"))
}

func TestCommitPreflightNoopAndStale(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.go")
	z := filepath.Join(dir, "z.go")
	old := []byte("generated old")
	if err := os.WriteFile(a, old, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(z, 0o700); err != nil {
		t.Fatal(err)
	}
	set, _ := artifact.New(artifact.File{
		Path:    a,
		Content: []byte("generated new"),
	}, artifact.File{
		Path:    z,
		Content: old,
	})
	if err := output.Commit(context.Background(), set, output.Options{
		Owns: owned,
	}); err == nil {
		t.Fatal("directory accepted")
	}
	data, _ := os.ReadFile(a)
	if !bytes.Equal(data, old) {
		t.Fatal("partial write before preflight")
	}
	set, _ = artifact.New(artifact.File{
		Path:    a,
		Content: old,
	})
	before, _ := os.Stat(a)
	if err := output.Commit(context.Background(), set, output.Options{
		Owns: owned,
	}); err != nil {
		t.Fatal(err)
	}
	after, _ := os.Stat(a)
	if !os.SameFile(before, after) || after.Mode().Perm() != 0o600 {
		t.Fatal("unchanged file replaced")
	}
	empty, _ := artifact.New()
	if err := output.Commit(context.Background(), empty, output.Options{
		Owns:  owned,
		Stale: []string{a},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(a); !os.IsNotExist(err) {
		t.Fatal("stale not removed")
	}
}

func TestRejectUnownedSymlinkAndCancellation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "user.go")
	if err := os.WriteFile(path, []byte("user"), 0o600); err != nil {
		t.Fatal(err)
	}
	set, _ := artifact.New(artifact.File{
		Path:    path,
		Content: []byte("generated"),
	})
	if err := output.Commit(context.Background(), set, output.Options{
		Owns: owned,
	}); err == nil {
		t.Fatal("unowned overwrite")
	}
	empty, _ := artifact.New()
	if err := output.Commit(context.Background(), empty, output.Options{
		Owns:  owned,
		Stale: []string{path},
	}); err == nil {
		t.Fatal("unowned deletion")
	}
	link := filepath.Join(dir, "link.go")
	if err := os.Symlink(path, link); err != nil {
		t.Fatal(err)
	}
	set, _ = artifact.New(artifact.File{
		Path:    link,
		Content: []byte("generated"),
	})
	if err := output.Commit(context.Background(), set, output.Options{
		Owns: owned,
	}); err == nil {
		t.Fatal("symlink followed")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := output.Commit(ctx, empty, output.Options{
		Owns: owned,
	}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestStagingCleanupWhenOwnershipChanges(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.go")
	if err := os.WriteFile(path, []byte("generated old"), 0o600); err != nil {
		t.Fatal(err)
	}
	calls := 0
	changingOwner := func(path string, data []byte) bool {
		calls++
		return calls == 1
	}
	set, _ := artifact.New(artifact.File{
		Path:    path,
		Content: []byte("generated new"),
	})
	if err := output.Commit(context.Background(), set, output.Options{
		Owns: changingOwner,
	}); err == nil {
		t.Fatal("lost ownership ignored")
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 {
		t.Fatal("staged file leaked", entries, err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != "generated old" {
		t.Fatal("unowned path replaced")
	}
}

func TestStalePreflightProtectsDesiredWrites(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.go")
	z := filepath.Join(dir, "z.go")
	if err := os.WriteFile(z, []byte("user"), 0o600); err != nil {
		t.Fatal(err)
	}
	set, _ := artifact.New(artifact.File{
		Path:    a,
		Content: []byte("generated new"),
	})
	if err := output.Commit(context.Background(), set, output.Options{
		Owns:  owned,
		Stale: []string{z},
	}); err == nil {
		t.Fatal("stale ownership ignored")
	}
	if _, err := os.Stat(a); !os.IsNotExist(err) {
		t.Fatal("publication happened before stale preflight")
	}
}

func TestStaleSelectionRejectsSymlinkParent(t *testing.T) {
	dir := t.TempDir()
	real := filepath.Join(dir, "real")
	if err := os.Mkdir(real, 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(real, "owned.go")
	if err := os.WriteFile(path, []byte("generated old"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link")
	if err := os.Symlink(real, link); err != nil {
		t.Fatal(err)
	}
	empty, _ := artifact.New()
	err := output.Commit(context.Background(), empty, output.Options{
		Owns:  owned,
		Stale: []string{filepath.Join(link, "owned.go")},
	})
	if err == nil {
		t.Fatal("followed symlink parent for stale removal")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("removed target outside selected path", err)
	}
}
