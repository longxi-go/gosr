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

package event_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/longxi-go/gosr/pkg/event"
)

type orderPlaced struct {
	event.BaseEvent
	ID string
}

type orderAudit struct{}

type publisherFunc func(context.Context, event.Event) error

func (f publisherFunc) PublishEvent(ctx context.Context, value event.Event) error {
	return f(ctx, value)
}

func (orderAudit) OnEvent(_ context.Context, placed orderPlaced) error {
	fmt.Println("audit", placed.ID)
	return nil
}

func ExampleListener() {
	// This is a plain typed receiver, not a reference event dispatcher.
	var audit event.Listener[orderPlaced] = orderAudit{}
	if err := audit.OnEvent(context.Background(), orderPlaced{
		ID: "order-42",
	}); err != nil {
		panic(err)
	}
	// Output: audit order-42
}

func ExampleBaseEvent() {
	var placed event.Event = orderPlaced{
		ID: "order-42",
	}
	value, ok := placed.(orderPlaced)
	fmt.Println(value.ID, ok)
	// Output: order-42 true
}

func ExamplePublisher() {
	// A producer accepts the narrow capability, not an application container.
	place := func(ctx context.Context, publisher event.Publisher, id string) error {
		return publisher.PublishEvent(ctx, orderPlaced{
			ID: id,
		})
	}
	// This local recording adapter illustrates injection only. It is not a
	// dispatcher and makes no runtime conformance or lifecycle claim.
	var publisher event.Publisher = publisherFunc(func(_ context.Context, value event.Event) error {
		placed, ok := value.(orderPlaced)
		if !ok {
			return fmt.Errorf("unexpected example event %T", value)
		}
		fmt.Println("published", placed.ID)
		return nil
	})
	if err := place(context.Background(), publisher, "order-42"); err != nil {
		panic(err)
	}
	// Output: published order-42
}

func ExampleError() {
	err := fmt.Errorf("publish: %w", &event.Error{
		Code:     event.ListenerFailed,
		Listener: "orders/audit.OnEvent",
		Cause:    context.Canceled,
	})
	fmt.Println(errors.Is(err, &event.Error{
		Code: event.ListenerFailed,
	}))
	fmt.Println(errors.Is(err, context.Canceled))
	// Output:
	// true
	// true
}
