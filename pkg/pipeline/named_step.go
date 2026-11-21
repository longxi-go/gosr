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

import (
	"context"
	"errors"
)

type namedStep[S any] struct {
	name    string
	execute StepFunc[S]
}

// NewStep adapts a named operation to Step. A nil operation is diagnosed when
// the step runs, preserving the same fail-fast error path as other step failures.
func NewStep[S any](name string, execute StepFunc[S]) Step[S] {
	return namedStep[S]{
		name:    name,
		execute: execute,
	}
}

func (s namedStep[S]) Name() string {
	return s.name
}

func (s namedStep[S]) Execute(ctx context.Context, state S) error {
	if s.execute == nil {
		return errors.New("step executor is nil")
	}
	return s.execute(ctx, state)
}
