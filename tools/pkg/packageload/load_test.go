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

package packageload_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/longxi-go/gosr/tools/pkg/packageload"
	"golang.org/x/tools/go/packages"
)

func TestOverlaySnapshotAndDiagnostics(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module fixture\ngo 1.27.0\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "a.go")
	if err := os.WriteFile(path, []byte("package fixture\nvar Value = missing\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	overlay := map[string][]byte{path: []byte("package fixture\nvar Value = 1\n")}
	config := packageload.Config(context.Background(), dir, overlay)
	// This temporary consumer owns its module policy independently of the
	// repository workspace, persistent Go settings and external package drivers.
	config.Env = append(os.Environ(),
		"GOWORK=off",
		"GOFLAGS=",
		"GOENV=off",
		"GOPROXY=off",
		"GOTOOLCHAIN=local",
		"GOPACKAGESDRIVER=off",
	)
	config.Mode = packages.NeedName | packages.NeedTypes
	overlay[path][0] = 'X'
	pkgs, err := packageload.Load(config, ".")
	if err != nil {
		t.Fatal(err)
	}
	if err := packageload.Errors("validate", pkgs); err != nil {
		t.Fatal(err)
	}
	config.Overlay = nil
	pkgs, err = packageload.Load(config, ".")
	if err != nil {
		t.Fatal(err)
	}
	if err := packageload.Errors("validate", pkgs); err == nil {
		t.Fatal("type errors lost")
	}
}
