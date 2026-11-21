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

package lifecycle

import "context"

// Initializer prepares an owned allocation; failure still permits Destroy.
type Initializer interface {
	Init(context.Context) error
}

// Starter activates a prepared object; even failed activation permits Stop.
type Starter interface {
	Start(context.Context) error
}

// Stopper releases activation work and must tolerate partial Start completion.
type Stopper interface {
	Stop(context.Context) error
}

// Destroyer releases an allocation and must tolerate partial Init completion.
type Destroyer interface {
	Destroy(context.Context) error
}
