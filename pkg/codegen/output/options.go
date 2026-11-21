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

// OwnershipFunc decides whether an existing file may be replaced or removed.
// Consumers own header/manifest policy; the committer owns filesystem safety.
// Each call receives detached bytes; mutation cannot alter the content snapshot.
type OwnershipFunc = func(path string, current []byte) bool

// Options selects explicit ownership and stale paths. Stale never means scan a
// directory. Missing paths are omitted during preflight; paths also present in
// the new set are retained. The zero ownership
// predicate refuses modifications to existing files.
type Options struct {
	Owns  OwnershipFunc
	Stale []string
}
