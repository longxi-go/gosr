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

package pipeline_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/longxi-go/gosr/pkg/pipeline"
)

type applicationStepFunc func(context.Context, *state) error

type state struct {
	visits []string
}

type invalidNamedStep struct {
	name  string
	calls int
}

func (s *invalidNamedStep) Name() string {
	return s.name
}

func (s *invalidNamedStep) Execute(context.Context, *state) error {
	s.calls++
	return nil
}

func TestRunExecutesNamedStepsInOrderOnSharedState(t *testing.T) {
	var applicationCallback applicationStepFunc = func(_ context.Context, current *state) error {
		current.visits = append(current.visits, "first")
		return nil
	}
	var execute pipeline.StepFunc[*state] = applicationCallback
	current := &state{}
	err := pipeline.Run(context.Background(), current,
		pipeline.NewStep("first", execute),
		pipeline.NewStep("second", func(_ context.Context, current *state) error {
			current.visits = append(current.visits, "second")
			return nil
		}),
	)
	if err != nil || strings.Join(current.visits, ",") != "first,second" {
		t.Fatalf("visits=%v err=%v", current.visits, err)
	}
}

func TestRunAcceptsAnEmptySequence(t *testing.T) {
	if err := pipeline.Run(context.Background(), &state{}); err != nil {
		t.Fatal(err)
	}
}

func TestRunStopsOnFailureAndPreservesCause(t *testing.T) {
	cause := errors.New("conversion failed")
	current := &state{}
	err := pipeline.Run(context.Background(), current,
		pipeline.NewStep("convert", func(_ context.Context, current *state) error {
			current.visits = append(current.visits, "convert")
			return cause
		}),
		pipeline.NewStep("write", func(_ context.Context, current *state) error {
			current.visits = append(current.visits, "write")
			return nil
		}),
	)
	if !errors.Is(err, cause) || !strings.Contains(err.Error(), "convert") ||
		strings.Join(current.visits, ",") != "convert" {
		t.Fatalf("visits=%v err=%v", current.visits, err)
	}
}

func TestRunStopsWhenContextIsCanceledBetweenSteps(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	current := &state{}
	err := pipeline.Run(ctx, current,
		pipeline.NewStep("cancel", func(_ context.Context, current *state) error {
			current.visits = append(current.visits, "cancel")
			cancel()
			return nil
		}),
		pipeline.NewStep("unreachable", func(_ context.Context, current *state) error {
			current.visits = append(current.visits, "unreachable")
			return nil
		}),
	)
	if !errors.Is(err, context.Canceled) || strings.Join(current.visits, ",") != "cancel" {
		t.Fatalf("visits=%v err=%v", current.visits, err)
	}
}

func TestRunRejectsInvalidStepsBeforeExecution(t *testing.T) {
	tests := []struct {
		name string
		step pipeline.Step[*state]
	}{
		{
			name: "nil",
		},
		{
			name: "blank name",
			step: &invalidNamedStep{},
		},
		{
			name: "whitespace name",
			step: &invalidNamedStep{
				name: " \t ",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var laterCalls int
			err := pipeline.Run(context.Background(), &state{}, test.step,
				pipeline.NewStep("later", func(context.Context, *state) error {
					laterCalls++
					return nil
				}),
			)
			if err == nil || !strings.Contains(err.Error(), "step 0") || laterCalls != 0 {
				t.Fatalf("err=%v laterCalls=%d", err, laterCalls)
			}
			if invalid, ok := test.step.(*invalidNamedStep); ok && invalid.calls != 0 {
				t.Fatalf("invalid step executed %d times", invalid.calls)
			}
		})
	}
}

func TestNamedStepRejectsMissingExecutor(t *testing.T) {
	err := pipeline.Run(context.Background(), &state{}, pipeline.NewStep[*state]("missing", nil))
	if err == nil || !strings.Contains(err.Error(), "missing") {
		t.Fatalf("expected named missing-executor error, got %v", err)
	}
}
