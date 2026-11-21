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

package pipeline

import "context"

// Step performs one operation on caller-owned state. Implementations must be
// safe for the usage pattern chosen by their owner and return errors rather than
// relying on the runner to recover panics.
type Step[S any] interface {
	Name() string
	Execute(context.Context, S) error
}

// StepFunc is the function form of Step execution. Its alias identity allows
// application-defined function types with the same signature to be used directly.
type StepFunc[S any] = func(context.Context, S) error
