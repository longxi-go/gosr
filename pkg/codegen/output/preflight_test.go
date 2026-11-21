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
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/longxi-go/gosr/pkg/codegen/artifact"
	"github.com/longxi-go/gosr/pkg/codegen/output"
)

func TestPreflightRejectsPlannedFileAsParent(t *testing.T) {
	dir := t.TempDir()
	sibling := filepath.Join(dir, "a.go")
	if err := os.WriteFile(sibling, []byte("generated old"), 0o600); err != nil {
		t.Fatal(err)
	}
	parent := filepath.Join(dir, "z.go")
	set, err := artifact.New(
		artifact.File{
			Path:    sibling,
			Content: []byte("generated new"),
		},
		artifact.File{
			Path:    parent,
			Content: []byte("parent"),
		},
		artifact.File{
			Path:    filepath.Join(parent, "child.go"),
			Content: []byte("child"),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	err = output.Commit(context.Background(), set, output.Options{
		Owns: owned,
	})
	if err == nil {
		t.Fatal("accepted an internally contradictory plan")
	}
	t.Logf("commit error: %v", err)
	if info, statErr := os.Stat(parent); statErr == nil {
		t.Errorf("preflight-invalid plan created a directory at planned file path: dir=%v", info.IsDir())
	}
	data, readErr := os.ReadFile(sibling)
	if readErr != nil || string(data) != "generated old" {
		t.Errorf("preflight-invalid plan changed sibling: %q, %v", data, readErr)
	}
}

func TestOwnershipPredicateCannotInfluenceNoopDecision(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.go")
	if err := os.WriteFile(path, []byte("generated old"), 0o600); err != nil {
		t.Fatal(err)
	}
	set, err := artifact.New(artifact.File{
		Path:    path,
		Content: []byte("generated new"),
	})
	if err != nil {
		t.Fatal(err)
	}
	mutateSnapshot := func(_ string, current []byte) bool {
		copy(current, "generated new")
		return true
	}
	err = output.Commit(context.Background(), set, output.Options{
		Owns: mutateSnapshot,
	})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "generated new" {
		t.Fatalf("commit reported success with %q, %v; want new content", data, err)
	}
}

func TestCancellationBeforeReplacingFinalFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.go")
	if err := os.WriteFile(path, []byte("generated old"), 0o600); err != nil {
		t.Fatal(err)
	}
	set, err := artifact.New(artifact.File{
		Path:    path,
		Content: []byte("generated new"),
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	calls := 0
	cancelDuringRecheck := func(_ string, _ []byte) bool {
		calls++
		if calls == 2 {
			cancel()
		}
		return true
	}
	err = output.Commit(ctx, set, output.Options{
		Owns: cancelDuringRecheck,
	})
	data, readErr := os.ReadFile(path)
	if !errors.Is(err, context.Canceled) || readErr != nil || string(data) != "generated old" {
		t.Fatalf("commit error = %v, file = %q, read error = %v, ownership calls = %d", err, data, readErr, calls)
	}
}

func TestCancellationBeforeRemovingFinalStaleFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "old.go")
	if err := os.WriteFile(path, []byte("generated old"), 0o600); err != nil {
		t.Fatal(err)
	}
	set, err := artifact.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	calls := 0
	cancelDuringRecheck := func(_ string, _ []byte) bool {
		calls++
		if calls == 2 {
			cancel()
		}
		return true
	}
	err = output.Commit(ctx, set, output.Options{
		Owns:  cancelDuringRecheck,
		Stale: []string{path},
	})
	data, readErr := os.ReadFile(path)
	if !errors.Is(err, context.Canceled) || readErr != nil || string(data) != "generated old" {
		t.Fatalf("commit error = %v, file = %q, read error = %v", err, data, readErr)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 {
		t.Fatalf("unexpected output after cancellation: %v, %v", entries, err)
	}
}

func TestMissingStalePathsDoNotConflictWithDesiredPaths(t *testing.T) {
	for _, ancestor := range []bool{true, false} {
		dir := t.TempDir()
		parent := filepath.Join(dir, "parent")
		child := filepath.Join(parent, "child.go")
		desired, stale := parent, child
		if ancestor {
			desired, stale = child, parent
		}
		set, err := artifact.New(artifact.File{
			Path:    desired,
			Content: []byte("generated new"),
		})
		if err != nil {
			t.Fatal(err)
		}
		err = output.Commit(context.Background(), set, output.Options{
			Stale: []string{stale},
		})
		data, readErr := os.ReadFile(desired)
		if err != nil || readErr != nil || string(data) != "generated new" {
			t.Errorf("stale ancestor=%v: commit error = %v, desired = %q, read err = %v", ancestor, err, data, readErr)
		}
	}
}
